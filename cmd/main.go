package main

import (
	"fmt"
	"miniproject/segment"
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
	LogStore, _ := segment.ParseLogSegments("../logs")
	// for _, segment := range LogStore.Segments {
	// 	fmt.Println(segment)
	// }
	fmt.Println(LogStore.Segments[0])
}
