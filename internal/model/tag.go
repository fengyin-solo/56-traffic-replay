package model

import (
	"strings"
	"time"
)

type Tag struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	TargetType  string    `json:"target_type"`
	TargetID    string    `json:"target_id"`
	CreatedAt   time.Time `json:"created_at"`
}

func (t *Tag) Validate() error {
	t.Name = strings.TrimSpace(t.Name)
	if t.Name == "" {
		return NewValidationError("name", "标签名称不能为空")
	}
	if t.TargetType == "" {
		return NewValidationError("target_type", "目标类型不能为空")
	}
	if t.TargetID == "" {
		return NewValidationError("target_id", "目标 ID 不能为空")
	}
	return nil
}

type TagFilter struct {
	Name       string
	TargetType string
	TargetID   string
	Keyword    string
}

func (f TagFilter) Match(t *Tag) bool {
	if f.Name != "" && t.Name != f.Name {
		return false
	}
	if f.TargetType != "" && t.TargetType != f.TargetType {
		return false
	}
	if f.TargetID != "" && t.TargetID != f.TargetID {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(t.Name), k) {
			return false
		}
	}
	return true
}
