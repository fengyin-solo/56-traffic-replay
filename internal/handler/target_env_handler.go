package handler

import (
	"net/http"

	"replay/internal/model"
	"replay/pkg/httpx"
)

func (s *Server) registerTargetEnvRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/target-envs", s.createTargetEnv)
	mux.HandleFunc("GET /api/target-envs", s.listTargetEnvs)
	mux.HandleFunc("GET /api/target-envs/{id}", s.getTargetEnv)
	mux.HandleFunc("PUT /api/target-envs/{id}", s.updateTargetEnv)
	mux.HandleFunc("DELETE /api/target-envs/{id}", s.deleteTargetEnv)
}

type createTargetEnvRequest struct {
	Name      string `json:"name"`
	BaseURL   string `json:"base_url"`
	AuthToken string `json:"auth_token"`
	Status    string `json:"status"`
}

func (s *Server) createTargetEnv(w http.ResponseWriter, r *http.Request) {
	var req createTargetEnvRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	env, err := s.svc.CreateTargetEnv(model.TargetEnv{
		Name:      req.Name,
		BaseURL:   req.BaseURL,
		AuthToken: req.AuthToken,
		Status:    req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, env)
}

func (s *Server) listTargetEnvs(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.TargetEnvFilter{
		Name:    r.URL.Query().Get("name"),
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListTargetEnvs(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getTargetEnv(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	env, err := s.svc.GetTargetEnv(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, env)
}

type updateTargetEnvRequest struct {
	Name      string `json:"name"`
	BaseURL   string `json:"base_url"`
	AuthToken string `json:"auth_token"`
	Status    string `json:"status"`
}

func (s *Server) updateTargetEnv(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateTargetEnvRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	env, err := s.svc.UpdateTargetEnv(id, model.TargetEnv{
		Name:      req.Name,
		BaseURL:   req.BaseURL,
		AuthToken: req.AuthToken,
		Status:    req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, env)
}

func (s *Server) deleteTargetEnv(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteTargetEnv(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
