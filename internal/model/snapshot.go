package model

// SnapshotExport 用于全量快照导出。
type SnapshotExport struct {
	RecordingSessions []*RecordingSession `json:"recording_sessions"`
	TrafficRecords    []*TrafficRecord    `json:"traffic_records"`
	ReplayTasks       []*ReplayTask       `json:"replay_tasks"`
	TargetEnvs        []*TargetEnv        `json:"target_envs"`
	FilterRules       []*FilterRule       `json:"filter_rules"`
	CompareReports    []*CompareReport    `json:"compare_reports"`
	DiffRecords       []*DiffRecord       `json:"diff_records"`
	Schedules         []*Schedule         `json:"schedules"`
	ReplayResults     []*ReplayResult     `json:"replay_results"`
	AuditLogs         []*AuditLog         `json:"audit_logs"`
	WebhookConfigs    []*WebhookConfig    `json:"webhook_configs"`
	ScriptTemplates   []*ScriptTemplate   `json:"script_templates"`
	EnvVariables      []*EnvVariable      `json:"env_variables"`
	NotificationRules []*NotificationRule `json:"notification_rules"`
	ProjectConfigs    []*ProjectConfig    `json:"project_configs"`
	Tags              []*Tag              `json:"tags"`
	Certificates      []*Certificate      `json:"certificates"`
	DataSources       []*DataSource       `json:"data_sources"`
}
