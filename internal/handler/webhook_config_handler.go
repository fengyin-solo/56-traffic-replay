package handler

import (
	"net/http"

	"replay/internal/model"
	"replay/pkg/httpx"
)

func (s *Server) registerWebhookConfigRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/webhook-configs", s.createWebhookConfig)
	mux.HandleFunc("GET /api/webhook-configs", s.listWebhookConfigs)
	mux.HandleFunc("GET /api/webhook-configs/{id}", s.getWebhookConfig)
	mux.HandleFunc("PUT /api/webhook-configs/{id}", s.updateWebhookConfig)
	mux.HandleFunc("DELETE /api/webhook-configs/{id}", s.deleteWebhookConfig)
}

type createWebhookConfigRequest struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Secret string `json:"secret"`
	Events string `json:"events"`
	Status string `json:"status"`
}

func (s *Server) createWebhookConfig(w http.ResponseWriter, r *http.Request) {
	var req createWebhookConfigRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	wc, err := s.svc.CreateWebhookConfig(model.WebhookConfig{
		Name:   req.Name,
		URL:    req.URL,
		Secret: req.Secret,
		Events: req.Events,
		Status: req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, wc)
}

func (s *Server) listWebhookConfigs(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.WebhookConfigFilter{
		Name:    r.URL.Query().Get("name"),
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListWebhookConfigs(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getWebhookConfig(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	wc, err := s.svc.GetWebhookConfig(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, wc)
}

type updateWebhookConfigRequest struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Secret string `json:"secret"`
	Events string `json:"events"`
	Status string `json:"status"`
}

func (s *Server) updateWebhookConfig(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateWebhookConfigRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	wc, err := s.svc.UpdateWebhookConfig(id, model.WebhookConfig{
		Name:   req.Name,
		URL:    req.URL,
		Secret: req.Secret,
		Events: req.Events,
		Status: req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, wc)
}

func (s *Server) deleteWebhookConfig(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteWebhookConfig(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
