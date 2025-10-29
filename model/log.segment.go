package models

import "time"

type Segment struct {
	Filename   string
	LogEntries []LogEntry
	Starttime  time.Time
	Endtime    time.Time
	Index      SegmentIndex
}

type SegmentIndex struct {
	ByLevel     map[string][]int
	ByComponent map[string][]int
	ByHost      map[string][]int
	ByRequestID map[string][]int
}
