package store

import "replay/internal/model"

func (s *MemoryStore) CreateWebhookConfig(w *model.WebhookConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.webhookConfigs {
		if exist.Name == w.Name {
			return ErrConflict
		}
	}
	s.webhookConfigs[w.ID] = w
	return nil
}

func (s *MemoryStore) GetWebhookConfig(id string) (*model.WebhookConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	w, ok := s.webhookConfigs[id]
	if !ok {
		return nil, ErrNotFound
	}
	return w, nil
}

func (s *MemoryStore) GetWebhookConfigByName(name string) (*model.WebhookConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, w := range s.webhookConfigs {
		if w.Name == name {
			return w, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListWebhookConfigs() []*model.WebhookConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.WebhookConfig, 0, len(s.webhookConfigs))
	for _, w := range s.webhookConfigs {
		list = append(list, w)
	}
	return list
}

func (s *MemoryStore) UpdateWebhookConfig(w *model.WebhookConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.webhookConfigs[w.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.webhookConfigs {
		if exist.ID != w.ID && exist.Name == w.Name {
			return ErrConflict
		}
	}
	s.webhookConfigs[w.ID] = w
	return nil
}

func (s *MemoryStore) DeleteWebhookConfig(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.webhookConfigs[id]; !ok {
		return ErrNotFound
	}
	delete(s.webhookConfigs, id)
	return nil
}
