package store

import "replay/internal/model"

func (s *MemoryStore) CreateNotificationRule(n *model.NotificationRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.notificationRules[n.ID] = n
	return nil
}

func (s *MemoryStore) GetNotificationRule(id string) (*model.NotificationRule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n, ok := s.notificationRules[id]
	if !ok {
		return nil, ErrNotFound
	}
	return n, nil
}

func (s *MemoryStore) ListNotificationRules() []*model.NotificationRule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.NotificationRule, 0, len(s.notificationRules))
	for _, n := range s.notificationRules {
		list = append(list, n)
	}
	return list
}

func (s *MemoryStore) UpdateNotificationRule(n *model.NotificationRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.notificationRules[n.ID]; !ok {
		return ErrNotFound
	}
	s.notificationRules[n.ID] = n
	return nil
}

func (s *MemoryStore) DeleteNotificationRule(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.notificationRules[id]; !ok {
		return ErrNotFound
	}
	delete(s.notificationRules, id)
	return nil
}
