package service

import (
	"sort"
	"time"

	"replay/internal/model"
	"replay/pkg/idgen"
)

func (s *Service) CreateReplayTask(input model.ReplayTask) (*model.ReplayTask, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetRecordingSession(input.SessionID); err != nil {
		return nil, model.NewValidationError("session_id", "关联会话不存在")
	}
	if _, err := s.store.GetTargetEnv(input.TargetEnvID); err != nil {
		return nil, model.NewValidationError("target_env_id", "目标环境不存在")
	}
	task := &model.ReplayTask{
		ID:            idgen.Hex(),
		SessionID:     input.SessionID,
		TargetEnvID:   input.TargetEnvID,
		Status:        model.TaskPending,
		TotalRequests: 0,
		SuccessCount:  0,
		FailedCount:   0,
		CreatedAt:     time.Now(),
	}
	if err := s.store.CreateReplayTask(task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *Service) ListReplayTasks(filter model.ReplayTaskFilter, page, size int) ([]*model.ReplayTask, int, error) {
	all := s.store.ListReplayTasks()
	matched := make([]*model.ReplayTask, 0, len(all))
	for _, t := range all {
		if filter.Match(t) {
			matched = append(matched, t)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.ReplayTask{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetReplayTask(id string) (*model.ReplayTask, error) {
	return s.store.GetReplayTask(id)
}

func (s *Service) UpdateReplayTask(id string, input model.ReplayTask) (*model.ReplayTask, error) {
	task, err := s.store.GetReplayTask(id)
	if err != nil {
		return nil, err
	}
	task.TargetEnvID = input.TargetEnvID
	if err := task.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateReplayTask(task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *Service) DeleteReplayTask(id string) error {
	return s.store.DeleteReplayTask(id)
}

func (s *Service) RunReplayTask(id string) (*model.ReplayTask, error) {
	task, err := s.store.GetReplayTask(id)
	if err != nil {
		return nil, err
	}
	if !model.TaskCanTransition(task.Status, model.TaskRunning) {
		return nil, model.NewValidationError("status", "当前状态不允许启动回放")
	}

	records := s.store.ListTrafficRecords()
	var matched []*model.TrafficRecord
	for _, r := range records {
		if r.SessionID == task.SessionID {
			matched = append(matched, r)
		}
	}

	rules := s.store.ListFilterRules()
	var activeRules []*model.FilterRule
	for _, rule := range rules {
		if rule.SessionID == task.SessionID && rule.Status == model.FilterStatusActive {
			activeRules = append(activeRules, rule)
		}
	}

	filtered := make([]*model.TrafficRecord, 0, len(matched))
	for _, r := range matched {
		pass := true
		for _, rule := range activeRules {
			if !rule.MatchRecord(r) {
				pass = false
				break
			}
		}
		if pass {
			filtered = append(filtered, r)
		}
	}

	task.Status = model.TaskRunning
	now := time.Now()
	task.StartedAt = &now
	task.TotalRequests = len(filtered)
	task.SuccessCount = 0
	task.FailedCount = 0

	for _, r := range filtered {
		status := model.ResultSuccess
		code := r.StatusCode
		duration := r.DurationMs
		if code >= 500 {
			status = model.ResultFailed
			task.FailedCount++
		} else {
			task.SuccessCount++
		}
		result := &model.ReplayResult{
			ID:              idgen.Hex(),
			ReplayTaskID:    task.ID,
			TrafficRecordID: r.ID,
			Status:          status,
			ResponseCode:    code,
			DurationMs:      duration,
			CreatedAt:       time.Now(),
		}
		_ = s.store.CreateReplayResult(result)
	}

	if err := s.store.UpdateReplayTask(task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *Service) CompleteReplayTask(id string) (*model.ReplayTask, error) {
	task, err := s.store.GetReplayTask(id)
	if err != nil {
		return nil, err
	}
	if !model.TaskCanTransition(task.Status, model.TaskCompleted) {
		return nil, model.NewValidationError("status", "当前状态不允许完成任务")
	}
	task.Status = model.TaskCompleted
	now := time.Now()
	task.FinishedAt = &now
	if err := s.store.UpdateReplayTask(task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *Service) FailReplayTask(id string) (*model.ReplayTask, error) {
	task, err := s.store.GetReplayTask(id)
	if err != nil {
		return nil, err
	}
	if !model.TaskCanTransition(task.Status, model.TaskFailed) {
		return nil, model.NewValidationError("status", "当前状态不允许标记失败")
	}
	task.Status = model.TaskFailed
	now := time.Now()
	task.FinishedAt = &now
	if err := s.store.UpdateReplayTask(task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *Service) BatchCreateReplayTasks(inputs []model.ReplayTask) ([]*model.ReplayTask, error) {
	results := make([]*model.ReplayTask, 0, len(inputs))
	for _, input := range inputs {
		if task, err := s.CreateReplayTask(input); err == nil {
			results = append(results, task)
		}
	}
	return results, nil
}
