package database

import (
	"fmt"
	"log/slog"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"siato_lms/internal/config"
)

func Open(c config.Config, log *slog.Logger) (*gorm.DB, error) {
	level := logger.Warn
	if c.Environment == "development" {
		level = logger.Info
	}
	db, err := gorm.Open(mysql.Open(c.DB.DSN()), &gorm.Config{Logger: logger.Default.LogMode(level), NowFunc: func() time.Time { return time.Now().UTC() }})
	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(c.DB.MaxOpen)
	sqlDB.SetMaxIdleConns(c.DB.MaxIdle)
	sqlDB.SetConnMaxLifetime(c.DB.ConnMaxLifetime)
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}
	log.Info("database connected", "host", c.DB.Host, "database", c.DB.Name)
	return db, nil
}
