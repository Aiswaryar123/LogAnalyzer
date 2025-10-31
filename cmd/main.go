package main

import (
	"flag"
	"fmt"
	"log/slog"
	"miniproject/filter"
	"miniproject/segment"
	"strings"
	"time"
)

// func main() {
// 	logLine := `2025-10-23 15:04:10.001 | DEBUG | auth | host=db01 | request_id=req-hyx6sa-8587 | msg="2FA verification completed"`
// 	entry, err := parser.LogParseEntry(logLine)

// 	if err != nil {
// 		fmt.Println("Error:", err)
// 	}

//		fmt.Println("Time:", entry.Time.Format("2006-01-02 15:04:05.000"))
//		fmt.Println("Level:", entry.Level)
//		fmt.Println("Component:", entry.Component)
//		fmt.Println("Host:", entry.Host)
//		fmt.Println("Request ID:", entry.Requestid)
//		fmt.Println("Message:", entry.Message)
//	}
func main() {
	// LogStore, _ := segment.ParseLogSegments("../logs")
	// // for _, segment := range LogStore.Segments {
	// // 	fmt.Println(segment)
	// // }
	// fmt.Println(LogStore.Segments[0])
	level := flag.String("level", "", "Filter by log level")
	component := flag.String("component", "", "Filter by component")
	host := flag.String("host", "", "Filter by host")
	reqID := flag.String("reqID", "", "Filter by requestID")
	startTime := flag.String("Starttime", "", "filter by start time")
	EndTime := flag.String("Endtime", "", "filter by End time")
	flag.Parse()

	logStore, err := segment.ParseLogSegments("../logs")
	// fmt.Println(logStore.Segments[0])
	if err != nil {
		slog.Error("Failed to parse logs\n")
	}
	split := func(s string) []string {
		if s == "" {
			return nil
		}
		parts := strings.Split(s, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		return parts
	}
	levels := split(*level)
	components := split(*component)
	hosts := split(*host)
	reqIDs := split(*reqID)
	// starttime, err := time.Parse("2006-01-02 15:04:05", *startTime)
	var starttime, endtime time.Time
	if *startTime != "" {
		starttime, err = time.Parse("2006-01-02 15:04:05", *startTime)

		if err != nil {
			slog.Error("Error parsing time: ", "error", err)
		}
	}
	// endTime, err := time.Parse("2006-01-02 15:04:05", *EndTime)
	if *EndTime != "" {
		endtime, err = time.Parse("2006-01-02 15:04:05", *EndTime)

		if err != nil {
			slog.Error("Error parsing time: ", "error", err)
		}
	}

	filteredLogs := filter.FilterLogs(logStore, levels, components, hosts, reqIDs, starttime, endtime)
	fmt.Printf("Found %d matching entries\n", len(filteredLogs))
	for _, entry := range filteredLogs {
		fmt.Println(entry.Raw)
	}
}
