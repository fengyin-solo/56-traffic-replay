package store

import "replay/internal/model"

func (s *MemoryStore) CreateEnvVariable(e *model.EnvVariable) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.envVariables[e.ID] = e
	return nil
}

func (s *MemoryStore) GetEnvVariable(id string) (*model.EnvVariable, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.envVariables[id]
	if !ok {
		return nil, ErrNotFound
	}
	return e, nil
}

func (s *MemoryStore) ListEnvVariables() []*model.EnvVariable {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.EnvVariable, 0, len(s.envVariables))
	for _, e := range s.envVariables {
		list = append(list, e)
	}
	return list
}

func (s *MemoryStore) UpdateEnvVariable(e *model.EnvVariable) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.envVariables[e.ID]; !ok {
		return ErrNotFound
	}
	s.envVariables[e.ID] = e
	return nil
}

func (s *MemoryStore) DeleteEnvVariable(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.envVariables[id]; !ok {
		return ErrNotFound
	}
	delete(s.envVariables, id)
	return nil
}
