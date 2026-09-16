package handler

import (
	"net/http"

	"replay/internal/model"
	"replay/pkg/httpx"
)

func (s *Server) registerDataSourceRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/data-sources", s.createDataSource)
	mux.HandleFunc("GET /api/data-sources", s.listDataSources)
	mux.HandleFunc("GET /api/data-sources/{id}", s.getDataSource)
	mux.HandleFunc("PUT /api/data-sources/{id}", s.updateDataSource)
	mux.HandleFunc("DELETE /api/data-sources/{id}", s.deleteDataSource)
}

type createDataSourceRequest struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Config string `json:"config"`
	Status string `json:"status"`
}

func (s *Server) createDataSource(w http.ResponseWriter, r *http.Request) {
	var req createDataSourceRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	d, err := s.svc.CreateDataSource(model.DataSource{
		Name:   req.Name,
		Type:   req.Type,
		Config: req.Config,
		Status: req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, d)
}

func (s *Server) listDataSources(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.DataSourceFilter{
		Name:    r.URL.Query().Get("name"),
		Type:    r.URL.Query().Get("type"),
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListDataSources(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getDataSource(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	d, err := s.svc.GetDataSource(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, d)
}

type updateDataSourceRequest struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Config string `json:"config"`
	Status string `json:"status"`
}

func (s *Server) updateDataSource(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateDataSourceRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	d, err := s.svc.UpdateDataSource(id, model.DataSource{
		Name:   req.Name,
		Type:   req.Type,
		Config: req.Config,
		Status: req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, d)
}

func (s *Server) deleteDataSource(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteDataSource(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
