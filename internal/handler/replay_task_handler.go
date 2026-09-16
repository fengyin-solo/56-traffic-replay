package handler

import (
	"net/http"

	"replay/internal/model"
	"replay/pkg/httpx"
)

func (s *Server) registerReplayTaskRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/replay-tasks", s.createReplayTask)
	mux.HandleFunc("GET /api/replay-tasks", s.listReplayTasks)
	mux.HandleFunc("GET /api/replay-tasks/{id}", s.getReplayTask)
	mux.HandleFunc("PUT /api/replay-tasks/{id}", s.updateReplayTask)
	mux.HandleFunc("DELETE /api/replay-tasks/{id}", s.deleteReplayTask)
	mux.HandleFunc("POST /api/replay-tasks/{id}/run", s.runReplayTask)
	mux.HandleFunc("POST /api/replay-tasks/{id}/complete", s.completeReplayTask)
	mux.HandleFunc("POST /api/replay-tasks/{id}/fail", s.failReplayTask)
	mux.HandleFunc("POST /api/replay-tasks/batch", s.batchCreateReplayTasks)
}

type createReplayTaskRequest struct {
	SessionID   string `json:"session_id"`
	TargetEnvID string `json:"target_env_id"`
}

func (s *Server) createReplayTask(w http.ResponseWriter, r *http.Request) {
	var req createReplayTaskRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	task, err := s.svc.CreateReplayTask(model.ReplayTask{
		SessionID:   req.SessionID,
		TargetEnvID: req.TargetEnvID,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, task)
}

func (s *Server) listReplayTasks(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ReplayTaskFilter{
		SessionID:   r.URL.Query().Get("session_id"),
		TargetEnvID: r.URL.Query().Get("target_env_id"),
		Status:      r.URL.Query().Get("status"),
		Keyword:     r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListReplayTasks(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getReplayTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	task, err := s.svc.GetReplayTask(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, task)
}

type updateReplayTaskRequest struct {
	TargetEnvID string `json:"target_env_id"`
}

func (s *Server) updateReplayTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateReplayTaskRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	task, err := s.svc.UpdateReplayTask(id, model.ReplayTask{TargetEnvID: req.TargetEnvID})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, task)
}

func (s *Server) deleteReplayTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteReplayTask(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) runReplayTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	task, err := s.svc.RunReplayTask(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, task)
}

func (s *Server) completeReplayTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	task, err := s.svc.CompleteReplayTask(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, task)
}

func (s *Server) failReplayTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	task, err := s.svc.FailReplayTask(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, task)
}

type batchCreateReplayTasksRequest struct {
	Tasks []createReplayTaskRequest `json:"tasks"`
}

func (s *Server) batchCreateReplayTasks(w http.ResponseWriter, r *http.Request) {
	var req batchCreateReplayTasksRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	inputs := make([]model.ReplayTask, 0, len(req.Tasks))
	for _, t := range req.Tasks {
		inputs = append(inputs, model.ReplayTask{
			SessionID:   t.SessionID,
			TargetEnvID: t.TargetEnvID,
		})
	}
	tasks, err := s.svc.BatchCreateReplayTasks(inputs)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, tasks)
}
