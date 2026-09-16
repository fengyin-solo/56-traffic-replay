package service

import (
	"sort"
	"time"

	"replay/internal/model"
	"replay/pkg/idgen"
)

func (s *Service) CreateNotificationRule(input model.NotificationRule) (*model.NotificationRule, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	n := &model.NotificationRule{
		ID:        idgen.Hex(),
		Name:      input.Name,
		Channel:   input.Channel,
		Target:    input.Target,
		Events:    input.Events,
		Status:    input.Status,
		CreatedAt: time.Now(),
	}
	if err := s.store.CreateNotificationRule(n); err != nil {
		return nil, err
	}
	return n, nil
}

func (s *Service) ListNotificationRules(filter model.NotificationRuleFilter, page, size int) ([]*model.NotificationRule, int, error) {
	all := s.store.ListNotificationRules()
	matched := make([]*model.NotificationRule, 0, len(all))
	for _, n := range all {
		if filter.Match(n) {
			matched = append(matched, n)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.NotificationRule{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetNotificationRule(id string) (*model.NotificationRule, error) {
	return s.store.GetNotificationRule(id)
}

func (s *Service) UpdateNotificationRule(id string, input model.NotificationRule) (*model.NotificationRule, error) {
	n, err := s.store.GetNotificationRule(id)
	if err != nil {
		return nil, err
	}
	n.Name = input.Name
	n.Channel = input.Channel
	n.Target = input.Target
	n.Events = input.Events
	n.Status = input.Status
	if err := n.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateNotificationRule(n); err != nil {
		return nil, err
	}
	return n, nil
}

func (s *Service) DeleteNotificationRule(id string) error {
	return s.store.DeleteNotificationRule(id)
}
