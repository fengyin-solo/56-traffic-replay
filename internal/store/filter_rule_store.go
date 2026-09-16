package store

import "replay/internal/model"

func (s *MemoryStore) CreateFilterRule(f *model.FilterRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.filterRules[f.ID] = f
	return nil
}

func (s *MemoryStore) GetFilterRule(id string) (*model.FilterRule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	f, ok := s.filterRules[id]
	if !ok {
		return nil, ErrNotFound
	}
	return f, nil
}

func (s *MemoryStore) ListFilterRules() []*model.FilterRule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.FilterRule, 0, len(s.filterRules))
	for _, f := range s.filterRules {
		list = append(list, f)
	}
	return list
}

func (s *MemoryStore) UpdateFilterRule(f *model.FilterRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.filterRules[f.ID]; !ok {
		return ErrNotFound
	}
	s.filterRules[f.ID] = f
	return nil
}

func (s *MemoryStore) DeleteFilterRule(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.filterRules[id]; !ok {
		return ErrNotFound
	}
	delete(s.filterRules, id)
	return nil
}
