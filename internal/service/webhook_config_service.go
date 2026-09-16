package service

import (
	"sort"
	"time"

	"replay/internal/model"
	"replay/pkg/idgen"
)

func (s *Service) CreateWebhookConfig(input model.WebhookConfig) (*model.WebhookConfig, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	w := &model.WebhookConfig{
		ID:        idgen.Hex(),
		Name:      input.Name,
		URL:       input.URL,
		Secret:    input.Secret,
		Events:    input.Events,
		Status:    input.Status,
		CreatedAt: time.Now(),
	}
	if err := s.store.CreateWebhookConfig(w); err != nil {
		return nil, err
	}
	return w, nil
}

func (s *Service) ListWebhookConfigs(filter model.WebhookConfigFilter, page, size int) ([]*model.WebhookConfig, int, error) {
	all := s.store.ListWebhookConfigs()
	matched := make([]*model.WebhookConfig, 0, len(all))
	for _, w := range all {
		if filter.Match(w) {
			matched = append(matched, w)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.WebhookConfig{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetWebhookConfig(id string) (*model.WebhookConfig, error) {
	return s.store.GetWebhookConfig(id)
}

func (s *Service) UpdateWebhookConfig(id string, input model.WebhookConfig) (*model.WebhookConfig, error) {
	w, err := s.store.GetWebhookConfig(id)
	if err != nil {
		return nil, err
	}
	w.Name = input.Name
	w.URL = input.URL
	w.Secret = input.Secret
	w.Events = input.Events
	w.Status = input.Status
	if err := w.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateWebhookConfig(w); err != nil {
		return nil, err
	}
	return w, nil
}

func (s *Service) DeleteWebhookConfig(id string) error {
	return s.store.DeleteWebhookConfig(id)
}
