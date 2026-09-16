package model

import (
	"strings"
	"time"
)

const (
	WebhookStatusActive   = "active"
	WebhookStatusInactive = "inactive"
)

type WebhookConfig struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	Secret    string    `json:"secret"`
	Events    string    `json:"events"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (w *WebhookConfig) Validate() error {
	w.Name = strings.TrimSpace(w.Name)
	w.URL = strings.TrimSpace(w.URL)
	if w.Name == "" {
		return NewValidationError("name", "Webhook 名称不能为空")
	}
	if w.URL == "" {
		return NewValidationError("url", "URL 不能为空")
	}
	if w.Status == "" {
		w.Status = WebhookStatusActive
	}
	if w.Status != WebhookStatusActive && w.Status != WebhookStatusInactive {
		return NewValidationError("status", "状态不合法")
	}
	return nil
}

type WebhookConfigFilter struct {
	Name    string
	Status  string
	Keyword string
}

func (f WebhookConfigFilter) Match(w *WebhookConfig) bool {
	if f.Name != "" && w.Name != f.Name {
		return false
	}
	if f.Status != "" && w.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(w.Name), k) &&
			!strings.Contains(strings.ToLower(w.URL), k) {
			return false
		}
	}
	return true
}
