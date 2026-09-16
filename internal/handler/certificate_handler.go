package handler

import (
	"net/http"

	"replay/internal/model"
	"replay/pkg/httpx"
)

func (s *Server) registerCertificateRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/certificates", s.createCertificate)
	mux.HandleFunc("GET /api/certificates", s.listCertificates)
	mux.HandleFunc("GET /api/certificates/{id}", s.getCertificate)
	mux.HandleFunc("PUT /api/certificates/{id}", s.updateCertificate)
	mux.HandleFunc("DELETE /api/certificates/{id}", s.deleteCertificate)
}

type createCertificateRequest struct {
	Name      string `json:"name"`
	CertPEM   string `json:"cert_pem"`
	KeyPEM    string `json:"key_pem"`
	ExpiresAt string `json:"expires_at"`
}

func (s *Server) createCertificate(w http.ResponseWriter, r *http.Request) {
	var req createCertificateRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.CreateCertificate(model.Certificate{
		Name:    req.Name,
		CertPEM: req.CertPEM,
		KeyPEM:  req.KeyPEM,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, c)
}

func (s *Server) listCertificates(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.CertificateFilter{
		Name:    r.URL.Query().Get("name"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListCertificates(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getCertificate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	c, err := s.svc.GetCertificate(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

type updateCertificateRequest struct {
	Name    string `json:"name"`
	CertPEM string `json:"cert_pem"`
	KeyPEM  string `json:"key_pem"`
}

func (s *Server) updateCertificate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateCertificateRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.UpdateCertificate(id, model.Certificate{
		Name:    req.Name,
		CertPEM: req.CertPEM,
		KeyPEM:  req.KeyPEM,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

func (s *Server) deleteCertificate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteCertificate(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
