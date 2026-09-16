package handler

import (
	"net/http"

	"replay/internal/model"
	"replay/pkg/httpx"
)

func (s *Server) registerTrafficRecordRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/traffic-records", s.createTrafficRecord)
	mux.HandleFunc("GET /api/traffic-records", s.listTrafficRecords)
	mux.HandleFunc("GET /api/traffic-records/{id}", s.getTrafficRecord)
	mux.HandleFunc("PUT /api/traffic-records/{id}", s.updateTrafficRecord)
	mux.HandleFunc("DELETE /api/traffic-records/{id}", s.deleteTrafficRecord)
	mux.HandleFunc("POST /api/traffic-records/import", s.importTrafficRecords)
}

type createTrafficRecordRequest struct {
	SessionID  string `json:"session_id"`
	Method     string `json:"method"`
	Path       string `json:"path"`
	Headers    string `json:"headers"`
	Body       string `json:"body"`
	DurationMs int    `json:"duration_ms"`
	StatusCode int    `json:"status_code"`
}

func (s *Server) createTrafficRecord(w http.ResponseWriter, r *http.Request) {
	var req createTrafficRecordRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	record, err := s.svc.CreateTrafficRecord(model.TrafficRecord{
		SessionID:  req.SessionID,
		Method:     req.Method,
		Path:       req.Path,
		Headers:    req.Headers,
		Body:       req.Body,
		DurationMs: req.DurationMs,
		StatusCode: req.StatusCode,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, record)
}

func (s *Server) listTrafficRecords(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.TrafficRecordFilter{
		SessionID: r.URL.Query().Get("session_id"),
		Method:    r.URL.Query().Get("method"),
		Path:      r.URL.Query().Get("path"),
		Keyword:   r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListTrafficRecords(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getTrafficRecord(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	record, err := s.svc.GetTrafficRecord(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, record)
}

type updateTrafficRecordRequest struct {
	Method     string `json:"method"`
	Path       string `json:"path"`
	Headers    string `json:"headers"`
	Body       string `json:"body"`
	DurationMs int    `json:"duration_ms"`
	StatusCode int    `json:"status_code"`
}

func (s *Server) updateTrafficRecord(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateTrafficRecordRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	record, err := s.svc.UpdateTrafficRecord(id, model.TrafficRecord{
		Method:     req.Method,
		Path:       req.Path,
		Headers:    req.Headers,
		Body:       req.Body,
		DurationMs: req.DurationMs,
		StatusCode: req.StatusCode,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, record)
}

func (s *Server) deleteTrafficRecord(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteTrafficRecord(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type importTrafficRecordsRequest struct {
	SessionID string                   `json:"session_id"`
	Records   []createTrafficRecordRequest `json:"records"`
}

func (s *Server) importTrafficRecords(w http.ResponseWriter, r *http.Request) {
	var req importTrafficRecordsRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	inputs := make([]model.TrafficRecord, 0, len(req.Records))
	for _, rec := range req.Records {
		inputs = append(inputs, model.TrafficRecord{
			Method:     rec.Method,
			Path:       rec.Path,
			Headers:    rec.Headers,
			Body:       rec.Body,
			DurationMs: rec.DurationMs,
			StatusCode: rec.StatusCode,
		})
	}
	count, err := s.svc.ImportTrafficRecords(req.SessionID, inputs)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]int{"imported": count})
}
