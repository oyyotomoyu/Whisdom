package logs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"
)

// RetentionDays is the maximum number of days a daily log file is kept.
const RetentionDays = 7

var dailyFileName = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}\.log$`)

// Service writes daily JSON Lines log files and enforces retention. It is
// safe for concurrent use.
type Service struct {
	dir string

	mu      sync.Mutex
	file    *os.File
	dateKey string
}

// NewService creates the log directory if needed, runs an initial retention
// cleanup, and returns a ready-to-use Service.
func NewService(dir string) (*Service, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("create log dir: %w", err)
	}
	s := &Service{dir: dir}
	if err := s.Cleanup(); err != nil {
		return nil, err
	}
	return s, nil
}

// Log validates status and appends one record to today's log file. IP and
// userID may be empty for background/system events.
func (s *Service) Log(status, ip, userID, content string) error {
	if err := validateStatus(status); err != nil {
		return err
	}

	now := time.Now()
	record := Record{
		Timestamp: now.Format(timestampLayout),
		Status:    status,
		IP:        ip,
		UserID:    userID,
		Content:   content,
	}
	line, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("encode log record: %w", err)
	}
	line = append(line, '\n')

	s.mu.Lock()
	defer s.mu.Unlock()

	f, err := s.fileForDateLocked(now)
	if err != nil {
		return err
	}
	if _, err := f.Write(line); err != nil {
		return fmt.Errorf("write log record: %w", err)
	}
	return nil
}

// fileForDateLocked returns the open file handle for the log file matching
// t's date, rotating to a new day's file (and running cleanup first) if the
// date has changed since the last write. Caller must hold s.mu.
func (s *Service) fileForDateLocked(t time.Time) (*os.File, error) {
	dateKey := t.Format("2006-01-02")
	if s.file != nil && s.dateKey == dateKey {
		return s.file, nil
	}

	if s.file != nil {
		_ = s.file.Close()
		s.file = nil
	}

	if err := s.cleanupLocked(); err != nil {
		return nil, err
	}

	path := filepath.Join(s.dir, dateKey+".log")
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}

	s.file = f
	s.dateKey = dateKey
	return f, nil
}

// Cleanup removes daily log files older than RetentionDays. Safe to call on
// a schedule (startup and hourly).
func (s *Service) Cleanup() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.cleanupLocked()
}

func (s *Service) cleanupLocked() error {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return fmt.Errorf("read log dir: %w", err)
	}

	cutoff := time.Now().AddDate(0, 0, -RetentionDays)

	for _, entry := range entries {
		if entry.IsDir() || !dailyFileName.MatchString(entry.Name()) {
			continue
		}
		dateKey := entry.Name()[:len("2006-01-02")]
		fileDate, err := time.ParseInLocation("2006-01-02", dateKey, time.Local)
		if err != nil {
			continue
		}
		if fileDate.Before(cutoff) {
			_ = os.Remove(filepath.Join(s.dir, entry.Name()))
		}
	}
	return nil
}

// StartRetentionLoop runs Cleanup once immediately and then every hour until
// stop is closed.
func (s *Service) StartRetentionLoop(stop <-chan struct{}) {
	ticker := time.NewTicker(time.Hour)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				_ = s.Cleanup()
			}
		}
	}()
}

// Close closes the currently open log file, if any.
func (s *Service) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.file == nil {
		return nil
	}
	err := s.file.Close()
	s.file = nil
	return err
}
