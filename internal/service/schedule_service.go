package service

import (
	"sort"
	"time"

	"replay/internal/model"
	"replay/pkg/idgen"
)

func (s *Service) CreateSchedule(input model.Schedule) (*model.Schedule, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetReplayTask(input.ReplayTaskID); err != nil {
		return nil, model.NewValidationError("replay_task_id", "关联回放任务不存在")
	}
	sch := &model.Schedule{
		ID:           idgen.Hex(),
		ReplayTaskID: input.ReplayTaskID,
		CronExpr:     input.CronExpr,
		Status:       input.Status,
		CreatedAt:    time.Now(),
	}
	if err := s.store.CreateSchedule(sch); err != nil {
		return nil, err
	}
	return sch, nil
}

func (s *Service) ListSchedules(filter model.ScheduleFilter, page, size int) ([]*model.Schedule, int, error) {
	all := s.store.ListSchedules()
	matched := make([]*model.Schedule, 0, len(all))
	for _, sch := range all {
		if filter.Match(sch) {
			matched = append(matched, sch)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Schedule{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetSchedule(id string) (*model.Schedule, error) {
	return s.store.GetSchedule(id)
}

func (s *Service) UpdateSchedule(id string, input model.Schedule) (*model.Schedule, error) {
	sch, err := s.store.GetSchedule(id)
	if err != nil {
		return nil, err
	}
	sch.CronExpr = input.CronExpr
	sch.Status = input.Status
	if err := sch.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateSchedule(sch); err != nil {
		return nil, err
	}
	return sch, nil
}

func (s *Service) DeleteSchedule(id string) error {
	return s.store.DeleteSchedule(id)
}
