package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Roles. admin may manage infrastructure (nodes/logic/topics); user may only
// create and view their own deliveries.
const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

// claimsKey is the gin context key under which verified JWT claims are stored.
const claimsKey = "claims"

// tokenTTL is how long an issued login token stays valid.
const tokenTTL = 24 * time.Hour

// AuthEnabled reports whether JWT authentication should be wired in. It is
// gated behind the AUTH_ENABLED env var so tests can disable it. Defaults ON;
// only an explicit "0"/"false"/"no" turns it off.
func AuthEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("AUTH_ENABLED"))) {
	case "0", "false", "no":
		return false
	default:
		return true
	}
}

// jwtSecret returns the HMAC signing secret from the environment.
func jwtSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}

// JWTAuthMiddleware validates the HS256 bearer token and stores the verified
// claims on the context. Unauthenticated requests are rejected.
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")

		claims, err := parseAndVerifyJWT(token, jwtSecret())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token: " + err.Error()})
			return
		}
		c.Set(claimsKey, claims)
		c.Next()
	}
}

// RequireRole returns middleware that allows the request only if the caller's
// "role" claim is one of the allowed roles. Must run after JWTAuthMiddleware.
func RequireRole(allowed ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := claimRole(c)
		for _, r := range allowed {
			if role == r {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient role"})
	}
}

// claimRole returns the caller's role from verified claims, or "" if absent.
func claimRole(c *gin.Context) string {
	if claims, ok := claimsFrom(c); ok {
		if role, ok := claims["role"].(string); ok {
			return role
		}
	}
	return ""
}

// claimSubject returns the caller's "sub" (the owning user's email), or "".
func claimSubject(c *gin.Context) string {
	if claims, ok := claimsFrom(c); ok {
		if sub, ok := claims["sub"].(string); ok {
			return sub
		}
	}
	return ""
}

// claimsFrom retrieves verified claims previously set by JWTAuthMiddleware.
func claimsFrom(c *gin.Context) (map[string]interface{}, bool) {
	v, exists := c.Get(claimsKey)
	if !exists {
		return nil, false
	}
	claims, ok := v.(map[string]interface{})
	return claims, ok
}

// authzEnforced reports whether claims-based authorization should run. When auth
// is disabled (tests/dev) there are no claims, so ownership/role checks pass.
func authzEnforced(c *gin.Context) bool {
	_, ok := claimsFrom(c)
	return ok
}

// requireOwnerOrAdmin aborts the request unless the caller is an admin or owns
// the resource (their subject email matches ownerEmail). Returns false (and
// has already written the response) when access is denied.
func requireOwnerOrAdmin(c *gin.Context, ownerEmail string) bool {
	if !authzEnforced(c) {
		return true
	}
	if claimRole(c) == RoleAdmin || claimSubject(c) == ownerEmail {
		return true
	}
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "not the owner"})
	return false
}

/**************************************************************/
/* Login                                                      */
/**************************************************************/

// loginRequest is the credential payload for issuing a token.
type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login validates credentials and issues a signed HS256 JWT. Credentials are
// checked against env-provided accounts so no secrets are hardcoded:
//   ADMIN_EMAIL / ADMIN_PASSWORD -> role admin
//   USER_EMAIL  / USER_PASSWORD  -> role user
func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, ok := authenticate(req.Email, req.Password)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := issueJWT(req.Email, role, jwtSecret())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token, "role": role})
}

// authenticate checks credentials against env-configured accounts and returns
// the matched role.
func authenticate(email, password string) (string, bool) {
	if matchAccount(email, password, "ADMIN_EMAIL", "ADMIN_PASSWORD") {
		return RoleAdmin, true
	}
	if matchAccount(email, password, "USER_EMAIL", "USER_PASSWORD") {
		return RoleUser, true
	}
	return "", false
}

// matchAccount reports whether email/password match the credentials in the
// given env vars. An unset email or password disables that account. Comparisons
// are constant-time to avoid leaking timing information about the secrets.
func matchAccount(email, password, emailEnv, passEnv string) bool {
	wantEmail := os.Getenv(emailEnv)
	wantPass := os.Getenv(passEnv)
	if wantEmail == "" || wantPass == "" {
		return false
	}
	emailOK := hmac.Equal([]byte(email), []byte(wantEmail))
	passOK := hmac.Equal([]byte(password), []byte(wantPass))
	return emailOK && passOK
}

// issueJWT builds and signs an HS256 token for the given subject and role.
func issueJWT(subject, role string, secret []byte) (string, error) {
	header := map[string]string{"alg": "HS256", "typ": "JWT"}
	claims := map[string]interface{}{
		"sub":  subject,
		"role": role,
		"exp":  time.Now().Add(tokenTTL).Unix(),
	}
	encHeader, err := encodeSegment(header)
	if err != nil {
		return "", err
	}
	encClaims, err := encodeSegment(claims)
	if err != nil {
		return "", err
	}
	signingInput := encHeader + "." + encClaims
	sig := base64.RawURLEncoding.EncodeToString(signHS256(signingInput, secret))
	return signingInput + "." + sig, nil
}

// encodeSegment JSON-encodes then base64url-encodes a JWT segment.
func encodeSegment(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// parseAndVerifyJWT verifies an HS256 JWT and returns its claims.
func parseAndVerifyJWT(token string, secret []byte) (map[string]interface{}, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("malformed token")
	}

	signingInput := parts[0] + "." + parts[1]
	expectedSig := signHS256(signingInput, secret)
	gotSig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, errors.New("bad signature encoding")
	}
	if !hmac.Equal(expectedSig, gotSig) {
		return nil, errors.New("signature mismatch")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("bad payload encoding")
	}
	var claims map[string]interface{}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, errors.New("bad claims json")
	}

	if expVal, ok := claims["exp"]; ok {
		if exp, ok := expVal.(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return nil, errors.New("token expired")
			}
		}
	}
	return claims, nil
}

func signHS256(signingInput string, secret []byte) []byte {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(signingInput))
	return mac.Sum(nil)
}
