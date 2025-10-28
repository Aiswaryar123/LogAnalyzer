package models

import (
	"time"
)

type Loglevel string

const (
	INFO  Loglevel = "INFO"
	WARN  Loglevel = "WARN"
	DEBUG Loglevel = "DEBUG"
	ERROR Loglevel = "ERROR"
)

type LogEntry struct {
	Raw       string
	Time      time.Time
	Level     Loglevel
	Component string
	Host      string
	Requestid string
	Message   string
}
