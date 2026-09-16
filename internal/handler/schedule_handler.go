package handler

import (
	"net/http"

	"replay/internal/model"
	"replay/pkg/httpx"
)

func (s *Server) registerScheduleRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/schedules", s.createSchedule)
	mux.HandleFunc("GET /api/schedules", s.listSchedules)
	mux.HandleFunc("GET /api/schedules/{id}", s.getSchedule)
	mux.HandleFunc("PUT /api/schedules/{id}", s.updateSchedule)
	mux.HandleFunc("DELETE /api/schedules/{id}", s.deleteSchedule)
}

type createScheduleRequest struct {
	ReplayTaskID string `json:"replay_task_id"`
	CronExpr     string `json:"cron_expr"`
	Status       string `json:"status"`
}

func (s *Server) createSchedule(w http.ResponseWriter, r *http.Request) {
	var req createScheduleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sch, err := s.svc.CreateSchedule(model.Schedule{
		ReplayTaskID: req.ReplayTaskID,
		CronExpr:     req.CronExpr,
		Status:       req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, sch)
}

func (s *Server) listSchedules(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ScheduleFilter{
		ReplayTaskID: r.URL.Query().Get("replay_task_id"),
		Status:       r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListSchedules(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getSchedule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	sch, err := s.svc.GetSchedule(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sch)
}

type updateScheduleRequest struct {
	CronExpr string `json:"cron_expr"`
	Status   string `json:"status"`
}

func (s *Server) updateSchedule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateScheduleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sch, err := s.svc.UpdateSchedule(id, model.Schedule{
		CronExpr: req.CronExpr,
		Status:   req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sch)
}

func (s *Server) deleteSchedule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteSchedule(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
