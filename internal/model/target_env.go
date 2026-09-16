package model

import (
	"strings"
	"time"
)

const (
	EnvActive    = "active"
	EnvDisabled  = "disabled"
)

type TargetEnv struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	BaseURL   string    `json:"base_url"`
	AuthToken string    `json:"auth_token"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (e *TargetEnv) Validate() error {
	e.Name = strings.TrimSpace(e.Name)
	e.BaseURL = strings.TrimSpace(e.BaseURL)
	if e.Name == "" {
		return NewValidationError("name", "环境名称不能为空")
	}
	if e.BaseURL == "" {
		return NewValidationError("base_url", "基础 URL 不能为空")
	}
	if e.Status == "" {
		e.Status = EnvActive
	}
	if e.Status != EnvActive && e.Status != EnvDisabled {
		return NewValidationError("status", "环境状态不合法")
	}
	return nil
}

type TargetEnvFilter struct {
	Name    string
	Status  string
	Keyword string
}

func (f TargetEnvFilter) Match(e *TargetEnv) bool {
	if f.Name != "" && e.Name != f.Name {
		return false
	}
	if f.Status != "" && e.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(e.Name), k) &&
			!strings.Contains(strings.ToLower(e.BaseURL), k) {
			return false
		}
	}
	return true
}
