package service

import (
	"testing"

	"replay/internal/config"
	"replay/internal/model"
	"replay/internal/store"
	"replay/pkg/logger"
)

func newTestService() *Service {
	cfg := &config.Config{MaxPageSize: 100}
	log := logger.NewLevel(logger.LevelError)
	return New(store.NewMemoryStore(), log, cfg)
}

func TestCreateRecordingSession(t *testing.T) {
	svc := newTestService()
	sess, err := svc.CreateRecordingSession(model.RecordingSession{Name: "sess1", Source: "src", TargetHost: "host", Protocol: "http"})
	if err != nil {
		t.Fatal(err)
	}
	if sess.Status != model.SessionRecording {
		t.Fatal("status should be recording")
	}
	_, err = svc.CreateRecordingSession(model.RecordingSession{Name: "sess1", Source: "src", TargetHost: "host"})
	if err == nil {
		t.Fatal("expected conflict")
	}
}

func TestRecordingSessionStateMachine(t *testing.T) {
	svc := newTestService()
	sess, _ := svc.CreateRecordingSession(model.RecordingSession{Name: "sm", Source: "s", TargetHost: "h"})
	if _, err := svc.StopRecording(sess.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.StartRecording(sess.ID); err == nil {
		t.Fatal("should not allow start after completed")
	}
}

func TestCreateTrafficRecord(t *testing.T) {
	svc := newTestService()
	_, err := svc.CreateTrafficRecord(model.TrafficRecord{SessionID: "x", Method: "GET", Path: "/"})
	if err == nil {
		t.Fatal("expected session not found")
	}
	sess, _ := svc.CreateRecordingSession(model.RecordingSession{Name: "tr", Source: "s", TargetHost: "h"})
	rec, err := svc.CreateTrafficRecord(model.TrafficRecord{SessionID: sess.ID, Method: "GET", Path: "/"})
	if err != nil {
		t.Fatal(err)
	}
	if rec.SessionID != sess.ID {
		t.Fatal("session id mismatch")
	}
}

func TestImportTrafficRecords(t *testing.T) {
	svc := newTestService()
	sess, _ := svc.CreateRecordingSession(model.RecordingSession{Name: "imp", Source: "s", TargetHost: "h"})
	count, err := svc.ImportTrafficRecords(sess.ID, []model.TrafficRecord{
		{Method: "GET", Path: "/a"},
		{Method: "POST", Path: "/b"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Fatalf("expected 2 imported, got %d", count)
	}
}

func TestCreateReplayTask(t *testing.T) {
	svc := newTestService()
	sess, _ := svc.CreateRecordingSession(model.RecordingSession{Name: "rt", Source: "s", TargetHost: "h"})
	env, _ := svc.CreateTargetEnv(model.TargetEnv{Name: "e1", BaseURL: "http://localhost"})
	task, err := svc.CreateReplayTask(model.ReplayTask{SessionID: sess.ID, TargetEnvID: env.ID})
	if err != nil {
		t.Fatal(err)
	}
	if task.Status != model.TaskPending {
		t.Fatal("status should be pending")
	}
}

func TestReplayTaskStateMachine(t *testing.T) {
	svc := newTestService()
	sess, _ := svc.CreateRecordingSession(model.RecordingSession{Name: "rtsm", Source: "s", TargetHost: "h"})
	env, _ := svc.CreateTargetEnv(model.TargetEnv{Name: "e2", BaseURL: "http://localhost"})
	task, _ := svc.CreateReplayTask(model.ReplayTask{SessionID: sess.ID, TargetEnvID: env.ID})
	if _, err := svc.CompleteReplayTask(task.ID); err == nil {
		t.Fatal("should not complete from pending")
	}
	if _, err := svc.RunReplayTask(task.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.CompleteReplayTask(task.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.FailReplayTask(task.ID); err == nil {
		t.Fatal("should not fail after completed")
	}
}

func TestRunReplayTaskWithFilter(t *testing.T) {
	svc := newTestService()
	sess, _ := svc.CreateRecordingSession(model.RecordingSession{Name: "fr", Source: "s", TargetHost: "h"})
	env, _ := svc.CreateTargetEnv(model.TargetEnv{Name: "e3", BaseURL: "http://localhost"})
	_, _ = svc.CreateTrafficRecord(model.TrafficRecord{SessionID: sess.ID, Method: "GET", Path: "/api"})
	_, _ = svc.CreateTrafficRecord(model.TrafficRecord{SessionID: sess.ID, Method: "POST", Path: "/api"})
	_, _ = svc.CreateFilterRule(model.FilterRule{SessionID: sess.ID, Field: "method", Operator: "eq", Value: "GET", Action: "include"})
	task, _ := svc.CreateReplayTask(model.ReplayTask{SessionID: sess.ID, TargetEnvID: env.ID})
	result, err := svc.RunReplayTask(task.ID)
	if err != nil {
		t.Fatal(err)
	}
	if result.TotalRequests != 1 {
		t.Fatalf("expected 1 request after filter, got %d", result.TotalRequests)
	}
}

func TestTargetEnvConflict(t *testing.T) {
	svc := newTestService()
	_, _ = svc.CreateTargetEnv(model.TargetEnv{Name: "ec", BaseURL: "http://a"})
	_, err := svc.CreateTargetEnv(model.TargetEnv{Name: "ec", BaseURL: "http://b"})
	if err == nil {
		t.Fatal("expected conflict")
	}
}

func TestCompareReportRun(t *testing.T) {
	svc := newTestService()
	sess, _ := svc.CreateRecordingSession(model.RecordingSession{Name: "cr", Source: "s", TargetHost: "h"})
	env, _ := svc.CreateTargetEnv(model.TargetEnv{Name: "e4", BaseURL: "http://localhost"})
	_, _ = svc.CreateTrafficRecord(model.TrafficRecord{SessionID: sess.ID, Method: "GET", Path: "/", StatusCode: 200})
	task, _ := svc.CreateReplayTask(model.ReplayTask{SessionID: sess.ID, TargetEnvID: env.ID})
	_, _ = svc.RunReplayTask(task.ID)
	_, _ = svc.CompleteReplayTask(task.ID)
	report, _ := svc.CreateCompareReport(model.CompareReport{ReplayTaskID: task.ID})
	report, err := svc.RunCompareReport(report.ID)
	if err != nil {
		t.Fatal(err)
	}
	if report.Status != model.ReportCompleted {
		t.Fatal("report should be completed")
	}
}

func TestFilterRuleMatch(t *testing.T) {
	rule := &model.FilterRule{Field: "method", Operator: "eq", Value: "GET", Action: "include", Status: model.FilterStatusActive}
	rec := &model.TrafficRecord{Method: "GET", Path: "/"}
	if !rule.MatchRecord(rec) {
		t.Fatal("should match GET")
	}
	rec.Method = "POST"
	if rule.MatchRecord(rec) {
		t.Fatal("should not match POST")
	}
}

func TestAuditLog(t *testing.T) {
	svc := newTestService()
	a, err := svc.CreateAuditLog(model.AuditLog{Operator: "admin", Action: "create", TargetType: "task", TargetID: "t1"})
	if err != nil {
		t.Fatal(err)
	}
	if a.Operator != "admin" {
		t.Fatal("operator mismatch")
	}
}

func TestStats(t *testing.T) {
	svc := newTestService()
	overview := svc.GetOverviewStats()
	if overview.TotalSessions != 0 {
		t.Fatal("expected 0 sessions")
	}
}
