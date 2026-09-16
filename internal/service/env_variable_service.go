package service

import (
	"sort"
	"time"

	"replay/internal/model"
	"replay/pkg/idgen"
)

func (s *Service) CreateEnvVariable(input model.EnvVariable) (*model.EnvVariable, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetTargetEnv(input.EnvID); err != nil {
		return nil, model.NewValidationError("env_id", "目标环境不存在")
	}
	e := &model.EnvVariable{
		ID:        idgen.Hex(),
		EnvID:     input.EnvID,
		Key:       input.Key,
		Value:     input.Value,
		IsSecret:  input.IsSecret,
		CreatedAt: time.Now(),
	}
	if err := s.store.CreateEnvVariable(e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *Service) ListEnvVariables(filter model.EnvVariableFilter, page, size int) ([]*model.EnvVariable, int, error) {
	all := s.store.ListEnvVariables()
	matched := make([]*model.EnvVariable, 0, len(all))
	for _, e := range all {
		if filter.Match(e) {
			matched = append(matched, e)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.EnvVariable{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetEnvVariable(id string) (*model.EnvVariable, error) {
	return s.store.GetEnvVariable(id)
}

func (s *Service) UpdateEnvVariable(id string, input model.EnvVariable) (*model.EnvVariable, error) {
	e, err := s.store.GetEnvVariable(id)
	if err != nil {
		return nil, err
	}
	e.Key = input.Key
	e.Value = input.Value
	e.IsSecret = input.IsSecret
	if err := e.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateEnvVariable(e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *Service) DeleteEnvVariable(id string) error {
	return s.store.DeleteEnvVariable(id)
}
