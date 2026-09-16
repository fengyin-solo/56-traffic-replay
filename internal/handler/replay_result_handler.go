package handler

import (
	"net/http"

	"replay/internal/model"
	"replay/pkg/httpx"
)

func (s *Server) registerReplayResultRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/replay-results", s.createReplayResult)
	mux.HandleFunc("GET /api/replay-results", s.listReplayResults)
	mux.HandleFunc("GET /api/replay-results/{id}", s.getReplayResult)
	mux.HandleFunc("PUT /api/replay-results/{id}", s.updateReplayResult)
	mux.HandleFunc("DELETE /api/replay-results/{id}", s.deleteReplayResult)
}

type createReplayResultRequest struct {
	ReplayTaskID    string `json:"replay_task_id"`
	TrafficRecordID string `json:"traffic_record_id"`
	Status          string `json:"status"`
	ResponseCode    int    `json:"response_code"`
	DurationMs      int    `json:"duration_ms"`
}

func (s *Server) createReplayResult(w http.ResponseWriter, r *http.Request) {
	var req createReplayResultRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	res, err := s.svc.CreateReplayResult(model.ReplayResult{
		ReplayTaskID:    req.ReplayTaskID,
		TrafficRecordID: req.TrafficRecordID,
		Status:          req.Status,
		ResponseCode:    req.ResponseCode,
		DurationMs:      req.DurationMs,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, res)
}

func (s *Server) listReplayResults(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ReplayResultFilter{
		ReplayTaskID:    r.URL.Query().Get("replay_task_id"),
		TrafficRecordID: r.URL.Query().Get("traffic_record_id"),
		Status:          r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListReplayResults(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getReplayResult(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	res, err := s.svc.GetReplayResult(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, res)
}

type updateReplayResultRequest struct {
	Status       string `json:"status"`
	ResponseCode int    `json:"response_code"`
	DurationMs   int    `json:"duration_ms"`
}

func (s *Server) updateReplayResult(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateReplayResultRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	res, err := s.svc.UpdateReplayResult(id, model.ReplayResult{
		Status:       req.Status,
		ResponseCode: req.ResponseCode,
		DurationMs:   req.DurationMs,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, res)
}

func (s *Server) deleteReplayResult(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteReplayResult(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
