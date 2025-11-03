package filter

import (
	models "miniproject/model"
	"testing"
	"time"
)

func makeTime(s string) time.Time {
	t, _ := time.Parse("2006-01-02 15:04:05", s)
	return t
}

func makeMockLogStore() models.LogStore {
	entries := []models.LogEntry{
		{Raw: "entry1", Level: "INFO", Component: "auth", Host: "server1", Requestid: "r1", Time: makeTime("2025-01-01 10:00:00")},
		{Raw: "entry2", Level: "ERROR", Component: "db", Host: "server2", Requestid: "r2", Time: makeTime("2025-01-01 11:00:00")},
		{Raw: "entry3", Level: "INFO", Component: "api", Host: "server1", Requestid: "r3", Time: makeTime("2025-01-01 12:00:00")},
	}

	index := models.SegmentIndex{
		ByLevel: map[string][]int{
			"INFO":  {0, 2},
			"ERROR": {1},
		},
		ByComponent: map[string][]int{
			"auth": {0},
			"db":   {1},
			"api":  {2},
		},
		ByHost: map[string][]int{
			"server1": {0, 2},
			"server2": {1},
		},
		ByRequestID: map[string][]int{
			"r1": {0},
			"r2": {1},
			"r3": {2},
		},
	}

	segment := models.Segment{
		LogEntries: entries,
		Starttime:  makeTime("2025-01-01 09:00:00"),
		Endtime:    makeTime("2025-01-01 13:00:00"),
		Index:      index,
	}

	return models.LogStore{Segments: []models.Segment{segment}}
}
func TestNoFiltersReturnsAllLogs(t *testing.T) {
	store := makeMockLogStore()

	got := FilterLogs(store, nil, nil, nil, nil, time.Time{}, time.Time{})
	if len(got) != 3 {
		t.Errorf("expected 3 logs, got %d", len(got))
	}
}
func TestFilterByLevel(t *testing.T) {
	store := makeMockLogStore()

	got := FilterLogs(store, []string{"ERROR"}, nil, nil, nil, time.Time{}, time.Time{})
	if len(got) != 1 || got[0].Raw != "entry2" {
		t.Errorf("expected entry2, got %+v", got)
	}
}
func TestFilterByComponent(t *testing.T) {
	store := makeMockLogStore()

	got := FilterLogs(store, nil, []string{"auth"}, nil, nil, time.Time{}, time.Time{})
	if len(got) != 1 || got[0].Raw != "entry1" {
		t.Errorf("expected entry1, got %+v", got)
	}
}
func TestFilterByHost(t *testing.T) {
	store := makeMockLogStore()

	got := FilterLogs(store, nil, nil, []string{"server2"}, nil, time.Time{}, time.Time{})
	if len(got) != 1 || got[0].Raw != "entry2" {
		t.Errorf("expected entry2, got %+v", got)
	}
}
func TestFilterByReqID(t *testing.T) {
	store := makeMockLogStore()

	got := FilterLogs(store, nil, nil, nil, []string{"r3"}, time.Time{}, time.Time{})
	if len(got) != 1 || got[0].Raw != "entry3" {
		t.Errorf("expected entry3, got %+v", got)
	}
}
func TestFilterByTimeRange(t *testing.T) {
	store := makeMockLogStore()

	start := makeTime("2025-01-01 10:30:00")
	end := makeTime("2025-01-01 12:30:00")
	got := FilterLogs(store, nil, nil, nil, nil, start, end)

	if len(got) != 2 {
		t.Errorf("expected 2 logs between 10:30 and 12:30, got %d", len(got))
	}
}
func TestCombinedLevelAndComponent(t *testing.T) {
	store := makeMockLogStore()

	got := FilterLogs(store, []string{"INFO"}, []string{"api"}, nil, nil, time.Time{}, time.Time{})
	if len(got) != 1 || got[0].Raw != "entry3" {
		t.Errorf("expected entry3, got %+v", got)
	}
}
func TestEmptyStore(t *testing.T) {
	store := models.LogStore{}
	got := FilterLogs(store, nil, nil, nil, nil, time.Time{}, time.Time{})
	if len(got) != 0 {
		t.Errorf("expected 0 logs for empty store, got %d", len(got))
	}
}
func TestSkipSegmentBeforeStartTime(t *testing.T) {
	now := time.Now()
	store := models.LogStore{
		Segments: []models.Segment{
			{
				Starttime:  now.Add(-2 * time.Hour),
				Endtime:    now.Add(-1 * time.Hour),
				LogEntries: []models.LogEntry{},
			},
		},
	}
	startTime := now
	endTime := time.Time{}

	result := FilterLogs(store, nil, nil, nil, nil, startTime, endTime)
	if len(result) != 0 {
		t.Errorf("expected 0 logs, got %d", len(result))
	}
}
func TestSkipSegmentAfterEndTime(t *testing.T) {
	now := time.Now()
	store := models.LogStore{
		Segments: []models.Segment{
			{
				Starttime:  now.Add(2 * time.Hour),
				Endtime:    now.Add(3 * time.Hour),
				LogEntries: []models.LogEntry{},
			},
		},
	}
	endTime := now
	startTime := time.Time{}

	result := FilterLogs(store, nil, nil, nil, nil, startTime, endTime)
	if len(result) != 0 {
		t.Errorf("expected 0 logs, got %d", len(result))
	}
}
func TestSkipEntryBeforeStartTime(t *testing.T) {
	now := time.Now()
	entry := models.LogEntry{Time: now.Add(-2 * time.Hour)}
	segment := models.Segment{
		Starttime:  now.Add(-3 * time.Hour),
		Endtime:    now.Add(-1 * time.Hour),
		LogEntries: []models.LogEntry{entry},
	}
	store := models.LogStore{Segments: []models.Segment{segment}}

	startTime := now
	endTime := time.Time{}

	result := FilterLogs(store, nil, nil, nil, nil, startTime, endTime)
	if len(result) != 0 {
		t.Errorf("expected 0 logs, got %d", len(result))
	}
}
func TestSkipEntryAfterEndTime(t *testing.T) {
	now := time.Now()
	entry := models.LogEntry{Time: now.Add(2 * time.Hour)}
	segment := models.Segment{
		Starttime:  now,
		Endtime:    now.Add(3 * time.Hour),
		LogEntries: []models.LogEntry{entry},
	}
	store := models.LogStore{Segments: []models.Segment{segment}}

	startTime := time.Time{}
	endTime := now

	result := FilterLogs(store, nil, nil, nil, nil, startTime, endTime)
	if len(result) != 0 {
		t.Errorf("expected 0 logs, got %d", len(result))
	}
}
func TestMatchedEntryBeforeStartTime(t *testing.T) {
	now := time.Now()

	entry1 := models.LogEntry{
		Time: now.Add(-2 * time.Hour),
		Raw:  "entry-before-start",
	}
	segment := models.Segment{
		Starttime:  now.Add(-3 * time.Hour),
		Endtime:    now,
		LogEntries: []models.LogEntry{entry1},
		Index: models.SegmentIndex{
			ByLevel: map[string][]int{"INFO": {0}},
		},
	}
	store := models.LogStore{Segments: []models.Segment{segment}}

	startTime := now.Add(-1 * time.Hour)
	endTime := time.Time{}

	result := FilterLogs(store, []string{"INFO"}, nil, nil, nil, startTime, endTime)
	if len(result) != 0 {
		t.Errorf("expected 0 logs (skipped due to before start), got %d", len(result))
	}
}
func TestMatchedEntryAfterEndTime(t *testing.T) {
	now := time.Now()

	entry1 := models.LogEntry{
		Time: now.Add(2 * time.Hour),
		Raw:  "entry-after-end",
	}
	segment := models.Segment{
		Starttime:  now,
		Endtime:    now.Add(3 * time.Hour),
		LogEntries: []models.LogEntry{entry1},
		Index: models.SegmentIndex{
			ByLevel: map[string][]int{"INFO": {0}},
		},
	}
	store := models.LogStore{Segments: []models.Segment{segment}}

	startTime := time.Time{}
	endTime := now.Add(1 * time.Hour)

	result := FilterLogs(store, []string{"INFO"}, nil, nil, nil, startTime, endTime)
	if len(result) != 0 {
		t.Errorf("expected 0 logs (skipped due to after end), got %d", len(result))
	}
}
