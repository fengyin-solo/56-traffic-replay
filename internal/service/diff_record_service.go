package service

import (
	"sort"
	"time"

	"replay/internal/model"
	"replay/pkg/idgen"
)

func (s *Service) CreateDiffRecord(input model.DiffRecord) (*model.DiffRecord, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetCompareReport(input.ReportID); err != nil {
		return nil, model.NewValidationError("report_id", "关联报告不存在")
	}
	d := &model.DiffRecord{
		ID:               idgen.Hex(),
		ReportID:         input.ReportID,
		TrafficRecordID:  input.TrafficRecordID,
		OriginalResponse: input.OriginalResponse,
		ReplayedResponse: input.ReplayedResponse,
		DiffDetail:       input.DiffDetail,
		CreatedAt:        time.Now(),
	}
	if err := s.store.CreateDiffRecord(d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) ListDiffRecords(filter model.DiffRecordFilter, page, size int) ([]*model.DiffRecord, int, error) {
	all := s.store.ListDiffRecords()
	matched := make([]*model.DiffRecord, 0, len(all))
	for _, d := range all {
		if filter.Match(d) {
			matched = append(matched, d)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.DiffRecord{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetDiffRecord(id string) (*model.DiffRecord, error) {
	return s.store.GetDiffRecord(id)
}

func (s *Service) UpdateDiffRecord(id string, input model.DiffRecord) (*model.DiffRecord, error) {
	d, err := s.store.GetDiffRecord(id)
	if err != nil {
		return nil, err
	}
	d.OriginalResponse = input.OriginalResponse
	d.ReplayedResponse = input.ReplayedResponse
	d.DiffDetail = input.DiffDetail
	if err := d.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateDiffRecord(d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) DeleteDiffRecord(id string) error {
	return s.store.DeleteDiffRecord(id)
}
