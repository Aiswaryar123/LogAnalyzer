package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func (e Entry) String() string {
	return fmt.Sprintf("%s : %s : %s : %s : %s : %s",
		e.TimeStamp.Format("2006-01-02 15:04:05"),
		e.Level.Level,
		e.Component.Component,
		e.Host.Host,
		e.RequestId,
		e.Message,
	)
}

func CreateDB(dbUrl string) (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{Logger: newLogger})
	if err != nil {
		return nil, fmt.Errorf("couldn't open database %v", err)
	}
	return db, nil

}
func InitDB(db *gorm.DB) error {
	db.AutoMigrate(&LogLevel{}, &LogComponent{}, &LogHost{}, &Entry{})
	db.Create(&[]LogLevel{
		{Level: "INFO"},
		{Level: "WARN"},
		{Level: "ERROR"},
		{Level: "DEBUG"},
	})
	db.Create(&[]LogComponent{
		{Component: "api-server"},
		{Component: "database"},
		{Component: "cache"},
		{Component: "worker"},
		{Component: "auth"},
	})
	db.Create(&[]LogHost{
		{Host: "web01"},
		{Host: "web02"},
		{Host: "cache01"},
		{Host: "worker01"},
		{Host: "db01"},
	})
	return nil
}

func AddDB(db *gorm.DB, e Entry) error {
	ctx := context.Background()
	err := gorm.G[Entry](db).Create(ctx, &e)
	if err != nil {
		return err
	}
	return nil
}
func parseQuery(parts []string) ([]queryComponent, error) {
	var ret []queryComponent

	pattern := `^(?P<key>[^\s=!<>]+)\s*(?P<operator>=|!=|>=|<=|>|<)\s*(?P<value>.+)$`
	r := regexp.MustCompile(pattern)

	for _, part := range parts {
		part = strings.TrimSpace(part)

		matches := r.FindStringSubmatch(part)
		if matches == nil {
			return nil, fmt.Errorf("invalid condition: %s", part)
		}

		rawValue := matches[r.SubexpIndex("value")]
		rawValue = strings.ReplaceAll(rawValue, "|", ",")

		vals := strings.Split(rawValue, ",")

		ret = append(ret, queryComponent{
			key:      matches[r.SubexpIndex("key")],
			operator: matches[r.SubexpIndex("operator")],
			value:    vals,
		})
	}

	return ret, nil
}

func QueryDB(db *gorm.DB, query []string) ([]Entry, error) {
	var ret []Entry

	parsed, err := parseQuery(query)
	if err != nil {
		return nil, err
	}

	fmt.Println("Parsed conditions:", parsed)

	q := db

	for _, c := range parsed {

		key := strings.ToLower(c.key)

		switch key {

		case "level":

			var ids []uint
			for _, v := range c.value {
				var lvl LogLevel
				if err := db.First(&lvl, "level = ?", v).Error; err != nil {
					return nil, fmt.Errorf("unknown level '%s'", v)
				}
				ids = append(ids, lvl.ID)
			}
			c.key = "level_id"
			c.value = toStringSlice(ids)

		case "component":
			var ids []uint
			for _, v := range c.value {
				var comp LogComponent
				if err := db.First(&comp, "component = ?", v).Error; err != nil {
					return nil, fmt.Errorf("unknown component '%s'", v)
				}
				ids = append(ids, comp.ID)
			}
			c.key = "component_id"
			c.value = toStringSlice(ids)

		case "host":
			var ids []uint
			for _, v := range c.value {
				var h LogHost
				if err := db.First(&h, "host = ?", v).Error; err != nil {
					return nil, fmt.Errorf("unknown host '%s'", v)
				}
				ids = append(ids, h.ID)
			}
			c.key = "host_id"
			c.value = toStringSlice(ids)
		}

		if len(c.value) == 1 {
			q = q.Where(fmt.Sprintf("%s %s ?", c.key, c.operator), c.value[0])
		} else {
			if c.operator == "!=" {
				q = q.Where(fmt.Sprintf("%s NOT IN ?", c.key), c.value)
			} else {
				q = q.Where(fmt.Sprintf("%s IN ?", c.key), c.value)
			}
		}
	}
	q = q.
		Preload("Level").
		Preload("Component").
		Preload("Host")
	if err := q.Find(&ret).Error; err != nil {
		return nil, err
	}

	return ret, nil
}
func SplitUserFilter(input string) []string {
	var parts []string
	current := ""
	tokens := strings.Fields(input)

	for _, tok := range tokens {

		if strings.Contains(tok, "=") ||
			strings.Contains(tok, ">=") ||
			strings.Contains(tok, "<=") ||
			strings.Contains(tok, ">") ||
			strings.Contains(tok, "<") {

			if current != "" {
				parts = append(parts, current)
			}
			current = tok
		} else {

			current += " " + tok
		}
	}

	if current != "" {
		parts = append(parts, current)
	}

	return parts
}

func toStringSlice(nums []uint) []string {
	s := make([]string, len(nums))
	for i, n := range nums {
		s[i] = fmt.Sprint(n)
	}
	return s
}
func GetAllLogs(db *gorm.DB) ([]Entry, error) {
	var result []Entry
	err := db.Preload("Level").
		Preload("Component").
		Preload("Host").
		Find(&result).Error
	return result, err
}
func FilterLogs(
	DB *gorm.DB,
	levels []string,
	components []string,
	hosts []string,
	requestId string,
	startTime string,
	endTime string,
) ([]Entry, error) {

	query := DB.Model(&Entry{})

	if len(levels) > 0 {
		var levelIDs []uint
		DB.Model(&LogLevel{}).Where("level IN ?", levels).Pluck("id", &levelIDs)
		if len(levelIDs) > 0 {
			query = query.Where("level_id IN ?", levelIDs)
		}
	}

	if len(components) > 0 {
		var compIDs []uint
		DB.Model(&LogComponent{}).Where("component IN ?", components).Pluck("id", &compIDs)
		if len(compIDs) > 0 {
			query = query.Where("component_id IN ?", compIDs)
		}
	}
	if len(hosts) > 0 {
		var hostIDs []uint
		DB.Model(&LogHost{}).Where("host IN ?", hosts).Pluck("id", &hostIDs)
		if len(hostIDs) > 0 {
			query = query.Where("host_id IN ?", hostIDs)
		}
	}

	if requestId != "" {
		query = query.Where("request_id = ?", requestId)
	}

	if startTime != "" && endTime != "" {
		query = query.Where("time_stamp BETWEEN ? AND ?", startTime, endTime)
	} else if startTime != "" {
		query = query.Where("time_stamp >= ?", startTime)
	} else if endTime != "" {
		query = query.Where("time_stamp <= ?", endTime)
	}
	var entries []Entry
	result := query.
		Preload("Level").
		Preload("Component").
		Preload("Host").
		Order("time_stamp DESC").
		Find(&entries)

	return entries, result.Error
}
