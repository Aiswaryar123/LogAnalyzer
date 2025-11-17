package database

import (
	models "Log_analyzer/model"
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

func InsertLogs(ctx context.Context, conn *pgx.Conn, store models.LogStore) error {
	for _, segment := range store.Segments {
		for _, entry := range segment.LogEntries {
			_, err := conn.Exec(ctx, `
				INSERT INTO log_entry (time, level, host, component, requestid, msg)
				VALUES (
					$1,
					(SELECT id FROM log_level WHERE level = $2),
					(SELECT id FROM log_host WHERE host = $3),
					(SELECT id FROM log_component WHERE component = $4),
					$5,
					$6
				)`,
				entry.Time,
				entry.Level,
				entry.Host,
				entry.Component,
				entry.Requestid,
				entry.Message,
			)
			if err != nil {
				slog.Warn("Failed to insert log entry", "error", err, "entry", entry.Raw)
				continue
			}
		}
	}
	return nil
}
