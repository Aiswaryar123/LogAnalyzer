package segment

import (
	"bufio"
	"fmt"

	indexer "miniproject/Indexer"
	models "miniproject/model"
	"miniproject/parser"
	"os"
	"path/filepath"
)

func ParseLogSegments(path string) (models.LogStore, error) {
	LogStore := models.LogStore{
		Segments: []models.Segment{},
	}
	files, err := os.ReadDir(path)
	if err != nil {
		return LogStore, fmt.Errorf("failed to read directory : %v", err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filepath := filepath.Join(path, file.Name())
		f, err := os.Open(filepath)
		if err != nil {
			fmt.Printf("Skipping file %s due to error: %v", filepath, err)
		}
		var LogEntries []models.LogEntry
		scanner := bufio.NewScanner(f)
		scanner.Buffer(make([]byte, 0, 1024*1024), 10*1024*1024)

		for scanner.Scan() {
			line := scanner.Text()
			entry, err := parser.LogParseEntry(line)
			if err == nil {
				LogEntries = append(LogEntries, *entry)
			}
		}
		f.Close()
		if len(LogEntries) == 0 {
			continue
		}
		index := indexer.BuildSegmentIndex(LogEntries)
		segment := models.Segment{
			Filename:   file.Name(),
			LogEntries: LogEntries,
			Starttime:  LogEntries[0].Time,
			Endtime:    LogEntries[len(LogEntries)-1].Time,
			Index:      index,
		}
		LogStore.Segments = append(LogStore.Segments, segment)
	}
	return LogStore, nil

}
