package handler

import (
	"net/http"

	"replay/internal/model"
	"replay/pkg/httpx"
)

func (s *Server) registerFilterRuleRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/filter-rules", s.createFilterRule)
	mux.HandleFunc("GET /api/filter-rules", s.listFilterRules)
	mux.HandleFunc("GET /api/filter-rules/{id}", s.getFilterRule)
	mux.HandleFunc("PUT /api/filter-rules/{id}", s.updateFilterRule)
	mux.HandleFunc("DELETE /api/filter-rules/{id}", s.deleteFilterRule)
}

type createFilterRuleRequest struct {
	SessionID string `json:"session_id"`
	Field     string `json:"field"`
	Operator  string `json:"operator"`
	Value     string `json:"value"`
	Action    string `json:"action"`
	Status    string `json:"status"`
}

func (s *Server) createFilterRule(w http.ResponseWriter, r *http.Request) {
	var req createFilterRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rule, err := s.svc.CreateFilterRule(model.FilterRule{
		SessionID: req.SessionID,
		Field:     req.Field,
		Operator:  req.Operator,
		Value:     req.Value,
		Action:    req.Action,
		Status:    req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, rule)
}

func (s *Server) listFilterRules(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.FilterRuleFilter{
		SessionID: r.URL.Query().Get("session_id"),
		Field:     r.URL.Query().Get("field"),
		Action:    r.URL.Query().Get("action"),
		Status:    r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListFilterRules(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getFilterRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	rule, err := s.svc.GetFilterRule(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rule)
}

type updateFilterRuleRequest struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
	Action   string `json:"action"`
	Status   string `json:"status"`
}

func (s *Server) updateFilterRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateFilterRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rule, err := s.svc.UpdateFilterRule(id, model.FilterRule{
		Field:    req.Field,
		Operator: req.Operator,
		Value:    req.Value,
		Action:   req.Action,
		Status:   req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rule)
}

func (s *Server) deleteFilterRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteFilterRule(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
