package service

import (
	"sort"
	"time"

	"replay/internal/model"
	"replay/pkg/idgen"
)

func (s *Service) CreateProjectConfig(input model.ProjectConfig) (*model.ProjectConfig, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	p := &model.ProjectConfig{
		ID:          idgen.Hex(),
		Key:         input.Key,
		Value:       input.Value,
		Description: input.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := s.store.CreateProjectConfig(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) ListProjectConfigs(filter model.ProjectConfigFilter, page, size int) ([]*model.ProjectConfig, int, error) {
	all := s.store.ListProjectConfigs()
	matched := make([]*model.ProjectConfig, 0, len(all))
	for _, p := range all {
		if filter.Match(p) {
			matched = append(matched, p)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.ProjectConfig{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetProjectConfig(id string) (*model.ProjectConfig, error) {
	return s.store.GetProjectConfig(id)
}

func (s *Service) UpdateProjectConfig(id string, input model.ProjectConfig) (*model.ProjectConfig, error) {
	p, err := s.store.GetProjectConfig(id)
	if err != nil {
		return nil, err
	}
	p.Key = input.Key
	p.Value = input.Value
	p.Description = input.Description
	p.UpdatedAt = time.Now()
	if err := p.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateProjectConfig(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) DeleteProjectConfig(id string) error {
	return s.store.DeleteProjectConfig(id)
}
