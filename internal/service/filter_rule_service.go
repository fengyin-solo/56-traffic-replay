package service

import (
	"sort"
	"time"

	"replay/internal/model"
	"replay/pkg/idgen"
)

func (s *Service) CreateFilterRule(input model.FilterRule) (*model.FilterRule, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetRecordingSession(input.SessionID); err != nil {
		return nil, model.NewValidationError("session_id", "关联会话不存在")
	}
	rule := &model.FilterRule{
		ID:        idgen.Hex(),
		SessionID: input.SessionID,
		Field:     input.Field,
		Operator:  input.Operator,
		Value:     input.Value,
		Action:    input.Action,
		Status:    input.Status,
		CreatedAt: time.Now(),
	}
	if err := s.store.CreateFilterRule(rule); err != nil {
		return nil, err
	}
	return rule, nil
}

func (s *Service) ListFilterRules(filter model.FilterRuleFilter, page, size int) ([]*model.FilterRule, int, error) {
	all := s.store.ListFilterRules()
	matched := make([]*model.FilterRule, 0, len(all))
	for _, f := range all {
		if filter.Match(f) {
			matched = append(matched, f)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.FilterRule{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetFilterRule(id string) (*model.FilterRule, error) {
	return s.store.GetFilterRule(id)
}

func (s *Service) UpdateFilterRule(id string, input model.FilterRule) (*model.FilterRule, error) {
	rule, err := s.store.GetFilterRule(id)
	if err != nil {
		return nil, err
	}
	rule.Field = input.Field
	rule.Operator = input.Operator
	rule.Value = input.Value
	rule.Action = input.Action
	rule.Status = input.Status
	if err := rule.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateFilterRule(rule); err != nil {
		return nil, err
	}
	return rule, nil
}

func (s *Service) DeleteFilterRule(id string) error {
	return s.store.DeleteFilterRule(id)
}
