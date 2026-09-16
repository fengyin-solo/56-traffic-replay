package store

import (
	"sync"

	"replay/internal/model"
)

type MemoryStore struct {
	mu                sync.RWMutex
	sessions          map[string]*model.RecordingSession
	trafficRecords    map[string]*model.TrafficRecord
	replayTasks       map[string]*model.ReplayTask
	targetEnvs        map[string]*model.TargetEnv
	filterRules       map[string]*model.FilterRule
	compareReports    map[string]*model.CompareReport
	diffRecords       map[string]*model.DiffRecord
	schedules         map[string]*model.Schedule
	replayResults     map[string]*model.ReplayResult
	auditLogs         map[string]*model.AuditLog
	webhookConfigs    map[string]*model.WebhookConfig
	scriptTemplates   map[string]*model.ScriptTemplate
	envVariables      map[string]*model.EnvVariable
	notificationRules map[string]*model.NotificationRule
	projectConfigs    map[string]*model.ProjectConfig
	tags              map[string]*model.Tag
	certificates      map[string]*model.Certificate
	dataSources       map[string]*model.DataSource
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		sessions:          make(map[string]*model.RecordingSession),
		trafficRecords:    make(map[string]*model.TrafficRecord),
		replayTasks:       make(map[string]*model.ReplayTask),
		targetEnvs:        make(map[string]*model.TargetEnv),
		filterRules:       make(map[string]*model.FilterRule),
		compareReports:    make(map[string]*model.CompareReport),
		diffRecords:       make(map[string]*model.DiffRecord),
		schedules:         make(map[string]*model.Schedule),
		replayResults:     make(map[string]*model.ReplayResult),
		auditLogs:         make(map[string]*model.AuditLog),
		webhookConfigs:    make(map[string]*model.WebhookConfig),
		scriptTemplates:   make(map[string]*model.ScriptTemplate),
		envVariables:      make(map[string]*model.EnvVariable),
		notificationRules: make(map[string]*model.NotificationRule),
		projectConfigs:    make(map[string]*model.ProjectConfig),
		tags:              make(map[string]*model.Tag),
		certificates:      make(map[string]*model.Certificate),
		dataSources:       make(map[string]*model.DataSource),
	}
}

var _ Store = (*MemoryStore)(nil)
