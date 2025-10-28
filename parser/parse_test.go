package parser

import (
	models "miniproject/model"
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
