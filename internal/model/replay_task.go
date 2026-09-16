package model

import (
	"strings"
	"time"
)

const (
	TaskPending    = "pending"
	TaskRunning    = "running"
	TaskCompleted  = "completed"
	TaskFailed     = "failed"
)

var taskTransitions = map[string]map[string]bool{
	TaskPending:   {TaskRunning: true, TaskFailed: true},
	TaskRunning:   {TaskCompleted: true, TaskFailed: true},
	TaskCompleted: {},
	TaskFailed:    {},
}

func TaskCanTransition(from, to string) bool {
	if m, ok := taskTransitions[from]; ok {
		return m[to]
	}
	return false
}

type ReplayTask struct {
	ID           string     `json:"id"`
	SessionID    string     `json:"session_id"`
	TargetEnvID  string     `json:"target_env_id"`
	Status       string     `json:"status"`
	TotalRequests int       `json:"total_requests"`
	SuccessCount int       `json:"success_count"`
	FailedCount  int       `json:"failed_count"`
	CreatedAt    time.Time  `json:"created_at"`
	StartedAt    *time.Time `json:"started_at,omitempty"`
	FinishedAt   *time.Time `json:"finished_at,omitempty"`
}

func (t *ReplayTask) Validate() error {
	if t.SessionID == "" {
		return NewValidationError("session_id", "会话 ID 不能为空")
	}
	if t.TargetEnvID == "" {
		return NewValidationError("target_env_id", "目标环境 ID 不能为空")
	}
	if t.Status == "" {
		t.Status = TaskPending
	}
	if t.Status != TaskPending && t.Status != TaskRunning && t.Status != TaskCompleted && t.Status != TaskFailed {
		return NewValidationError("status", "任务状态不合法")
	}
	return nil
}

type ReplayTaskFilter struct {
	SessionID   string
	TargetEnvID string
	Status      string
	Keyword     string
}

func (f ReplayTaskFilter) Match(t *ReplayTask) bool {
	if f.SessionID != "" && t.SessionID != f.SessionID {
		return false
	}
	if f.TargetEnvID != "" && t.TargetEnvID != f.TargetEnvID {
		return false
	}
	if f.Status != "" && t.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(t.SessionID), k) {
			return false
		}
	}
	return true
}
