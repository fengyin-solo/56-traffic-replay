package store

import "replay/internal/model"

func (s *MemoryStore) CreateTrafficRecord(t *model.TrafficRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.trafficRecords[t.ID] = t
	return nil
}

func (s *MemoryStore) GetTrafficRecord(id string) (*model.TrafficRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.trafficRecords[id]
	if !ok {
		return nil, ErrNotFound
	}
	return t, nil
}

func (s *MemoryStore) ListTrafficRecords() []*model.TrafficRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.TrafficRecord, 0, len(s.trafficRecords))
	for _, t := range s.trafficRecords {
		list = append(list, t)
	}
	return list
}

func (s *MemoryStore) UpdateTrafficRecord(t *model.TrafficRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.trafficRecords[t.ID]; !ok {
		return ErrNotFound
	}
	s.trafficRecords[t.ID] = t
	return nil
}

func (s *MemoryStore) DeleteTrafficRecord(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.trafficRecords[id]; !ok {
		return ErrNotFound
	}
	delete(s.trafficRecords, id)
	return nil
}
