package main

import (
	"Log_analyzer/database"
	"Log_analyzer/segment"
	"context"
	"flag"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	logPath := flag.String("path", "/home/user/miniproject/logs", "Path to the log directory")
	flag.Parse()
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("database_url"))
	if err != nil {
		slog.Error("Unable to connect to database", "error", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)
	logStore, err := segment.ParseLogSegments(*logPath)
	if err != nil {
		slog.Error("Failed to parse logs", "error", err)

	}

	err = database.InsertLogs(ctx, conn, logStore)
	if err != nil {
		slog.Error("Failed to insert logs", "error", err)
	} else {
		slog.Info("All logs inserted successfully")
	}
}
