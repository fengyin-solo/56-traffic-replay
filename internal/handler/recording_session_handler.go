package handler

import (
	"net/http"

	"replay/internal/model"
	"replay/pkg/httpx"
)

func (s *Server) registerRecordingSessionRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/recording-sessions", s.createRecordingSession)
	mux.HandleFunc("GET /api/recording-sessions", s.listRecordingSessions)
	mux.HandleFunc("GET /api/recording-sessions/{id}", s.getRecordingSession)
	mux.HandleFunc("PUT /api/recording-sessions/{id}", s.updateRecordingSession)
	mux.HandleFunc("DELETE /api/recording-sessions/{id}", s.deleteRecordingSession)
	mux.HandleFunc("POST /api/recording-sessions/{id}/start", s.startRecording)
	mux.HandleFunc("POST /api/recording-sessions/{id}/pause", s.pauseRecording)
	mux.HandleFunc("POST /api/recording-sessions/{id}/stop", s.stopRecording)
}

type createRecordingSessionRequest struct {
	Name       string `json:"name"`
	Source     string `json:"source"`
	TargetHost string `json:"target_host"`
	Protocol   string `json:"protocol"`
}

func (s *Server) createRecordingSession(w http.ResponseWriter, r *http.Request) {
	var req createRecordingSessionRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sess, err := s.svc.CreateRecordingSession(model.RecordingSession{
		Name:       req.Name,
		Source:     req.Source,
		TargetHost: req.TargetHost,
		Protocol:   req.Protocol,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, sess)
}

func (s *Server) listRecordingSessions(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.RecordingSessionFilter{
		Name:     r.URL.Query().Get("name"),
		Protocol: r.URL.Query().Get("protocol"),
		Status:   r.URL.Query().Get("status"),
		Keyword:  r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListRecordingSessions(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getRecordingSession(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	sess, err := s.svc.GetRecordingSession(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sess)
}

type updateRecordingSessionRequest struct {
	Name       string `json:"name"`
	Source     string `json:"source"`
	TargetHost string `json:"target_host"`
	Protocol   string `json:"protocol"`
}

func (s *Server) updateRecordingSession(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateRecordingSessionRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sess, err := s.svc.UpdateRecordingSession(id, model.RecordingSession{
		Name:       req.Name,
		Source:     req.Source,
		TargetHost: req.TargetHost,
		Protocol:   req.Protocol,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sess)
}

func (s *Server) deleteRecordingSession(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteRecordingSession(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) startRecording(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	sess, err := s.svc.StartRecording(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sess)
}

func (s *Server) pauseRecording(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	sess, err := s.svc.PauseRecording(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sess)
}

func (s *Server) stopRecording(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	sess, err := s.svc.StopRecording(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sess)
}
