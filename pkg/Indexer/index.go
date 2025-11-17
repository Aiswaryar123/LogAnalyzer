package indexer

import models "Log_analyzer/model"

func BuildSegmentIndex(LogEntries []models.LogEntry) models.SegmentIndex {
	Index := models.SegmentIndex{
		ByLevel:     make(map[string][]int),
		ByComponent: make(map[string][]int),
		ByHost:      make(map[string][]int),
		ByRequestID: make(map[string][]int),
	}
	for index, LogEntry := range LogEntries {
		Index.ByLevel[string(LogEntry.Level)] = append(Index.ByLevel[string(LogEntry.Level)], index)
		Index.ByComponent[LogEntry.Component] = append(Index.ByComponent[LogEntry.Component], index)
		Index.ByHost[LogEntry.Host] = append(Index.ByHost[LogEntry.Host], index)
		Index.ByRequestID[LogEntry.Requestid] = append(Index.ByRequestID[LogEntry.Requestid], index)
	}
	return Index
}
