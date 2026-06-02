// Package notify sends the "delivery complete" email triggered when a flight
// reports a successful delivery. All connection details and credentials come
// from the environment so no secrets live in source.
package notify

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

// SMTPConfig holds the mail server connection details. Credentials are optional:
// dev brokers like MailHog accept unauthenticated mail, so auth is used only
// when both user and password are set.
type SMTPConfig struct {
	Host string
	Port string
	User string
	Pass string
	From string
}

// LoadSMTPConfig reads SMTP settings from SMTP_* env vars, defaulting to a local
// MailHog instance (localhost:1025) for development.
func LoadSMTPConfig() SMTPConfig {
	return SMTPConfig{
		Host: getenv("SMTP_HOST", "localhost"),
		Port: getenv("SMTP_PORT", "1025"),
		User: os.Getenv("SMTP_USER"),
		Pass: os.Getenv("SMTP_PASS"),
		From: getenv("SMTP_FROM", "airpost@localhost"),
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// addr returns the host:port the SMTP client dials.
func (c SMTPConfig) addr() string {
	return c.Host + ":" + c.Port
}

// auth returns SMTP auth, or nil when credentials are absent (e.g. MailHog).
func (c SMTPConfig) auth() smtp.Auth {
	if c.User == "" || c.Pass == "" {
		return nil
	}
	return smtp.PlainAuth("", c.User, c.Pass, c.Host)
}

// SendDeliveredEmail emails the recipient that their order was delivered.
func SendDeliveredEmail(cfg SMTPConfig, to, orderNum string) error {
	subject := "AirPost delivery complete - order " + orderNum
	body := fmt.Sprintf("Your AirPost order %s has been delivered.", orderNum)
	msg := buildMessage(cfg.From, to, subject, body)
	return smtp.SendMail(cfg.addr(), cfg.auth(), cfg.From, []string{to}, []byte(msg))
}

// buildMessage assembles a minimal RFC 5322 message.
func buildMessage(from, to, subject, body string) string {
	headers := []string{
		"From: " + from,
		"To: " + to,
		"Subject: " + subject,
	}
	return strings.Join(headers, "\r\n") + "\r\n\r\n" + body
}
