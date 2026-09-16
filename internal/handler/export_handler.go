package handler

import (
	"net/http"

	"replay/pkg/httpx"
)

func (s *Server) registerExportRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/export/snapshot", s.exportSnapshot)
}

func (s *Server) exportSnapshot(w http.ResponseWriter, r *http.Request) {
	snapshot := s.svc.GetSnapshotExport()
	httpx.OK(w, snapshot)
}
