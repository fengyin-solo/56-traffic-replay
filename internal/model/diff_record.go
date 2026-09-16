package model

import (
	"time"
)

type DiffRecord struct {
	ID               string    `json:"id"`
	ReportID         string    `json:"report_id"`
	TrafficRecordID  string    `json:"traffic_record_id"`
	OriginalResponse string    `json:"original_response"`
	ReplayedResponse string    `json:"replayed_response"`
	DiffDetail       string    `json:"diff_detail"`
	CreatedAt        time.Time `json:"created_at"`
}

func (d *DiffRecord) Validate() error {
	if d.ReportID == "" {
		return NewValidationError("report_id", "报告 ID 不能为空")
	}
	if d.TrafficRecordID == "" {
		return NewValidationError("traffic_record_id", "流量记录 ID 不能为空")
	}
	return nil
}

type DiffRecordFilter struct {
	ReportID        string
	TrafficRecordID string
}

func (f DiffRecordFilter) Match(d *DiffRecord) bool {
	if f.ReportID != "" && d.ReportID != f.ReportID {
		return false
	}
	if f.TrafficRecordID != "" && d.TrafficRecordID != f.TrafficRecordID {
		return false
	}
	return true
}
