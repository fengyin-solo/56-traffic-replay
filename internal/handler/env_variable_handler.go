package handler

import (
	"net/http"

	"replay/internal/model"
	"replay/pkg/httpx"
)

func (s *Server) registerEnvVariableRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/env-variables", s.createEnvVariable)
	mux.HandleFunc("GET /api/env-variables", s.listEnvVariables)
	mux.HandleFunc("GET /api/env-variables/{id}", s.getEnvVariable)
	mux.HandleFunc("PUT /api/env-variables/{id}", s.updateEnvVariable)
	mux.HandleFunc("DELETE /api/env-variables/{id}", s.deleteEnvVariable)
}

type createEnvVariableRequest struct {
	EnvID    string `json:"env_id"`
	Key      string `json:"key"`
	Value    string `json:"value"`
	IsSecret bool   `json:"is_secret"`
}

func (s *Server) createEnvVariable(w http.ResponseWriter, r *http.Request) {
	var req createEnvVariableRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	e, err := s.svc.CreateEnvVariable(model.EnvVariable{
		EnvID:    req.EnvID,
		Key:      req.Key,
		Value:    req.Value,
		IsSecret: req.IsSecret,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, e)
}

func (s *Server) listEnvVariables(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.EnvVariableFilter{
		EnvID: r.URL.Query().Get("env_id"),
		Key:   r.URL.Query().Get("key"),
	}
	items, total, err := s.svc.ListEnvVariables(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getEnvVariable(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	e, err := s.svc.GetEnvVariable(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, e)
}

type updateEnvVariableRequest struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	IsSecret bool   `json:"is_secret"`
}

func (s *Server) updateEnvVariable(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateEnvVariableRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	e, err := s.svc.UpdateEnvVariable(id, model.EnvVariable{
		Key:      req.Key,
		Value:    req.Value,
		IsSecret: req.IsSecret,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, e)
}

func (s *Server) deleteEnvVariable(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteEnvVariable(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
