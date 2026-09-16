package service

import (
	"sort"
	"time"

	"replay/internal/model"
	"replay/pkg/idgen"
)

func (s *Service) CreateTargetEnv(input model.TargetEnv) (*model.TargetEnv, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	env := &model.TargetEnv{
		ID:        idgen.Hex(),
		Name:      input.Name,
		BaseURL:   input.BaseURL,
		AuthToken: input.AuthToken,
		Status:    input.Status,
		CreatedAt: time.Now(),
	}
	if err := s.store.CreateTargetEnv(env); err != nil {
		return nil, err
	}
	return env, nil
}

func (s *Service) ListTargetEnvs(filter model.TargetEnvFilter, page, size int) ([]*model.TargetEnv, int, error) {
	all := s.store.ListTargetEnvs()
	matched := make([]*model.TargetEnv, 0, len(all))
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
		return []*model.TargetEnv{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetTargetEnv(id string) (*model.TargetEnv, error) {
	return s.store.GetTargetEnv(id)
}

func (s *Service) UpdateTargetEnv(id string, input model.TargetEnv) (*model.TargetEnv, error) {
	env, err := s.store.GetTargetEnv(id)
	if err != nil {
		return nil, err
	}
	env.Name = input.Name
	env.BaseURL = input.BaseURL
	env.AuthToken = input.AuthToken
	env.Status = input.Status
	if err := env.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateTargetEnv(env); err != nil {
		return nil, err
	}
	return env, nil
}

func (s *Service) DeleteTargetEnv(id string) error {
	return s.store.DeleteTargetEnv(id)
}
