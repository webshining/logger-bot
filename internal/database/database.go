package database

import (
	"bot/internal/settings"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect(settings settings.DatabaseConfig) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	switch settings.Driver {
	case "postgres":
		db, err = gorm.Open(postgres.Open(settings.Url))
	default:
		db, err = gorm.Open(sqlite.Open(settings.Url))
	}

	return db, err
}

func MustConnect(config settings.DatabaseConfig, logger *zap.Logger) *gorm.DB {
	db, err := Connect(config)
	if err != nil {
		logger.Named("database").Fatal("failed to connect to database", zap.String("driver", config.Driver), zap.String("url", config.Url), zap.Error(err))
	}

	return db
}
