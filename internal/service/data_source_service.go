package service

import (
	"sort"
	"time"

	"replay/internal/model"
	"replay/pkg/idgen"
)

func (s *Service) CreateDataSource(input model.DataSource) (*model.DataSource, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	d := &model.DataSource{
		ID:        idgen.Hex(),
		Name:      input.Name,
		Type:      input.Type,
		Config:    input.Config,
		Status:    input.Status,
		CreatedAt: time.Now(),
	}
	if err := s.store.CreateDataSource(d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) ListDataSources(filter model.DataSourceFilter, page, size int) ([]*model.DataSource, int, error) {
	all := s.store.ListDataSources()
	matched := make([]*model.DataSource, 0, len(all))
	for _, d := range all {
		if filter.Match(d) {
			matched = append(matched, d)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.DataSource{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetDataSource(id string) (*model.DataSource, error) {
	return s.store.GetDataSource(id)
}

func (s *Service) UpdateDataSource(id string, input model.DataSource) (*model.DataSource, error) {
	d, err := s.store.GetDataSource(id)
	if err != nil {
		return nil, err
	}
	d.Name = input.Name
	d.Type = input.Type
	d.Config = input.Config
	d.Status = input.Status
	if err := d.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateDataSource(d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) DeleteDataSource(id string) error {
	return s.store.DeleteDataSource(id)
}
