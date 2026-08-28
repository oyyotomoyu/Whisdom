package logs

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func newTestService(t *testing.T) *Service {
	t.Helper()
	dir := t.TempDir()
	svc, err := NewService(dir)
	if err != nil {
		t.Fatalf("NewService: %v", err)
	}
	t.Cleanup(func() { svc.Close() })
	return svc
}

func TestLogRejectsInvalidStatus(t *testing.T) {
	svc := newTestService(t)
	if err := svc.Log("critical", "127.0.0.1", "usr_1", "should fail"); err == nil {
		t.Error("Log should reject an invalid status")
	}
}

func TestLogWritesJSONLine(t *testing.T) {
	svc := newTestService(t)
	if err := svc.Log(StatusInfo, "127.0.0.1", "usr_1", "hello"); err != nil {
		t.Fatalf("Log: %v", err)
	}

	entries, err := os.ReadDir(svc.dir)
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected exactly one log file, got %d", len(entries))
	}
	if !dailyFileName.MatchString(entries[0].Name()) {
		t.Errorf("log file name %q does not match daily pattern", entries[0].Name())
	}
}

func TestQueryFiltersAndPaginates(t *testing.T) {
	svc := newTestService(t)
	svc.Log(StatusInfo, "1.1.1.1", "usr_1", "logged in")
	svc.Log(StatusWarning, "1.1.1.1", "usr_1", "slow model response")
	svc.Log(StatusError, "1.1.1.1", "usr_2", "model request failed")

	result, err := svc.Query(QueryOptions{Statuses: []string{StatusWarning, StatusError}})
	if err != nil {
		t.Fatalf("Query: %v", err)
	}
	if len(result.Logs) != 2 {
		t.Fatalf("expected 2 filtered records, got %d", len(result.Logs))
	}

	result, err = svc.Query(QueryOptions{Keyword: "model"})
	if err != nil {
		t.Fatalf("Query: %v", err)
	}
	if len(result.Logs) != 2 {
		t.Fatalf("expected 2 keyword-matched records, got %d", len(result.Logs))
	}

	result, err = svc.Query(QueryOptions{Limit: 1})
	if err != nil {
		t.Fatalf("Query: %v", err)
	}
	if len(result.Logs) != 1 || result.NextCursor == "" {
		t.Fatalf("expected a paginated first page with a next cursor, got %d logs, cursor=%q", len(result.Logs), result.NextCursor)
	}
}

func TestCleanupRemovesOldFilesOnly(t *testing.T) {
	dir := t.TempDir()

	old := time.Now().AddDate(0, 0, -RetentionDays-1).Format("2006-01-02")
	recent := time.Now().Format("2006-01-02")

	writeEmptyLogFile(t, dir, old+".log")
	writeEmptyLogFile(t, dir, recent+".log")
	writeEmptyLogFile(t, dir, "not-a-log-file.txt")

	svc, err := NewService(dir) // NewService runs an initial cleanup
	if err != nil {
		t.Fatalf("NewService: %v", err)
	}
	defer svc.Close()

	if _, err := os.Stat(filepath.Join(dir, old+".log")); !os.IsNotExist(err) {
		t.Error("expected old daily log file to be removed")
	}
	if _, err := os.Stat(filepath.Join(dir, recent+".log")); err != nil {
		t.Error("expected recent daily log file to survive cleanup")
	}
	if _, err := os.Stat(filepath.Join(dir, "not-a-log-file.txt")); err != nil {
		t.Error("cleanup must not touch files outside the YYYY-MM-DD.log pattern")
	}
}

func writeEmptyLogFile(t *testing.T, dir, name string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte("{}\n"), 0o644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}
