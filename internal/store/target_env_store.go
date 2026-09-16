package store

import "replay/internal/model"

func (s *MemoryStore) CreateTargetEnv(e *model.TargetEnv) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.targetEnvs {
		if exist.Name == e.Name {
			return ErrConflict
		}
	}
	s.targetEnvs[e.ID] = e
	return nil
}

func (s *MemoryStore) GetTargetEnv(id string) (*model.TargetEnv, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.targetEnvs[id]
	if !ok {
		return nil, ErrNotFound
	}
	return e, nil
}

func (s *MemoryStore) GetTargetEnvByName(name string) (*model.TargetEnv, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, e := range s.targetEnvs {
		if e.Name == name {
			return e, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListTargetEnvs() []*model.TargetEnv {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.TargetEnv, 0, len(s.targetEnvs))
	for _, e := range s.targetEnvs {
		list = append(list, e)
	}
	return list
}

func (s *MemoryStore) UpdateTargetEnv(e *model.TargetEnv) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.targetEnvs[e.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.targetEnvs {
		if exist.ID != e.ID && exist.Name == e.Name {
			return ErrConflict
		}
	}
	s.targetEnvs[e.ID] = e
	return nil
}

func (s *MemoryStore) DeleteTargetEnv(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.targetEnvs[id]; !ok {
		return ErrNotFound
	}
	delete(s.targetEnvs, id)
	return nil
}
