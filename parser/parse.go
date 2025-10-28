package parser

import (
	"fmt"
	models "miniproject/model"

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
