package model

import (
	"strings"
	"time"
)

type ProjectConfig struct {
	ID              string    `json:"id"`
	Key             string    `json:"key"`
	Value           string    `json:"value"`
	Description     string    `json:"description"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (p *ProjectConfig) Validate() error {
	p.Key = strings.TrimSpace(p.Key)
	if p.Key == "" {
		return NewValidationError("key", "配置键不能为空")
	}
	return nil
}

type ProjectConfigFilter struct {
	Key     string
	Keyword string
}

func (f ProjectConfigFilter) Match(p *ProjectConfig) bool {
	if f.Key != "" && p.Key != f.Key {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(p.Key), k) &&
			!strings.Contains(strings.ToLower(p.Description), k) {
			return false
		}
	}
	return true
}
