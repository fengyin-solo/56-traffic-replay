package model

import (
	"strings"
	"time"
)

type EnvVariable struct {
	ID        string    `json:"id"`
	EnvID     string    `json:"env_id"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	IsSecret  bool      `json:"is_secret"`
	CreatedAt time.Time `json:"created_at"`
}

func (e *EnvVariable) Validate() error {
	e.Key = strings.TrimSpace(e.Key)
	if e.EnvID == "" {
		return NewValidationError("env_id", "环境 ID 不能为空")
	}
	if e.Key == "" {
		return NewValidationError("key", "变量名不能为空")
	}
	return nil
}

type EnvVariableFilter struct {
	EnvID string
	Key   string
}

func (f EnvVariableFilter) Match(e *EnvVariable) bool {
	if f.EnvID != "" && e.EnvID != f.EnvID {
		return false
	}
	if f.Key != "" && e.Key != f.Key {
		return false
	}
	return true
}
