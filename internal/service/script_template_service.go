package service

import (
	"sort"
	"time"

	"replay/internal/model"
	"replay/pkg/idgen"
)

func (s *Service) CreateScriptTemplate(input model.ScriptTemplate) (*model.ScriptTemplate, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	st := &model.ScriptTemplate{
		ID:          idgen.Hex(),
		Name:        input.Name,
		ScriptType:  input.ScriptType,
		Content:     input.Content,
		Description: input.Description,
		CreatedAt:   time.Now(),
	}
	if err := s.store.CreateScriptTemplate(st); err != nil {
		return nil, err
	}
	return st, nil
}

func (s *Service) ListScriptTemplates(filter model.ScriptTemplateFilter, page, size int) ([]*model.ScriptTemplate, int, error) {
	all := s.store.ListScriptTemplates()
	matched := make([]*model.ScriptTemplate, 0, len(all))
	for _, st := range all {
		if filter.Match(st) {
			matched = append(matched, st)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.ScriptTemplate{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetScriptTemplate(id string) (*model.ScriptTemplate, error) {
	return s.store.GetScriptTemplate(id)
}

func (s *Service) UpdateScriptTemplate(id string, input model.ScriptTemplate) (*model.ScriptTemplate, error) {
	st, err := s.store.GetScriptTemplate(id)
	if err != nil {
		return nil, err
	}
	st.Name = input.Name
	st.ScriptType = input.ScriptType
	st.Content = input.Content
	st.Description = input.Description
	if err := st.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateScriptTemplate(st); err != nil {
		return nil, err
	}
	return st, nil
}

func (s *Service) DeleteScriptTemplate(id string) error {
	return s.store.DeleteScriptTemplate(id)
}
