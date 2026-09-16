package store

import "replay/internal/model"

func (s *MemoryStore) CreateScriptTemplate(st *model.ScriptTemplate) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.scriptTemplates[st.ID] = st
	return nil
}

func (s *MemoryStore) GetScriptTemplate(id string) (*model.ScriptTemplate, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	st, ok := s.scriptTemplates[id]
	if !ok {
		return nil, ErrNotFound
	}
	return st, nil
}

func (s *MemoryStore) ListScriptTemplates() []*model.ScriptTemplate {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.ScriptTemplate, 0, len(s.scriptTemplates))
	for _, st := range s.scriptTemplates {
		list = append(list, st)
	}
	return list
}

func (s *MemoryStore) UpdateScriptTemplate(st *model.ScriptTemplate) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.scriptTemplates[st.ID]; !ok {
		return ErrNotFound
	}
	s.scriptTemplates[st.ID] = st
	return nil
}

func (s *MemoryStore) DeleteScriptTemplate(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.scriptTemplates[id]; !ok {
		return ErrNotFound
	}
	delete(s.scriptTemplates, id)
	return nil
}
