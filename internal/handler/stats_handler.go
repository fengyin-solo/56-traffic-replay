package handler

import (
	"net/http"

	"replay/pkg/httpx"
)

func (s *Server) registerStatsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/stats/overview", s.getOverviewStats)
	mux.HandleFunc("GET /api/stats/method-distribution", s.getMethodDistribution)
	mux.HandleFunc("GET /api/stats/env-distribution", s.getEnvDistribution)
	mux.HandleFunc("GET /api/stats/endpoint-distribution", s.getEndpointDistribution)
	mux.HandleFunc("GET /api/stats/success-rate-trend", s.getSuccessRateTrend)
}

func (s *Server) getOverviewStats(w http.ResponseWriter, r *http.Request) {
	stats := s.svc.GetOverviewStats()
	httpx.OK(w, stats)
}

func (s *Server) getMethodDistribution(w http.ResponseWriter, r *http.Request) {
	dist := s.svc.GetMethodDistribution()
	httpx.OK(w, dist)
}

func (s *Server) getEnvDistribution(w http.ResponseWriter, r *http.Request) {
	dist := s.svc.GetEnvDistribution()
	httpx.OK(w, dist)
}

func (s *Server) getEndpointDistribution(w http.ResponseWriter, r *http.Request) {
	dist := s.svc.GetEndpointDistribution()
	httpx.OK(w, dist)
}

func (s *Server) getSuccessRateTrend(w http.ResponseWriter, r *http.Request) {
	trend := s.svc.GetSuccessRateTrend()
	httpx.OK(w, trend)
}
