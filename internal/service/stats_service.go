package service

import (
	"replay/internal/model"
)

type OverviewStats struct {
	TotalSessions    int     `json:"total_sessions"`
	TotalTasks       int     `json:"total_tasks"`
	TotalRecords     int     `json:"total_records"`
	TotalEnvs        int     `json:"total_envs"`
	SuccessRate      float64 `json:"success_rate"`
	MatchRate        float64 `json:"match_rate"`
}

func (s *Service) GetOverviewStats() *OverviewStats {
	sessions := s.store.ListRecordingSessions()
	tasks := s.store.ListReplayTasks()
	records := s.store.ListTrafficRecords()
	envs := s.store.ListTargetEnvs()
	results := s.store.ListReplayResults()
	reports := s.store.ListCompareReports()

	var totalSuccess, totalResults int
	for _, r := range results {
		totalResults++
		if r.Status == model.ResultSuccess {
			totalSuccess++
		}
	}
	var successRate float64
	if totalResults > 0 {
		successRate = float64(totalSuccess) / float64(totalResults) * 100
	}

	var totalMatched, totalCompared int
	for _, rep := range reports {
		totalCompared += rep.TotalCompared
		totalMatched += rep.Matched
	}
	var matchRate float64
	if totalCompared > 0 {
		matchRate = float64(totalMatched) / float64(totalCompared) * 100
	}

	return &OverviewStats{
		TotalSessions: len(sessions),
		TotalTasks:    len(tasks),
		TotalRecords:  len(records),
		TotalEnvs:     len(envs),
		SuccessRate:   successRate,
		MatchRate:     matchRate,
	}
}

type MethodDistribution struct {
	Method string `json:"method"`
	Count  int    `json:"count"`
}

func (s *Service) GetMethodDistribution() []MethodDistribution {
	records := s.store.ListTrafficRecords()
	m := make(map[string]int)
	for _, r := range records {
		m[r.Method]++
	}
	res := make([]MethodDistribution, 0, len(m))
	for method, count := range m {
		res = append(res, MethodDistribution{Method: method, Count: count})
	}
	return res
}

type EnvDistribution struct {
	EnvID string `json:"env_id"`
	Count int    `json:"count"`
}

func (s *Service) GetEnvDistribution() []EnvDistribution {
	tasks := s.store.ListReplayTasks()
	m := make(map[string]int)
	for _, t := range tasks {
		m[t.TargetEnvID]++
	}
	res := make([]EnvDistribution, 0, len(m))
	for envID, count := range m {
		res = append(res, EnvDistribution{EnvID: envID, Count: count})
	}
	return res
}

type EndpointDistribution struct {
	Path  string `json:"path"`
	Count int    `json:"count"`
}

func (s *Service) GetEndpointDistribution() []EndpointDistribution {
	records := s.store.ListTrafficRecords()
	m := make(map[string]int)
	for _, r := range records {
		m[r.Path]++
	}
	res := make([]EndpointDistribution, 0, len(m))
	for path, count := range m {
		res = append(res, EndpointDistribution{Path: path, Count: count})
	}
	return res
}

type SuccessRateTrend struct {
	TaskID      string  `json:"task_id"`
	SuccessRate float64 `json:"success_rate"`
}

func (s *Service) GetSuccessRateTrend() []SuccessRateTrend {
	tasks := s.store.ListReplayTasks()
	res := make([]SuccessRateTrend, 0, len(tasks))
	for _, t := range tasks {
		if t.TotalRequests > 0 {
			res = append(res, SuccessRateTrend{
				TaskID:      t.ID,
				SuccessRate: float64(t.SuccessCount) / float64(t.TotalRequests) * 100,
			})
		}
	}
	return res
}

func (s *Service) GetSnapshotExport() *model.SnapshotExport {
	return &model.SnapshotExport{
		RecordingSessions: s.store.ListRecordingSessions(),
		TrafficRecords:    s.store.ListTrafficRecords(),
		ReplayTasks:       s.store.ListReplayTasks(),
		TargetEnvs:        s.store.ListTargetEnvs(),
		FilterRules:       s.store.ListFilterRules(),
		CompareReports:    s.store.ListCompareReports(),
		DiffRecords:       s.store.ListDiffRecords(),
		Schedules:         s.store.ListSchedules(),
		ReplayResults:     s.store.ListReplayResults(),
		AuditLogs:         s.store.ListAuditLogs(),
		WebhookConfigs:    s.store.ListWebhookConfigs(),
		ScriptTemplates:   s.store.ListScriptTemplates(),
		EnvVariables:      s.store.ListEnvVariables(),
		NotificationRules: s.store.ListNotificationRules(),
		ProjectConfigs:    s.store.ListProjectConfigs(),
		Tags:              s.store.ListTags(),
		Certificates:      s.store.ListCertificates(),
		DataSources:       s.store.ListDataSources(),
	}
}
