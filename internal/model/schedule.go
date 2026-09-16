package model

import (
	"strings"
	"time"
)

const (
	ScheduleActive = "active"
	SchedulePaused = "paused"
)

type Schedule struct {
	ID         string     `json:"id"`
	ReplayTaskID string   `json:"replay_task_id"`
	CronExpr   string     `json:"cron_expr"`
	Status     string     `json:"status"`
	LastRunAt  *time.Time `json:"last_run_at,omitempty"`
	NextRunAt  *time.Time `json:"next_run_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

func (s *Schedule) Validate() error {
	s.CronExpr = strings.TrimSpace(s.CronExpr)
	if s.ReplayTaskID == "" {
		return NewValidationError("replay_task_id", "回放任务 ID 不能为空")
	}
	if s.CronExpr == "" {
		return NewValidationError("cron_expr", "Cron 表达式不能为空")
	}
	if s.Status == "" {
		s.Status = ScheduleActive
	}
	if s.Status != ScheduleActive && s.Status != SchedulePaused {
		return NewValidationError("status", "调度状态不合法")
	}
	return nil
}

type ScheduleFilter struct {
	ReplayTaskID string
	Status       string
}

func (f ScheduleFilter) Match(s *Schedule) bool {
	if f.ReplayTaskID != "" && s.ReplayTaskID != f.ReplayTaskID {
		return false
	}
	if f.Status != "" && s.Status != f.Status {
		return false
	}
	return true
}
