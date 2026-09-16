package handler

import (
	"net/http"

	"replay/internal/model"
	"replay/pkg/httpx"
)

func (s *Server) registerDiffRecordRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/diff-records", s.createDiffRecord)
	mux.HandleFunc("GET /api/diff-records", s.listDiffRecords)
	mux.HandleFunc("GET /api/diff-records/{id}", s.getDiffRecord)
	mux.HandleFunc("PUT /api/diff-records/{id}", s.updateDiffRecord)
	mux.HandleFunc("DELETE /api/diff-records/{id}", s.deleteDiffRecord)
}

type createDiffRecordRequest struct {
	ReportID         string `json:"report_id"`
	TrafficRecordID  string `json:"traffic_record_id"`
	OriginalResponse string `json:"original_response"`
	ReplayedResponse string `json:"replayed_response"`
	DiffDetail       string `json:"diff_detail"`
}

func (s *Server) createDiffRecord(w http.ResponseWriter, r *http.Request) {
	var req createDiffRecordRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	d, err := s.svc.CreateDiffRecord(model.DiffRecord{
		ReportID:         req.ReportID,
		TrafficRecordID:  req.TrafficRecordID,
		OriginalResponse: req.OriginalResponse,
		ReplayedResponse: req.ReplayedResponse,
		DiffDetail:       req.DiffDetail,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, d)
}

func (s *Server) listDiffRecords(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.DiffRecordFilter{
		ReportID:        r.URL.Query().Get("report_id"),
		TrafficRecordID: r.URL.Query().Get("traffic_record_id"),
	}
	items, total, err := s.svc.ListDiffRecords(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getDiffRecord(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	d, err := s.svc.GetDiffRecord(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, d)
}

type updateDiffRecordRequest struct {
	OriginalResponse string `json:"original_response"`
	ReplayedResponse string `json:"replayed_response"`
	DiffDetail       string `json:"diff_detail"`
}

func (s *Server) updateDiffRecord(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateDiffRecordRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	d, err := s.svc.UpdateDiffRecord(id, model.DiffRecord{
		OriginalResponse: req.OriginalResponse,
		ReplayedResponse: req.ReplayedResponse,
		DiffDetail:       req.DiffDetail,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, d)
}

func (s *Server) deleteDiffRecord(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteDiffRecord(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
