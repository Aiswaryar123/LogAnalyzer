package segment

import (
	"bufio"
	"fmt"
	"sync"
	"time"

	models "Log_analyzer/model"
	indexer "Log_analyzer/pkg/Indexer"
	"Log_analyzer/pkg/parser"
	"os"
	"path/filepath"
)

func ParseLogSegments(path string) (models.LogStore, error) {
	start := time.Now()
	LogStore := models.LogStore{
		Segments: []models.Segment{},
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return LogStore, fmt.Errorf("failed to read directory : %v", err)
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		wg.Add(1)
		go func(file os.DirEntry) {
			defer wg.Done()

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
				return
			}
			index := indexer.BuildSegmentIndex(LogEntries)
			segment := models.Segment{
				Filename:   file.Name(),
				LogEntries: LogEntries,
				Starttime:  LogEntries[0].Time,
				Endtime:    LogEntries[len(LogEntries)-1].Time,
				Index:      index,
			}
			mu.Lock()
			LogStore.Segments = append(LogStore.Segments, segment)
			mu.Unlock()
		}(file)

	}
	wg.Wait()
	end := time.Since(start)
	fmt.Println(end)
	return LogStore, nil

}
