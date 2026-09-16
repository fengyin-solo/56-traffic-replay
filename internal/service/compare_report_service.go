package service

import (
	"sort"
	"strings"
	"time"

	"replay/internal/model"
	"replay/pkg/idgen"
)

func (s *Service) CreateCompareReport(input model.CompareReport) (*model.CompareReport, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetReplayTask(input.ReplayTaskID); err != nil {
		return nil, model.NewValidationError("replay_task_id", "关联回放任务不存在")
	}
	report := &model.CompareReport{
		ID:            idgen.Hex(),
		ReplayTaskID:  input.ReplayTaskID,
		TotalCompared: 0,
		Matched:       0,
		Mismatched:    0,
		DiffCount:     0,
		Status:        model.ReportPending,
		CreatedAt:     time.Now(),
	}
	if err := s.store.CreateCompareReport(report); err != nil {
		return nil, err
	}
	return report, nil
}

func (s *Service) ListCompareReports(filter model.CompareReportFilter, page, size int) ([]*model.CompareReport, int, error) {
	all := s.store.ListCompareReports()
	matched := make([]*model.CompareReport, 0, len(all))
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
		return []*model.CompareReport{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetCompareReport(id string) (*model.CompareReport, error) {
	return s.store.GetCompareReport(id)
}

func (s *Service) UpdateCompareReport(id string, input model.CompareReport) (*model.CompareReport, error) {
	report, err := s.store.GetCompareReport(id)
	if err != nil {
		return nil, err
	}
	report.TotalCompared = input.TotalCompared
	report.Matched = input.Matched
	report.Mismatched = input.Mismatched
	report.DiffCount = input.DiffCount
	if err := report.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateCompareReport(report); err != nil {
		return nil, err
	}
	return report, nil
}

func (s *Service) DeleteCompareReport(id string) error {
	return s.store.DeleteCompareReport(id)
}

func (s *Service) RunCompareReport(id string) (*model.CompareReport, error) {
	report, err := s.store.GetCompareReport(id)
	if err != nil {
		return nil, err
	}
	if !model.ReportCanTransition(report.Status, model.ReportRunning) {
		return nil, model.NewValidationError("status", "当前状态不允许启动对比")
	}
	report.Status = model.ReportRunning
	if err := s.store.UpdateCompareReport(report); err != nil {
		return nil, err
	}

	results := s.store.ListReplayResults()
	var taskResults []*model.ReplayResult
	for _, r := range results {
		if r.ReplayTaskID == report.ReplayTaskID {
			taskResults = append(taskResults, r)
		}
	}

	report.TotalCompared = len(taskResults)
	report.Matched = 0
	report.Mismatched = 0
	report.DiffCount = 0

	for _, r := range taskResults {
		record, err := s.store.GetTrafficRecord(r.TrafficRecordID)
		if err != nil {
			continue
		}
		original := record.Headers + "\n" + record.Body
		replayed := "status=" + string(rune('0'+r.ResponseCode/100)) + "\n"
		matched := r.Status == model.ResultSuccess
		if matched {
			report.Matched++
		} else {
			report.Mismatched++
		}

		diffDetail := ""
		if !matched {
			diffDetail = computeDiff(original, replayed)
			report.DiffCount++
		}

		diff := &model.DiffRecord{
			ID:               idgen.Hex(),
			ReportID:         report.ID,
			TrafficRecordID:  r.TrafficRecordID,
			OriginalResponse: original,
			ReplayedResponse: replayed,
			DiffDetail:       diffDetail,
			CreatedAt:        time.Now(),
		}
		_ = s.store.CreateDiffRecord(diff)
	}

	report.Status = model.ReportCompleted
	now := time.Now()
	report.FinishedAt = &now
	if err := s.store.UpdateCompareReport(report); err != nil {
		return nil, err
	}
	return report, nil
}

func computeDiff(a, b string) string {
	if a == b {
		return ""
	}
	var diffs []string
	linesA := strings.Split(a, "\n")
	linesB := strings.Split(b, "\n")
	maxLen := len(linesA)
	if len(linesB) > maxLen {
		maxLen = len(linesB)
	}
	for i := 0; i < maxLen; i++ {
		la := ""
		if i < len(linesA) {
			la = linesA[i]
		}
		lb := ""
		if i < len(linesB) {
			lb = linesB[i]
		}
		if la != lb {
			diffs = append(diffs, "- "+la)
			diffs = append(diffs, "+ "+lb)
		}
	}
	return strings.Join(diffs, "\n")
}
