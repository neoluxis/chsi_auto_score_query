package db

import (
"chsi-auto-score-query/internal/logger"
"chsi-auto-score-query/internal/model"
"chsi-auto-score-query/pkg/config"
"gorm.io/driver/sqlite"
"gorm.io/gorm"
)

var DB *gorm.DB

func Init(cfg *config.Config) (*gorm.DB, error) {
	database, err := gorm.Open(sqlite.Open(cfg.DatabaseDSN), &gorm.Config{})
	if err != nil {
		logger.Error("Failed to connect database: %v", err)
		return nil, err
	}

	// 自动迁移
	err = database.AutoMigrate(&model.User{})
	if err != nil {
		logger.Error("Failed to auto migrate: %v", err)
		return nil, err
	}

	// 清空数据库（如果需要）
	if cfg.ClearDBOnStart {
		database.Exec("DELETE FROM users")
		logger.Info("Database cleared on start")
	}

	DB = database
	return database, nil
}
