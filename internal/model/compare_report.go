package model

import (
	"time"
)

const (
	ReportPending   = "pending"
	ReportRunning   = "running"
	ReportCompleted = "completed"
)

var reportTransitions = map[string]map[string]bool{
	ReportPending:   {ReportRunning: true},
	ReportRunning:   {ReportCompleted: true},
	ReportCompleted: {},
}

func ReportCanTransition(from, to string) bool {
	if m, ok := reportTransitions[from]; ok {
		return m[to]
	}
	return false
}

type CompareReport struct {
	ID             string     `json:"id"`
	ReplayTaskID   string     `json:"replay_task_id"`
	TotalCompared  int        `json:"total_compared"`
	Matched        int        `json:"matched"`
	Mismatched     int        `json:"mismatched"`
	DiffCount      int        `json:"diff_count"`
	Status         string     `json:"status"`
	CreatedAt      time.Time  `json:"created_at"`
	FinishedAt     *time.Time `json:"finished_at,omitempty"`
}

func (r *CompareReport) Validate() error {
	if r.ReplayTaskID == "" {
		return NewValidationError("replay_task_id", "回放任务 ID 不能为空")
	}
	if r.Status == "" {
		r.Status = ReportPending
	}
	if r.Status != ReportPending && r.Status != ReportRunning && r.Status != ReportCompleted {
		return NewValidationError("status", "报告状态不合法")
	}
	return nil
}

type CompareReportFilter struct {
	ReplayTaskID string
	Status       string
}

func (f CompareReportFilter) Match(r *CompareReport) bool {
	if f.ReplayTaskID != "" && r.ReplayTaskID != f.ReplayTaskID {
		return false
	}
	if f.Status != "" && r.Status != f.Status {
		return false
	}
	return true
}
