package model

import (
	"time"
)

const (
	ResultSuccess   = "success"
	ResultFailed    = "failed"
	ResultMismatch  = "mismatch"
)

type ReplayResult struct {
	ID              string    `json:"id"`
	ReplayTaskID    string    `json:"replay_task_id"`
	TrafficRecordID string    `json:"traffic_record_id"`
	Status          string    `json:"status"`
	ResponseCode    int       `json:"response_code"`
	DurationMs      int       `json:"duration_ms"`
	CreatedAt       time.Time `json:"created_at"`
}

func (r *ReplayResult) Validate() error {
	if r.ReplayTaskID == "" {
		return NewValidationError("replay_task_id", "回放任务 ID 不能为空")
	}
	if r.TrafficRecordID == "" {
		return NewValidationError("traffic_record_id", "流量记录 ID 不能为空")
	}
	if r.Status == "" {
		return NewValidationError("status", "结果状态不能为空")
	}
	if r.Status != ResultSuccess && r.Status != ResultFailed && r.Status != ResultMismatch {
		return NewValidationError("status", "结果状态不合法")
	}
	return nil
}

type ReplayResultFilter struct {
	ReplayTaskID    string
	TrafficRecordID string
	Status          string
}

func (f ReplayResultFilter) Match(r *ReplayResult) bool {
	if f.ReplayTaskID != "" && r.ReplayTaskID != f.ReplayTaskID {
		return false
	}
	if f.TrafficRecordID != "" && r.TrafficRecordID != f.TrafficRecordID {
		return false
	}
	if f.Status != "" && r.Status != f.Status {
		return false
	}
	return true
}
