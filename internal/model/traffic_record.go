package model

import (
	"strings"
	"time"
)

type TrafficRecord struct {
	ID         string    `json:"id"`
	SessionID  string    `json:"session_id"`
	Method     string    `json:"method"`
	Path       string    `json:"path"`
	Headers    string    `json:"headers"`
	Body       string    `json:"body"`
	Timestamp  time.Time `json:"timestamp"`
	DurationMs int       `json:"duration_ms"`
	StatusCode int       `json:"status_code"`
	CreatedAt  time.Time `json:"created_at"`
}

func (t *TrafficRecord) Validate() error {
	t.Method = strings.TrimSpace(t.Method)
	t.Path = strings.TrimSpace(t.Path)
	if t.SessionID == "" {
		return NewValidationError("session_id", "会话 ID 不能为空")
	}
	if t.Method == "" {
		return NewValidationError("method", "请求方法不能为空")
	}
	if t.Path == "" {
		return NewValidationError("path", "请求路径不能为空")
	}
	return nil
}

type TrafficRecordFilter struct {
	SessionID string
	Method    string
	Path      string
	Keyword   string
}

func (f TrafficRecordFilter) Match(t *TrafficRecord) bool {
	if f.SessionID != "" && t.SessionID != f.SessionID {
		return false
	}
	if f.Method != "" && t.Method != f.Method {
		return false
	}
	if f.Path != "" && t.Path != f.Path {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(t.Path), k) &&
			!strings.Contains(strings.ToLower(t.Method), k) &&
			!strings.Contains(strings.ToLower(t.Body), k) {
			return false
		}
	}
	return true
}
