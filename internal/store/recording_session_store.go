package store

import "replay/internal/model"

func (s *MemoryStore) CreateRecordingSession(sess *model.RecordingSession) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.sessions {
		if exist.Name == sess.Name {
			return ErrConflict
		}
	}
	s.sessions[sess.ID] = sess
	return nil
}

func (s *MemoryStore) GetRecordingSession(id string) (*model.RecordingSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sess, ok := s.sessions[id]
	if !ok {
		return nil, ErrNotFound
	}
	return sess, nil
}

func (s *MemoryStore) GetRecordingSessionByName(name string) (*model.RecordingSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, sess := range s.sessions {
		if sess.Name == name {
			return sess, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListRecordingSessions() []*model.RecordingSession {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.RecordingSession, 0, len(s.sessions))
	for _, sess := range s.sessions {
		list = append(list, sess)
	}
	return list
}

func (s *MemoryStore) UpdateRecordingSession(sess *model.RecordingSession) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.sessions[sess.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.sessions {
		if exist.ID != sess.ID && exist.Name == sess.Name {
			return ErrConflict
		}
	}
	s.sessions[sess.ID] = sess
	return nil
}

func (s *MemoryStore) DeleteRecordingSession(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.sessions[id]; !ok {
		return ErrNotFound
	}
	delete(s.sessions, id)
	return nil
}
