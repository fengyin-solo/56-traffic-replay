package store

import "replay/internal/model"

func (s *MemoryStore) CreateDiffRecord(d *model.DiffRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.diffRecords[d.ID] = d
	return nil
}

func (s *MemoryStore) GetDiffRecord(id string) (*model.DiffRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	d, ok := s.diffRecords[id]
	if !ok {
		return nil, ErrNotFound
	}
	return d, nil
}

func (s *MemoryStore) ListDiffRecords() []*model.DiffRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.DiffRecord, 0, len(s.diffRecords))
	for _, d := range s.diffRecords {
		list = append(list, d)
	}
	return list
}

func (s *MemoryStore) UpdateDiffRecord(d *model.DiffRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.diffRecords[d.ID]; !ok {
		return ErrNotFound
	}
	s.diffRecords[d.ID] = d
	return nil
}

func (s *MemoryStore) DeleteDiffRecord(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.diffRecords[id]; !ok {
		return ErrNotFound
	}
	delete(s.diffRecords, id)
	return nil
}
