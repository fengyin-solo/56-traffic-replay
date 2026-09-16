package store

import "replay/internal/model"

func (s *MemoryStore) CreateCompareReport(r *model.CompareReport) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.compareReports[r.ID] = r
	return nil
}

func (s *MemoryStore) GetCompareReport(id string) (*model.CompareReport, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.compareReports[id]
	if !ok {
		return nil, ErrNotFound
	}
	return r, nil
}

func (s *MemoryStore) ListCompareReports() []*model.CompareReport {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.CompareReport, 0, len(s.compareReports))
	for _, r := range s.compareReports {
		list = append(list, r)
	}
	return list
}

func (s *MemoryStore) UpdateCompareReport(r *model.CompareReport) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.compareReports[r.ID]; !ok {
		return ErrNotFound
	}
	s.compareReports[r.ID] = r
	return nil
}

func (s *MemoryStore) DeleteCompareReport(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.compareReports[id]; !ok {
		return ErrNotFound
	}
	delete(s.compareReports, id)
	return nil
}
