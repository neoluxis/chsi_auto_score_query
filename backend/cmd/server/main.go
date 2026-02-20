package main

import (
"log"

"chsi-auto-score-query/internal/api"
"chsi-auto-score-query/internal/db"
"chsi-auto-score-query/internal/logger"
"chsi-auto-score-query/pkg/config"
)

func main() {
	// 初始化配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	logger.Init(cfg.LogLevel)
	logger.Info("Application starting")

	// 初始化数据库
	database, err := db.Init(cfg)
	if err != nil {
		logger.Error("Failed to init database: %v", err)
		log.Fatalf("Database init failed: %v", err)
	}
	logger.Info("Database initialized")

	// 启动API服务
	server := api.NewServer(cfg, database)
	if err := server.Start(); err != nil {
		logger.Error("Failed to start server: %v", err)
		log.Fatalf("Server start failed: %v", err)
	}
}
