package store

import (
	"testing"

	"replay/internal/model"
)

func TestRecordingSessionStore(t *testing.T) {
	s := NewMemoryStore()
	sess := &model.RecordingSession{ID: "s1", Name: "test"}
	if err := s.CreateRecordingSession(sess); err != nil {
		t.Fatal(err)
	}
	if err := s.CreateRecordingSession(&model.RecordingSession{ID: "s2", Name: "test"}); err != ErrConflict {
		t.Fatalf("expected conflict, got %v", err)
	}
	got, err := s.GetRecordingSession("s1")
	if err != nil || got.Name != "test" {
		t.Fatal("get failed")
	}
	if _, err := s.GetRecordingSession("x"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	if len(s.ListRecordingSessions()) != 1 {
		t.Fatal("list count mismatch")
	}
	sess.Name = "updated"
	if err := s.UpdateRecordingSession(sess); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteRecordingSession("s1"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.GetRecordingSession("s1"); err != ErrNotFound {
		t.Fatal("expected not found after delete")
	}
}

func TestTrafficRecordStore(t *testing.T) {
	s := NewMemoryStore()
	rec := &model.TrafficRecord{ID: "t1", SessionID: "s1", Method: "GET", Path: "/"}
	if err := s.CreateTrafficRecord(rec); err != nil {
		t.Fatal(err)
	}
	got, err := s.GetTrafficRecord("t1")
	if err != nil || got.Method != "GET" {
		t.Fatal("get failed")
	}
	if _, err := s.GetTrafficRecord("x"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	rec.Method = "POST"
	if err := s.UpdateTrafficRecord(rec); err != nil {
		t.Fatal(err)
	}
	if len(s.ListTrafficRecords()) != 1 {
		t.Fatal("list count mismatch")
	}
	if err := s.DeleteTrafficRecord("t1"); err != nil {
		t.Fatal(err)
	}
}

func TestReplayTaskStore(t *testing.T) {
	s := NewMemoryStore()
	task := &model.ReplayTask{ID: "r1", SessionID: "s1", TargetEnvID: "e1"}
	if err := s.CreateReplayTask(task); err != nil {
		t.Fatal(err)
	}
	got, err := s.GetReplayTask("r1")
	if err != nil || got.SessionID != "s1" {
		t.Fatal("get failed")
	}
	if _, err := s.GetReplayTask("x"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	if len(s.ListReplayTasks()) != 1 {
		t.Fatal("list count mismatch")
	}
	if err := s.DeleteReplayTask("r1"); err != nil {
		t.Fatal(err)
	}
}

func TestTargetEnvStore(t *testing.T) {
	s := NewMemoryStore()
	e := &model.TargetEnv{ID: "e1", Name: "env1", BaseURL: "http://localhost"}
	if err := s.CreateTargetEnv(e); err != nil {
		t.Fatal(err)
	}
	if err := s.CreateTargetEnv(&model.TargetEnv{ID: "e2", Name: "env1"}); err != ErrConflict {
		t.Fatalf("expected conflict, got %v", err)
	}
	if _, err := s.GetTargetEnv("x"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	if err := s.DeleteTargetEnv("e1"); err != nil {
		t.Fatal(err)
	}
}

func TestFilterRuleStore(t *testing.T) {
	s := NewMemoryStore()
	f := &model.FilterRule{ID: "f1", SessionID: "s1", Field: "method", Operator: "eq", Value: "GET"}
	if err := s.CreateFilterRule(f); err != nil {
		t.Fatal(err)
	}
	if _, err := s.GetFilterRule("x"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	if len(s.ListFilterRules()) != 1 {
		t.Fatal("list count mismatch")
	}
	if err := s.DeleteFilterRule("f1"); err != nil {
		t.Fatal(err)
	}
}

func TestCompareReportStore(t *testing.T) {
	s := NewMemoryStore()
	rep := &model.CompareReport{ID: "c1", ReplayTaskID: "r1"}
	if err := s.CreateCompareReport(rep); err != nil {
		t.Fatal(err)
	}
	if _, err := s.GetCompareReport("x"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	if err := s.DeleteCompareReport("c1"); err != nil {
		t.Fatal(err)
	}
}

func TestDiffRecordStore(t *testing.T) {
	s := NewMemoryStore()
	d := &model.DiffRecord{ID: "d1", ReportID: "c1", TrafficRecordID: "t1"}
	if err := s.CreateDiffRecord(d); err != nil {
		t.Fatal(err)
	}
	if _, err := s.GetDiffRecord("x"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	if err := s.DeleteDiffRecord("d1"); err != nil {
		t.Fatal(err)
	}
}

func TestScheduleStore(t *testing.T) {
	s := NewMemoryStore()
	sch := &model.Schedule{ID: "sch1", ReplayTaskID: "r1", CronExpr: "0 0 * * *"}
	if err := s.CreateSchedule(sch); err != nil {
		t.Fatal(err)
	}
	if _, err := s.GetSchedule("x"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	if err := s.DeleteSchedule("sch1"); err != nil {
		t.Fatal(err)
	}
}

func TestReplayResultStore(t *testing.T) {
	s := NewMemoryStore()
	res := &model.ReplayResult{ID: "res1", ReplayTaskID: "r1", TrafficRecordID: "t1", Status: "success"}
	if err := s.CreateReplayResult(res); err != nil {
		t.Fatal(err)
	}
	if _, err := s.GetReplayResult("x"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	if err := s.DeleteReplayResult("res1"); err != nil {
		t.Fatal(err)
	}
}

func TestAuditLogStore(t *testing.T) {
	s := NewMemoryStore()
	a := &model.AuditLog{ID: "a1", Operator: "admin", Action: "create", TargetType: "task", TargetID: "r1"}
	if err := s.CreateAuditLog(a); err != nil {
		t.Fatal(err)
	}
	if _, err := s.GetAuditLog("x"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	if err := s.DeleteAuditLog("a1"); err != nil {
		t.Fatal(err)
	}
}
