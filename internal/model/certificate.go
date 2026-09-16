package model

import (
	"strings"
	"time"
)

type Certificate struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	CertPEM     string    `json:"cert_pem"`
	KeyPEM      string    `json:"key_pem"`
	ExpiresAt   time.Time `json:"expires_at"`
	CreatedAt   time.Time `json:"created_at"`
}

func (c *Certificate) Validate() error {
	c.Name = strings.TrimSpace(c.Name)
	if c.Name == "" {
		return NewValidationError("name", "证书名称不能为空")
	}
	if c.CertPEM == "" {
		return NewValidationError("cert_pem", "证书内容不能为空")
	}
	return nil
}

type CertificateFilter struct {
	Name    string
	Keyword string
}

func (f CertificateFilter) Match(c *Certificate) bool {
	if f.Name != "" && c.Name != f.Name {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(c.Name), k) {
			return false
		}
	}
	return true
}
