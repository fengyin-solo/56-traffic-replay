package service

import (
	"sort"
	"time"

	"replay/internal/model"
	"replay/pkg/idgen"
)

func (s *Service) CreateReplayResult(input model.ReplayResult) (*model.ReplayResult, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetReplayTask(input.ReplayTaskID); err != nil {
		return nil, model.NewValidationError("replay_task_id", "关联回放任务不存在")
	}
	r := &model.ReplayResult{
		ID:              idgen.Hex(),
		ReplayTaskID:    input.ReplayTaskID,
		TrafficRecordID: input.TrafficRecordID,
		Status:          input.Status,
		ResponseCode:    input.ResponseCode,
		DurationMs:      input.DurationMs,
		CreatedAt:       time.Now(),
	}
	if err := s.store.CreateReplayResult(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) ListReplayResults(filter model.ReplayResultFilter, page, size int) ([]*model.ReplayResult, int, error) {
	all := s.store.ListReplayResults()
	matched := make([]*model.ReplayResult, 0, len(all))
	for _, r := range all {
		if filter.Match(r) {
			matched = append(matched, r)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.ReplayResult{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetReplayResult(id string) (*model.ReplayResult, error) {
	return s.store.GetReplayResult(id)
}

func (s *Service) UpdateReplayResult(id string, input model.ReplayResult) (*model.ReplayResult, error) {
	r, err := s.store.GetReplayResult(id)
	if err != nil {
		return nil, err
	}
	r.Status = input.Status
	r.ResponseCode = input.ResponseCode
	r.DurationMs = input.DurationMs
	if err := r.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateReplayResult(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) DeleteReplayResult(id string) error {
	return s.store.DeleteReplayResult(id)
}
