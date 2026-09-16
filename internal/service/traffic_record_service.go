package service

import (
	"sort"
	"time"

	"replay/internal/model"
	"replay/pkg/idgen"
)

func (s *Service) CreateTrafficRecord(input model.TrafficRecord) (*model.TrafficRecord, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetRecordingSession(input.SessionID); err != nil {
		return nil, model.NewValidationError("session_id", "关联会话不存在")
	}
	record := &model.TrafficRecord{
		ID:         idgen.Hex(),
		SessionID:  input.SessionID,
		Method:     input.Method,
		Path:       input.Path,
		Headers:    input.Headers,
		Body:       input.Body,
		Timestamp:  time.Now(),
		DurationMs: input.DurationMs,
		StatusCode: input.StatusCode,
		CreatedAt:  time.Now(),
	}
	if err := s.store.CreateTrafficRecord(record); err != nil {
		return nil, err
	}
	return record, nil
}

func (s *Service) ListTrafficRecords(filter model.TrafficRecordFilter, page, size int) ([]*model.TrafficRecord, int, error) {
	all := s.store.ListTrafficRecords()
	matched := make([]*model.TrafficRecord, 0, len(all))
	for _, r := range all {
		if filter.Match(r) {
			matched = append(matched, r)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.TrafficRecord{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetTrafficRecord(id string) (*model.TrafficRecord, error) {
	return s.store.GetTrafficRecord(id)
}

func (s *Service) UpdateTrafficRecord(id string, input model.TrafficRecord) (*model.TrafficRecord, error) {
	record, err := s.store.GetTrafficRecord(id)
	if err != nil {
		return nil, err
	}
	record.Method = input.Method
	record.Path = input.Path
	record.Headers = input.Headers
	record.Body = input.Body
	record.DurationMs = input.DurationMs
	record.StatusCode = input.StatusCode
	if err := record.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateTrafficRecord(record); err != nil {
		return nil, err
	}
	return record, nil
}

func (s *Service) DeleteTrafficRecord(id string) error {
	return s.store.DeleteTrafficRecord(id)
}

func (s *Service) ImportTrafficRecords(sessionID string, inputs []model.TrafficRecord) (int, error) {
	if _, err := s.store.GetRecordingSession(sessionID); err != nil {
		return 0, model.NewValidationError("session_id", "关联会话不存在")
	}
	count := 0
	for _, input := range inputs {
		input.SessionID = sessionID
		if _, err := s.CreateTrafficRecord(input); err == nil {
			count++
		}
	}
	return count, nil
}
