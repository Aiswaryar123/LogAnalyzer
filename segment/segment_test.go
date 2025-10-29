package segment

import (
	"os"
	"path/filepath"
	"testing"
)

func createTempLogFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp log file : %v\n", err)
	}
	return path
}
func TestParseLogSegments(t *testing.T) {
	tempDir := t.TempDir()
	logContent := `2025-10-27 10:00:00.123 | INFO | parser | host=server1 | request_id=req1 | msg="Started"
2025-10-27 10:01:00.456 | ERROR | database | host=server2 | request_id=req2 | msg="Failed to connect"`
	createTempLogFile(t, tempDir, "sample.log", logContent)
	logStore, err := ParseLogSegments(tempDir)
	if err != nil {
		t.Fatalf("Expected no error but got %v\n", err)
	}
	if len(logStore.Segments) != 1 {
		t.Errorf("Expected 1 segment but got %d\n", len(logStore.Segments))
	}
	segment := logStore.Segments[0]
	if segment.Filename != "sample.log" {
		t.Errorf("Expected file name 'sample.log' but got %s\n", segment.Filename)
	}
	if len(segment.LogEntries) != 2 {
		t.Errorf("Expected 2 log entries but got %d\n", len(segment.LogEntries))
	}
	if len(segment.Index.ByLevel["INFO"]) == 0 || len(segment.Index.ByLevel["ERROR"]) == 0 {
		t.Errorf("Expected indexes for INFO and ERROR levels")
	}

	if len(segment.Index.ByComponent["parser"]) == 0 || len(segment.Index.ByComponent["database"]) == 0 {
		t.Errorf("Expected indexes for parser and database components")
	}

	if len(segment.Index.ByHost["server1"]) == 0 || len(segment.Index.ByHost["server2"]) == 0 {
		t.Errorf("Expected indexes for both hosts")
	}

	if len(segment.Index.ByRequestID["req1"]) == 0 || len(segment.Index.ByRequestID["req2"]) == 0 {
		t.Errorf("Expected indexes for both request IDs")
	}
}
func TestParseLogSegments_BadDir(t *testing.T) {
	_, err := ParseLogSegments("nonexistent_dir")
	if err == nil {
		t.Errorf("Expected error for nonexistent directory, got nil")
	}
}
func TestParseLogSegments_SkipSubdirectory(t *testing.T) {
	tmpDir := t.TempDir()
	os.Mkdir(filepath.Join(tmpDir, "subdir"), 0755)
	_, err := ParseLogSegments(tmpDir)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
func TestParseLogSegments_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "empty.log"), []byte("not a log line"), 0644)

	logstore, err := ParseLogSegments(tmpDir)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(logstore.Segments) != 0 {
		t.Errorf("Expected 0 segments, got %d", len(logstore.Segments))
	}
}
