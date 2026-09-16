package model

import (
	"strings"
	"time"
)

const (
	SessionRecording = "recording"
	SessionPaused    = "paused"
	SessionCompleted = "completed"
	SessionFailed    = "failed"
)

var sessionTransitions = map[string]map[string]bool{
	SessionRecording: {SessionPaused: true, SessionCompleted: true, SessionFailed: true},
	SessionPaused:    {SessionRecording: true, SessionCompleted: true, SessionFailed: true},
	SessionCompleted: {},
	SessionFailed:    {},
}

func SessionCanTransition(from, to string) bool {
	if m, ok := sessionTransitions[from]; ok {
		return m[to]
	}
	return false
}

type RecordingSession struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Source    string    `json:"source"`
	TargetHost string   `json:"target_host"`
	Protocol  string    `json:"protocol"`
	Status    string    `json:"status"`
	StartedAt time.Time `json:"started_at"`
	StoppedAt *time.Time `json:"stopped_at,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *RecordingSession) Validate() error {
	s.Name = strings.TrimSpace(s.Name)
	s.Source = strings.TrimSpace(s.Source)
	s.TargetHost = strings.TrimSpace(s.TargetHost)
	if s.Name == "" {
		return NewValidationError("name", "会话名称不能为空")
	}
	if s.Source == "" {
		return NewValidationError("source", "来源不能为空")
	}
	if s.TargetHost == "" {
		return NewValidationError("target_host", "目标主机不能为空")
	}
	if s.Protocol == "" {
		s.Protocol = "http"
	}
	if s.Protocol != "http" && s.Protocol != "tcp" {
		return NewValidationError("protocol", "协议只能是 http 或 tcp")
	}
	if s.Status == "" {
		s.Status = SessionRecording
	}
	if s.Status != SessionRecording && s.Status != SessionPaused && s.Status != SessionCompleted && s.Status != SessionFailed {
		return NewValidationError("status", "会话状态不合法")
	}
	return nil
}

type RecordingSessionFilter struct {
	Name     string
	Protocol string
	Status   string
	Keyword  string
}

func (f RecordingSessionFilter) Match(s *RecordingSession) bool {
	if f.Name != "" && s.Name != f.Name {
		return false
	}
	if f.Protocol != "" && s.Protocol != f.Protocol {
		return false
	}
	if f.Status != "" && s.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(s.Name), k) &&
			!strings.Contains(strings.ToLower(s.Source), k) {
			return false
		}
	}
	return true
}
