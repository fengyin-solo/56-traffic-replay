package model

import (
	"regexp"
	"strings"
	"time"
)

const (
	FilterFieldMethod = "method"
	FilterFieldPath   = "path"
	FilterFieldHeader = "header"
	FilterFieldBody   = "body"

	FilterOpEq       = "eq"
	FilterOpNe       = "ne"
	FilterOpContains = "contains"
	FilterOpRegex    = "regex"

	FilterActionInclude = "include"
	FilterActionExclude = "exclude"

	FilterStatusActive   = "active"
	FilterStatusDisabled = "disabled"
)

type FilterRule struct {
	ID        string    `json:"id"`
	SessionID string    `json:"session_id"`
	Field     string    `json:"field"`
	Operator  string    `json:"operator"`
	Value     string    `json:"value"`
	Action    string    `json:"action"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (f *FilterRule) Validate() error {
	f.Value = strings.TrimSpace(f.Value)
	if f.SessionID == "" {
		return NewValidationError("session_id", "会话 ID 不能为空")
	}
	if f.Field == "" {
		return NewValidationError("field", "字段不能为空")
	}
	if f.Field != FilterFieldMethod && f.Field != FilterFieldPath && f.Field != FilterFieldHeader && f.Field != FilterFieldBody {
		return NewValidationError("field", "字段类型不合法")
	}
	if f.Operator == "" {
		return NewValidationError("operator", "操作符不能为空")
	}
	if f.Operator != FilterOpEq && f.Operator != FilterOpNe && f.Operator != FilterOpContains && f.Operator != FilterOpRegex {
		return NewValidationError("operator", "操作符不合法")
	}
	if f.Action == "" {
		f.Action = FilterActionInclude
	}
	if f.Action != FilterActionInclude && f.Action != FilterActionExclude {
		return NewValidationError("action", "动作只能是 include 或 exclude")
	}
	if f.Status == "" {
		f.Status = FilterStatusActive
	}
	if f.Status != FilterStatusActive && f.Status != FilterStatusDisabled {
		return NewValidationError("status", "状态不合法")
	}
	return nil
}

func (f *FilterRule) MatchRecord(record *TrafficRecord) bool {
	if f.Status != FilterStatusActive {
		return true
	}
	var target string
	switch f.Field {
	case FilterFieldMethod:
		target = record.Method
	case FilterFieldPath:
		target = record.Path
	case FilterFieldHeader:
		target = record.Headers
	case FilterFieldBody:
		target = record.Body
	}

	matched := false
	switch f.Operator {
	case FilterOpEq:
		matched = target == f.Value
	case FilterOpNe:
		matched = target != f.Value
	case FilterOpContains:
		matched = strings.Contains(target, f.Value)
	case FilterOpRegex:
		if re, err := regexp.Compile(f.Value); err == nil {
			matched = re.MatchString(target)
		}
	}

	if f.Action == FilterActionInclude {
		return matched
	}
	return !matched
}

type FilterRuleFilter struct {
	SessionID string
	Field     string
	Action    string
	Status    string
}

func (ff FilterRuleFilter) Match(f *FilterRule) bool {
	if ff.SessionID != "" && f.SessionID != ff.SessionID {
		return false
	}
	if ff.Field != "" && f.Field != ff.Field {
		return false
	}
	if ff.Action != "" && f.Action != ff.Action {
		return false
	}
	if ff.Status != "" && f.Status != ff.Status {
		return false
	}
	return true
}
