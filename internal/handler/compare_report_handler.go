package handler

import (
	"net/http"

	"replay/internal/model"
	"replay/pkg/httpx"
)

func (s *Server) registerCompareReportRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/compare-reports", s.createCompareReport)
	mux.HandleFunc("GET /api/compare-reports", s.listCompareReports)
	mux.HandleFunc("GET /api/compare-reports/{id}", s.getCompareReport)
	mux.HandleFunc("PUT /api/compare-reports/{id}", s.updateCompareReport)
	mux.HandleFunc("DELETE /api/compare-reports/{id}", s.deleteCompareReport)
	mux.HandleFunc("POST /api/compare-reports/{id}/run", s.runCompareReport)
}

type createCompareReportRequest struct {
	ReplayTaskID string `json:"replay_task_id"`
}

func (s *Server) createCompareReport(w http.ResponseWriter, r *http.Request) {
	var req createCompareReportRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	report, err := s.svc.CreateCompareReport(model.CompareReport{ReplayTaskID: req.ReplayTaskID})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, report)
}

func (s *Server) listCompareReports(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.CompareReportFilter{
		ReplayTaskID: r.URL.Query().Get("replay_task_id"),
		Status:       r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListCompareReports(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getCompareReport(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	report, err := s.svc.GetCompareReport(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, report)
}

type updateCompareReportRequest struct {
	TotalCompared int `json:"total_compared"`
	Matched       int `json:"matched"`
	Mismatched    int `json:"mismatched"`
	DiffCount     int `json:"diff_count"`
}

func (s *Server) updateCompareReport(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateCompareReportRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	report, err := s.svc.UpdateCompareReport(id, model.CompareReport{
		TotalCompared: req.TotalCompared,
		Matched:       req.Matched,
		Mismatched:    req.Mismatched,
		DiffCount:     req.DiffCount,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, report)
}

func (s *Server) deleteCompareReport(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteCompareReport(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) runCompareReport(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	report, err := s.svc.RunCompareReport(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, report)
}
