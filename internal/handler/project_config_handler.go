package handler

import (
	"net/http"

	"replay/internal/model"
	"replay/pkg/httpx"
)

func (s *Server) registerProjectConfigRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/project-configs", s.createProjectConfig)
	mux.HandleFunc("GET /api/project-configs", s.listProjectConfigs)
	mux.HandleFunc("GET /api/project-configs/{id}", s.getProjectConfig)
	mux.HandleFunc("PUT /api/project-configs/{id}", s.updateProjectConfig)
	mux.HandleFunc("DELETE /api/project-configs/{id}", s.deleteProjectConfig)
}

type createProjectConfigRequest struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

func (s *Server) createProjectConfig(w http.ResponseWriter, r *http.Request) {
	var req createProjectConfigRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	p, err := s.svc.CreateProjectConfig(model.ProjectConfig{
		Key:         req.Key,
		Value:       req.Value,
		Description: req.Description,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, p)
}

func (s *Server) listProjectConfigs(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ProjectConfigFilter{
		Key:     r.URL.Query().Get("key"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListProjectConfigs(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getProjectConfig(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	p, err := s.svc.GetProjectConfig(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, p)
}

type updateProjectConfigRequest struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

func (s *Server) updateProjectConfig(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateProjectConfigRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	p, err := s.svc.UpdateProjectConfig(id, model.ProjectConfig{
		Key:         req.Key,
		Value:       req.Value,
		Description: req.Description,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, p)
}

func (s *Server) deleteProjectConfig(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteProjectConfig(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
