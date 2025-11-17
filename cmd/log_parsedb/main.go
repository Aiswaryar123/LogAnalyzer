package main

import (
	database "Log_analyzer/pkg/dbmodels"
	"Log_analyzer/pkg/parser"
	"Log_analyzer/web"
	"fmt"
	"log"
	"log/slog"
	"os"
)

const dbUrl = "postgresql:///Log_analyzer_db?host=/var/run/postgresql/"

func commandHandler(args []string) error {
	db, err := database.CreateDB(dbUrl)
	if err != nil {
		return err
	}
	switch args[0] {
	case "init":
		err := database.InitDB(db)
		if err != nil {
			return err
		}
	case "add":
		dirPath := args[1]
		if dirPath == "" {
			slog.Error("Specify directory!")
		}
		entries, err := parser.ParseLogFiles(dirPath)
		if err != nil {
			return err
		}

		for _, p := range entries {

			var level database.LogLevel
			if err := db.First(&level, "level = ?", string(p.Level)).Error; err != nil {
				return fmt.Errorf("unknown level %s: %w", p.Level, err)
			}

			var component database.LogComponent
			if err := db.First(&component, "component = ?", p.Component).Error; err != nil {
				return fmt.Errorf("unknown component %s: %w", p.Component, err)
			}

			var host database.LogHost
			if err := db.First(&host, "host = ?", p.Host).Error; err != nil {
				return fmt.Errorf("unknown host %s: %w", p.Host, err)
			}

			dbEntry := database.Entry{
				TimeStamp:   p.Time,
				LevelID:     level.ID,
				ComponentID: component.ID,
				HostID:      host.ID,
				RequestId:   p.Requestid,
				Message:     p.Message,
			}

			if err := database.AddDB(db, dbEntry); err != nil {
				return err
			}
		}
		return nil

	case "query":
		queryList := args[1:]
		fmt.Println(queryList)

		entries, err := database.QueryDB(db, queryList)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			fmt.Println(entry)
		}
		slog.Info("Filtering successful!", "no. of entries:", len(entries))
		return nil
	case "web":

		db, err := database.CreateDB(dbUrl)
		if err != nil {
			return nil
		}

		r := web.SetupRouter(db)

		log.Println("Server running at http://localhost:8080")
		r.Run(":8080")
	default:
		slog.Warn("Unknown command!")
		return fmt.Errorf("unknown command %v", args[0])

	}
	return nil
}
func main() {
	err := commandHandler(os.Args[1:])
	if err != nil {
		slog.Error("Error in invocation", "error", err)
		os.Exit(-1)
	}
}
