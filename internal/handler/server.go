// Package handler 实现 HTTP 处理器层。
package handler

import (
	"errors"
	"net"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	"replay/internal/config"
	"replay/internal/model"
	"replay/internal/service"
	"replay/internal/store"
	"replay/pkg/httpx"
	"replay/pkg/logger"
)

type Server struct {
	svc *service.Service
	log *logger.Logger
	cfg *config.Config
}

func NewServer(svc *service.Service, log *logger.Logger, cfg *config.Config) *Server {
	return &Server{svc: svc, log: log, cfg: cfg}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	s.registerRecordingSessionRoutes(mux)
	s.registerTrafficRecordRoutes(mux)
	s.registerReplayTaskRoutes(mux)
	s.registerTargetEnvRoutes(mux)
	s.registerFilterRuleRoutes(mux)
	s.registerCompareReportRoutes(mux)
	s.registerDiffRecordRoutes(mux)
	s.registerScheduleRoutes(mux)
	s.registerReplayResultRoutes(mux)
	s.registerAuditLogRoutes(mux)
	s.registerStatsRoutes(mux)
	s.registerExportRoutes(mux)
	s.registerWebhookConfigRoutes(mux)
	s.registerScriptTemplateRoutes(mux)
	s.registerEnvVariableRoutes(mux)
	s.registerNotificationRuleRoutes(mux)
	s.registerProjectConfigRoutes(mux)
	s.registerTagRoutes(mux)
	s.registerCertificateRoutes(mux)
	s.registerDataSourceRoutes(mux)
	mux.Handle("GET /", http.FileServer(http.Dir("web")))
	return s.loggingMiddleware(s.recoveryMiddleware(s.apiKeyMiddleware(s.rateLimitMiddleware(mux))))
}

func (s *Server) maxPageSize() int {
	if s.cfg != nil && s.cfg.MaxPageSize > 0 {
		return s.cfg.MaxPageSize
	}
	return 100
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		s.log.Infof("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func (s *Server) recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				s.log.Errorf("panic: %v\n%s", rec, debug.Stack())
				httpx.InternalError(w, "服务器内部错误")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (s *Server) apiKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.cfg.APIKey == "" {
			next.ServeHTTP(w, r)
			return
		}
		key := r.Header.Get("X-API-Key")
		if key == "" {
			key = r.URL.Query().Get("api_key")
		}
		if key != s.cfg.APIKey {
			httpx.Unauthorized(w, "API Key 无效")
			return
		}
		next.ServeHTTP(w, r)
	})
}

var (
	rateLimiter = make(map[string][]time.Time)
	rateMu      sync.Mutex
)

func (s *Server) rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			host = r.RemoteAddr
		}
		now := time.Now()
		rateMu.Lock()
		var times []time.Time
		for _, t := range rateLimiter[host] {
			if now.Sub(t) < time.Minute {
				times = append(times, t)
			}
		}
		times = append(times, now)
		rateLimiter[host] = times
		count := len(times)
		rateMu.Unlock()
		if count > 120 {
			httpx.Error(w, http.StatusTooManyRequests, 429, "请求过于频繁")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case model.IsValidationError(err):
		httpx.BadRequest(w, err.Error())
	case errors.Is(err, store.ErrNotFound):
		httpx.NotFound(w, err.Error())
	case errors.Is(err, store.ErrConflict):
		httpx.Conflict(w, err.Error())
	default:
		httpx.InternalError(w, err.Error())
	}
}
