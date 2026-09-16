package handler

import (
	"net/http"

	"replay/internal/model"
	"replay/pkg/httpx"
)

func (s *Server) registerScriptTemplateRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/script-templates", s.createScriptTemplate)
	mux.HandleFunc("GET /api/script-templates", s.listScriptTemplates)
	mux.HandleFunc("GET /api/script-templates/{id}", s.getScriptTemplate)
	mux.HandleFunc("PUT /api/script-templates/{id}", s.updateScriptTemplate)
	mux.HandleFunc("DELETE /api/script-templates/{id}", s.deleteScriptTemplate)
}

type createScriptTemplateRequest struct {
	Name        string `json:"name"`
	ScriptType  string `json:"script_type"`
	Content     string `json:"content"`
	Description string `json:"description"`
}

func (s *Server) createScriptTemplate(w http.ResponseWriter, r *http.Request) {
	var req createScriptTemplateRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	st, err := s.svc.CreateScriptTemplate(model.ScriptTemplate{
		Name:        req.Name,
		ScriptType:  req.ScriptType,
		Content:     req.Content,
		Description: req.Description,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, st)
}

func (s *Server) listScriptTemplates(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ScriptTemplateFilter{
		Name:       r.URL.Query().Get("name"),
		ScriptType: r.URL.Query().Get("script_type"),
		Keyword:    r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListScriptTemplates(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getScriptTemplate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	st, err := s.svc.GetScriptTemplate(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, st)
}

type updateScriptTemplateRequest struct {
	Name        string `json:"name"`
	ScriptType  string `json:"script_type"`
	Content     string `json:"content"`
	Description string `json:"description"`
}

func (s *Server) updateScriptTemplate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateScriptTemplateRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	st, err := s.svc.UpdateScriptTemplate(id, model.ScriptTemplate{
		Name:        req.Name,
		ScriptType:  req.ScriptType,
		Content:     req.Content,
		Description: req.Description,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, st)
}

func (s *Server) deleteScriptTemplate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteScriptTemplate(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
