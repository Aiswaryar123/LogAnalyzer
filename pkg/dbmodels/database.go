package database

import (
	"time"

	"gorm.io/gorm"
)

type Entry struct {
	gorm.Model
	TimeStamp   time.Time
	LevelID     uint
	ComponentID uint
	HostID      uint
	RequestId   string
	Message     string

	Level     LogLevel     `gorm:"foreignKey:LevelID;references:ID"`
	Component LogComponent `gorm:"foreignKey:ComponentID;references:ID"`
	Host      LogHost      `gorm:"foreignKey:HostID;references:ID"`
}
type queryComponent struct {
	key      string
	value    []string
	operator string
}
type LogLevel struct {
	ID    uint   `gorm:"primaryKey"`
	Level string `gorm:"unique;not null"`
}

type LogComponent struct {
	ID        uint   `gorm:"primaryKey"`
	Component string `gorm:"unique;not null"`
}

type LogHost struct {
	ID   uint   `gorm:"primaryKey"`
	Host string `gorm:"unique;not null"`
}
