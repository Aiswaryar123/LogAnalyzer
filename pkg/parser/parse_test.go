package parser

import (
	models "Log_analyzer/model"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLogParseEntryLine(t *testing.T) {
	line := `2025-10-24 12:34:56.789 | INFO | api-server | host=web01 | request_id=req-ab12cd-1234 | msg="Request started GET /api"`

	expectedTime, _ := time.Parse("2006-01-02 15:04:05.000", "2025-10-24 12:34:56.789")

	expected := &models.LogEntry{
		Time:      expectedTime,
		Level:     "INFO",
		Component: "api-server",
		Host:      "web01",
		Requestid: "req-ab12cd-1234",
		Message:   "Request started GET /api",
	}

	got, err := LogParseEntry(line)
	if err != nil {
		t.Errorf("log parsing failed")
	}

	if !got.Time.Equal(expected.Time) {
		t.Errorf("Expected time %v, got %v", expected.Time, got.Time)
	}
	if got.Level != expected.Level {
		t.Errorf("Expected Level %s, got %s", expected.Level, got.Level)
	}
	if got.Component != expected.Component {
		t.Errorf("Expected Component %s, got %s", expected.Component, got.Component)
	}
	if got.Host != expected.Host {
		t.Errorf("Expected Host %s, got %s", expected.Host, got.Host)
	}
	if got.Requestid != expected.Requestid {
		t.Errorf("Expected Requestid %s, got %s", expected.Requestid, got.Requestid)
	}
	if got.Message != expected.Message {
		t.Errorf("Expected Message %s, got %s", expected.Message, got.Message)
	}
}
func TestParseInvalidLogEntry(t *testing.T) {
	invalidLine := `invalid log line`
	_, err := LogParseEntry(invalidLine)
	if err == nil {
		t.Errorf("Expected error for invalid format but got none")
	}
}
func TestParseLogEntryBadTime(t *testing.T) {
	badTimeLine := `2025-10-23 15:17:08.636000 | WARN | api-server | host=worker01 | request_id=req-4leuyy-5910 | msg="Cache cleared"`
	_, err := LogParseEntry(badTimeLine)
	if err == nil || !strings.Contains(err.Error(), "failed to parse time") {
		t.Errorf("Expected time parsing error, got %v", err)
	}
}
func TestParseLogFilesBadDirectory(t *testing.T) {
	path := "../logss"
	_, err := ParseLogFiles(path)
	if err == nil {
		t.Errorf("Expected 'no such directory' error but got none.")
	}
}
func TestParseLogFilesValidFile(t *testing.T) {
	tmpDir := t.TempDir()
	logContent := `2025-10-23 15:17:08.636 | INFO | api-server | host=worker01 | request_id=req-xyz | msg="Cache cleared"`
	tmpFile := filepath.Join(tmpDir, "valid.log")

	err := os.WriteFile(tmpFile, []byte(logContent+"\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp log file: %v", err)
	}

	entries, err := ParseLogFiles(tmpDir)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(entries) != 1 {
		t.Errorf("Expected 1 log entry, got %d", len(entries))
	}
	if entries[0].Host != "worker01" {
		t.Errorf("Expected host=worker01, got %s", entries[0].Host)
	}
}

func TestParseLogFilesInvalidLog(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "invalid.log")

	err := os.WriteFile(tmpFile, []byte("invalid log line\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp invalid log file: %v", err)
	}

	entries, err := ParseLogFiles(tmpDir)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(entries) != 0 {
		t.Errorf("Expected 0 valid entries, got %d", len(entries))
	}
}
func TestParseLogFilesUnreadableFile(t *testing.T) {
	tmpDir := t.TempDir()
	badFile := filepath.Join(tmpDir, "unreadable.log")

	err := os.WriteFile(badFile, []byte("data"), 0000)
	if err != nil {
		t.Fatalf("Failed to create unreadable file: %v", err)
	}
	defer os.Chmod(badFile, 0644)

	_, err = ParseLogFiles(tmpDir)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}
func TestParseLogFilesSkipsSubdir(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "nested")
	os.Mkdir(subDir, 0755)

	logContent := `2025-10-23 15:17:08.636 | WARN | scheduler | host=worker02 | request_id=req-abc | msg="Job delayed"`
	tmpFile := filepath.Join(tmpDir, "log1.log")
	os.WriteFile(tmpFile, []byte(logContent+"\n"), 0644)

	entries, err := ParseLogFiles(tmpDir)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(entries) != 1 {
		t.Errorf("Expected 1 entry (subdir skipped), got %d", len(entries))
	}
}
