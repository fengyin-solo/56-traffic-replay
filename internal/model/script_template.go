package model

import (
	"strings"
	"time"
)

const (
	ScriptTypePre  = "pre"
	ScriptTypePost = "post"
)

type ScriptTemplate struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	ScriptType  string    `json:"script_type"`
	Content     string    `json:"content"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func (s *ScriptTemplate) Validate() error {
	s.Name = strings.TrimSpace(s.Name)
	s.Content = strings.TrimSpace(s.Content)
	if s.Name == "" {
		return NewValidationError("name", "脚本模板名称不能为空")
	}
	if s.Content == "" {
		return NewValidationError("content", "脚本内容不能为空")
	}
	if s.ScriptType == "" {
		s.ScriptType = ScriptTypePre
	}
	if s.ScriptType != ScriptTypePre && s.ScriptType != ScriptTypePost {
		return NewValidationError("script_type", "脚本类型不合法")
	}
	return nil
}

type ScriptTemplateFilter struct {
	Name       string
	ScriptType string
	Keyword    string
}

func (f ScriptTemplateFilter) Match(s *ScriptTemplate) bool {
	if f.Name != "" && s.Name != f.Name {
		return false
	}
	if f.ScriptType != "" && s.ScriptType != f.ScriptType {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(s.Name), k) &&
			!strings.Contains(strings.ToLower(s.Description), k) {
			return false
		}
	}
	return true
}
