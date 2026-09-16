// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"replay/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

type Store interface {
	// RecordingSession
	CreateRecordingSession(s *model.RecordingSession) error
	GetRecordingSession(id string) (*model.RecordingSession, error)
	GetRecordingSessionByName(name string) (*model.RecordingSession, error)
	ListRecordingSessions() []*model.RecordingSession
	UpdateRecordingSession(s *model.RecordingSession) error
	DeleteRecordingSession(id string) error

	// TrafficRecord
	CreateTrafficRecord(t *model.TrafficRecord) error
	GetTrafficRecord(id string) (*model.TrafficRecord, error)
	ListTrafficRecords() []*model.TrafficRecord
	UpdateTrafficRecord(t *model.TrafficRecord) error
	DeleteTrafficRecord(id string) error

	// ReplayTask
	CreateReplayTask(t *model.ReplayTask) error
	GetReplayTask(id string) (*model.ReplayTask, error)
	ListReplayTasks() []*model.ReplayTask
	UpdateReplayTask(t *model.ReplayTask) error
	DeleteReplayTask(id string) error

	// TargetEnv
	CreateTargetEnv(e *model.TargetEnv) error
	GetTargetEnv(id string) (*model.TargetEnv, error)
	GetTargetEnvByName(name string) (*model.TargetEnv, error)
	ListTargetEnvs() []*model.TargetEnv
	UpdateTargetEnv(e *model.TargetEnv) error
	DeleteTargetEnv(id string) error

	// FilterRule
	CreateFilterRule(f *model.FilterRule) error
	GetFilterRule(id string) (*model.FilterRule, error)
	ListFilterRules() []*model.FilterRule
	UpdateFilterRule(f *model.FilterRule) error
	DeleteFilterRule(id string) error

	// CompareReport
	CreateCompareReport(r *model.CompareReport) error
	GetCompareReport(id string) (*model.CompareReport, error)
	ListCompareReports() []*model.CompareReport
	UpdateCompareReport(r *model.CompareReport) error
	DeleteCompareReport(id string) error

	// DiffRecord
	CreateDiffRecord(d *model.DiffRecord) error
	GetDiffRecord(id string) (*model.DiffRecord, error)
	ListDiffRecords() []*model.DiffRecord
	UpdateDiffRecord(d *model.DiffRecord) error
	DeleteDiffRecord(id string) error

	// Schedule
	CreateSchedule(s *model.Schedule) error
	GetSchedule(id string) (*model.Schedule, error)
	ListSchedules() []*model.Schedule
	UpdateSchedule(s *model.Schedule) error
	DeleteSchedule(id string) error

	// ReplayResult
	CreateReplayResult(r *model.ReplayResult) error
	GetReplayResult(id string) (*model.ReplayResult, error)
	ListReplayResults() []*model.ReplayResult
	UpdateReplayResult(r *model.ReplayResult) error
	DeleteReplayResult(id string) error

	// AuditLog
	CreateAuditLog(a *model.AuditLog) error
	GetAuditLog(id string) (*model.AuditLog, error)
	ListAuditLogs() []*model.AuditLog
	UpdateAuditLog(a *model.AuditLog) error
	DeleteAuditLog(id string) error

	// WebhookConfig
	CreateWebhookConfig(w *model.WebhookConfig) error
	GetWebhookConfig(id string) (*model.WebhookConfig, error)
	GetWebhookConfigByName(name string) (*model.WebhookConfig, error)
	ListWebhookConfigs() []*model.WebhookConfig
	UpdateWebhookConfig(w *model.WebhookConfig) error
	DeleteWebhookConfig(id string) error

	// ScriptTemplate
	CreateScriptTemplate(st *model.ScriptTemplate) error
	GetScriptTemplate(id string) (*model.ScriptTemplate, error)
	ListScriptTemplates() []*model.ScriptTemplate
	UpdateScriptTemplate(st *model.ScriptTemplate) error
	DeleteScriptTemplate(id string) error

	// EnvVariable
	CreateEnvVariable(e *model.EnvVariable) error
	GetEnvVariable(id string) (*model.EnvVariable, error)
	ListEnvVariables() []*model.EnvVariable
	UpdateEnvVariable(e *model.EnvVariable) error
	DeleteEnvVariable(id string) error

	// NotificationRule
	CreateNotificationRule(n *model.NotificationRule) error
	GetNotificationRule(id string) (*model.NotificationRule, error)
	ListNotificationRules() []*model.NotificationRule
	UpdateNotificationRule(n *model.NotificationRule) error
	DeleteNotificationRule(id string) error

	// ProjectConfig
	CreateProjectConfig(p *model.ProjectConfig) error
	GetProjectConfig(id string) (*model.ProjectConfig, error)
	GetProjectConfigByKey(key string) (*model.ProjectConfig, error)
	ListProjectConfigs() []*model.ProjectConfig
	UpdateProjectConfig(p *model.ProjectConfig) error
	DeleteProjectConfig(id string) error

	// Tag
	CreateTag(t *model.Tag) error
	GetTag(id string) (*model.Tag, error)
	ListTags() []*model.Tag
	UpdateTag(t *model.Tag) error
	DeleteTag(id string) error

	// Certificate
	CreateCertificate(c *model.Certificate) error
	GetCertificate(id string) (*model.Certificate, error)
	ListCertificates() []*model.Certificate
	UpdateCertificate(c *model.Certificate) error
	DeleteCertificate(id string) error

	// DataSource
	CreateDataSource(d *model.DataSource) error
	GetDataSource(id string) (*model.DataSource, error)
	ListDataSources() []*model.DataSource
	UpdateDataSource(d *model.DataSource) error
	DeleteDataSource(id string) error
}
