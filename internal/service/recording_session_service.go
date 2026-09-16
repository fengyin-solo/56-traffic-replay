package service

import (
	"sort"
	"time"

	"replay/internal/model"
	"replay/pkg/idgen"
)

func (s *Service) CreateRecordingSession(input model.RecordingSession) (*model.RecordingSession, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	sess := &model.RecordingSession{
		ID:         idgen.Hex(),
		Name:       input.Name,
		Source:     input.Source,
		TargetHost: input.TargetHost,
		Protocol:   input.Protocol,
		Status:     model.SessionRecording,
		StartedAt:  time.Now(),
		CreatedAt:  time.Now(),
	}
	if err := s.store.CreateRecordingSession(sess); err != nil {
		return nil, err
	}
	return sess, nil
}

func (s *Service) ListRecordingSessions(filter model.RecordingSessionFilter, page, size int) ([]*model.RecordingSession, int, error) {
	all := s.store.ListRecordingSessions()
	matched := make([]*model.RecordingSession, 0, len(all))
	for _, sess := range all {
		if filter.Match(sess) {
			matched = append(matched, sess)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.RecordingSession{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetRecordingSession(id string) (*model.RecordingSession, error) {
	return s.store.GetRecordingSession(id)
}

func (s *Service) UpdateRecordingSession(id string, input model.RecordingSession) (*model.RecordingSession, error) {
	sess, err := s.store.GetRecordingSession(id)
	if err != nil {
		return nil, err
	}
	sess.Name = input.Name
	sess.Source = input.Source
	sess.TargetHost = input.TargetHost
	sess.Protocol = input.Protocol
	if err := sess.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateRecordingSession(sess); err != nil {
		return nil, err
	}
	return sess, nil
}

func (s *Service) DeleteRecordingSession(id string) error {
	return s.store.DeleteRecordingSession(id)
}

func (s *Service) StartRecording(id string) (*model.RecordingSession, error) {
	sess, err := s.store.GetRecordingSession(id)
	if err != nil {
		return nil, err
	}
	if !model.SessionCanTransition(sess.Status, model.SessionRecording) {
		return nil, model.NewValidationError("status", "当前状态不允许开始录制")
	}
	sess.Status = model.SessionRecording
	sess.StartedAt = time.Now()
	sess.StoppedAt = nil
	if err := s.store.UpdateRecordingSession(sess); err != nil {
		return nil, err
	}
	return sess, nil
}

func (s *Service) PauseRecording(id string) (*model.RecordingSession, error) {
	sess, err := s.store.GetRecordingSession(id)
	if err != nil {
		return nil, err
	}
	if !model.SessionCanTransition(sess.Status, model.SessionPaused) {
		return nil, model.NewValidationError("status", "当前状态不允许暂停录制")
	}
	sess.Status = model.SessionPaused
	if err := s.store.UpdateRecordingSession(sess); err != nil {
		return nil, err
	}
	return sess, nil
}

func (s *Service) StopRecording(id string) (*model.RecordingSession, error) {
	sess, err := s.store.GetRecordingSession(id)
	if err != nil {
		return nil, err
	}
	if !model.SessionCanTransition(sess.Status, model.SessionCompleted) {
		return nil, model.NewValidationError("status", "当前状态不允许停止录制")
	}
	sess.Status = model.SessionCompleted
	now := time.Now()
	sess.StoppedAt = &now
	if err := s.store.UpdateRecordingSession(sess); err != nil {
		return nil, err
	}
	return sess, nil
}
