package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

const testSecret = "test-signing-secret"

// makeToken builds an HS256 JWT for tests with an arbitrary role and expiry,
// so we can exercise valid, expired, and role-mismatch cases.
func makeToken(t *testing.T, secret, sub, role string, exp time.Time) string {
	t.Helper()
	enc := func(v interface{}) string {
		b, _ := json.Marshal(v)
		return base64.RawURLEncoding.EncodeToString(b)
	}
	header := enc(map[string]string{"alg": "HS256", "typ": "JWT"})
	claims := enc(map[string]interface{}{"sub": sub, "role": role, "exp": exp.Unix()})
	signingInput := header + "." + claims
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signingInput))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return signingInput + "." + sig
}

// protectedRouter builds a router with the JWT middleware and an admin-gated
// route, mirroring how main.go wires authorization.
func protectedRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(JWTAuthMiddleware())
	r.GET("/admin", RequireRole(RoleAdmin), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	r.GET("/any", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"sub": claimSubject(c)})
	})
	return r
}

func doGet(r *gin.Engine, path, token string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestJWTAuthMiddleware(t *testing.T) {
	t.Setenv("JWT_SECRET", testSecret)
	r := protectedRouter()

	validAdmin := makeToken(t, testSecret, "admin@x", RoleAdmin, time.Now().Add(time.Hour))
	validUser := makeToken(t, testSecret, "user@x", RoleUser, time.Now().Add(time.Hour))
	expired := makeToken(t, testSecret, "admin@x", RoleAdmin, time.Now().Add(-time.Hour))
	wrongSecret := makeToken(t, "other-secret", "admin@x", RoleAdmin, time.Now().Add(time.Hour))

	tests := []struct {
		name     string
		path     string
		token    string
		wantCode int
	}{
		{"missing token rejected", "/any", "", http.StatusUnauthorized},
		{"malformed token rejected", "/any", "not-a-jwt", http.StatusUnauthorized},
		{"bad signature rejected", "/any", wrongSecret, http.StatusUnauthorized},
		{"expired token rejected", "/any", expired, http.StatusUnauthorized},
		{"valid token accepted", "/any", validUser, http.StatusOK},
		{"admin role passes admin gate", "/admin", validAdmin, http.StatusOK},
		{"user role blocked from admin gate", "/admin", validUser, http.StatusForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := doGet(r, tt.path, tt.token)
			if w.Code != tt.wantCode {
				t.Errorf("%s: got %d, want %d (body=%s)", tt.path, w.Code, tt.wantCode, w.Body.String())
			}
		})
	}
}

func TestLoginIssuesUsableToken(t *testing.T) {
	t.Setenv("JWT_SECRET", testSecret)
	t.Setenv("ADMIN_EMAIL", "admin@x")
	t.Setenv("ADMIN_PASSWORD", "s3cret")

	gin.SetMode(gin.TestMode)
	h := &Handler{}
	login := gin.New()
	login.POST("/auth/login", h.Login)

	// Wrong credentials are rejected.
	bad := httptest.NewRequest(http.MethodPost, "/auth/login",
		jsonBody(`{"email":"admin@x","password":"wrong"}`))
	bad.Header.Set("Content-Type", "application/json")
	wbad := httptest.NewRecorder()
	login.ServeHTTP(wbad, bad)
	if wbad.Code != http.StatusUnauthorized {
		t.Fatalf("bad creds: got %d, want 401", wbad.Code)
	}

	// Correct credentials yield a token that the middleware accepts.
	good := httptest.NewRequest(http.MethodPost, "/auth/login",
		jsonBody(`{"email":"admin@x","password":"s3cret"}`))
	good.Header.Set("Content-Type", "application/json")
	wgood := httptest.NewRecorder()
	login.ServeHTTP(wgood, good)
	if wgood.Code != http.StatusOK {
		t.Fatalf("good creds: got %d, want 200 (body=%s)", wgood.Code, wgood.Body.String())
	}

	var resp struct {
		Token string `json:"token"`
		Role  string `json:"role"`
	}
	if err := json.Unmarshal(wgood.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode login response: %v", err)
	}
	if resp.Role != RoleAdmin {
		t.Errorf("role = %q, want admin", resp.Role)
	}

	r := protectedRouter()
	if w := doGet(r, "/admin", resp.Token); w.Code != http.StatusOK {
		t.Errorf("issued token rejected by admin gate: %d (body=%s)", w.Code, w.Body.String())
	}
}

// jsonBody wraps a JSON string as a request body reader.
func jsonBody(s string) *strings.Reader { return strings.NewReader(s) }
