package handler

import (
	"net/http"

	"replay/internal/model"
	"replay/pkg/httpx"
)

func (s *Server) registerNotificationRuleRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/notification-rules", s.createNotificationRule)
	mux.HandleFunc("GET /api/notification-rules", s.listNotificationRules)
	mux.HandleFunc("GET /api/notification-rules/{id}", s.getNotificationRule)
	mux.HandleFunc("PUT /api/notification-rules/{id}", s.updateNotificationRule)
	mux.HandleFunc("DELETE /api/notification-rules/{id}", s.deleteNotificationRule)
}

type createNotificationRuleRequest struct {
	Name    string `json:"name"`
	Channel string `json:"channel"`
	Target  string `json:"target"`
	Events  string `json:"events"`
	Status  string `json:"status"`
}

func (s *Server) createNotificationRule(w http.ResponseWriter, r *http.Request) {
	var req createNotificationRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	n, err := s.svc.CreateNotificationRule(model.NotificationRule{
		Name:    req.Name,
		Channel: req.Channel,
		Target:  req.Target,
		Events:  req.Events,
		Status:  req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, n)
}

func (s *Server) listNotificationRules(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.NotificationRuleFilter{
		Name:    r.URL.Query().Get("name"),
		Channel: r.URL.Query().Get("channel"),
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListNotificationRules(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getNotificationRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	n, err := s.svc.GetNotificationRule(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, n)
}

type updateNotificationRuleRequest struct {
	Name    string `json:"name"`
	Channel string `json:"channel"`
	Target  string `json:"target"`
	Events  string `json:"events"`
	Status  string `json:"status"`
}

func (s *Server) updateNotificationRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateNotificationRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	n, err := s.svc.UpdateNotificationRule(id, model.NotificationRule{
		Name:    req.Name,
		Channel: req.Channel,
		Target:  req.Target,
		Events:  req.Events,
		Status:  req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, n)
}

func (s *Server) deleteNotificationRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteNotificationRule(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
