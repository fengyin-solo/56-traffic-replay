package store

import "replay/internal/model"

func (s *MemoryStore) CreateProjectConfig(p *model.ProjectConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.projectConfigs {
		if exist.Key == p.Key {
			return ErrConflict
		}
	}
	s.projectConfigs[p.ID] = p
	return nil
}

func (s *MemoryStore) GetProjectConfig(id string) (*model.ProjectConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.projectConfigs[id]
	if !ok {
		return nil, ErrNotFound
	}
	return p, nil
}

func (s *MemoryStore) GetProjectConfigByKey(key string) (*model.ProjectConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, p := range s.projectConfigs {
		if p.Key == key {
			return p, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListProjectConfigs() []*model.ProjectConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.ProjectConfig, 0, len(s.projectConfigs))
	for _, p := range s.projectConfigs {
		list = append(list, p)
	}
	return list
}

func (s *MemoryStore) UpdateProjectConfig(p *model.ProjectConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.projectConfigs[p.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.projectConfigs {
		if exist.ID != p.ID && exist.Key == p.Key {
			return ErrConflict
		}
	}
	s.projectConfigs[p.ID] = p
	return nil
}

func (s *MemoryStore) DeleteProjectConfig(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.projectConfigs[id]; !ok {
		return ErrNotFound
	}
	delete(s.projectConfigs, id)
	return nil
}
