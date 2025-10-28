package parser

import (
	"fmt"
	models "miniproject/model"
	"path/filepath"

	"bufio"
	"os"
	"regexp"
	"time"
)

func LogParseEntry(s string) (*models.LogEntry, error) {
	pattern := `^(?P<time>\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d+)\s+\|\s+(?P<level>[A-Z]+)\s+\|\s+(?P<component>[\w-]+)\s+\|\s+host=(?P<host>[\w-]+)\s+\|\s+request_id=(?P<request_id>[\w-]+)\s+\|\s+msg="(?P<msg>.*)"$`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(s)
	if matches == nil {
		return nil, fmt.Errorf("Invalid format")
	}
	result := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if name != "" {
			result[name] = matches[i]
		}
	}
	t, err := time.Parse("2006-01-02 15:04:05.000", result["time"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse time: %v", err)
	}
	entry := models.LogEntry{
		Raw:       matches[0],
		Time:      t,
		Level:     models.Loglevel(result["level"]),
		Component: result["component"],
		Host:      result["host"],
		Requestid: result["request_id"],
		Message:   result["msg"],
	}
	return &entry, nil
}
func ParseLogFiles(directory string) ([]models.LogEntry, error) {
	var allEntries []models.LogEntry
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory : %v", err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filepath := filepath.Join(directory, file.Name())
		f, err := os.Open(filepath)
		if err != nil {
			fmt.Printf("Skipping file %s due to error: %v\n", filepath, err)
			continue
		}
		scanner := bufio.NewScanner(f)
		scanner.Buffer(make([]byte, 0, 1024*1024), 10*1024*1024)

		for scanner.Scan() {
			line := scanner.Text()
			entry, err := LogParseEntry(line)
			if err == nil {
				allEntries = append(allEntries, *entry)
			}
		}
		err = scanner.Err()
		if err != nil {
			fmt.Printf("Error reading file %s due to error :%v.", filepath, err)
		}
		f.Close()

	}
	return allEntries, nil
}
