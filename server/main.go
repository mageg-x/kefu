package main

import (
	"flag"
	"log"

	"kefu-server/config"
	"kefu-server/models"
	"kefu-server/router"
	"kefu-server/store"
	"kefu-server/utils/logger"
)

func main() {
	logger.SetLevel(logger.INFO)

	// 解析命令行参数
	configPath := flag.String("config", "./config/config.yaml", "Path to config file")
	flag.Parse()

	// 打印配置文件路径
	logger.Infof("Using config file: %s", *configPath)

	// 加载配置
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		logger.Errorf("failed to load config: %v", err)
		log.Fatal(err)
	}

	// 初始化数据库
	db, err := store.InitDB(cfg.Admin.Database)
	if err != nil {
		logger.Errorf("database connection failed: %v", err)
		log.Fatal(err)
	}

	// 数据库迁移
	if err := db.AutoMigrate(&models.User{}, &models.App{}); err != nil {
		logger.Errorf("database migration failed: %v", err)
		log.Fatal(err)
	}

	// 创建默认用户
	if err := models.CreateDefaultUsers(db); err != nil {
		logger.Errorf("failed to create default users: %v", err)
		log.Fatal(err)
	}

	// 设置路由
	r := router.SetupRouter()

	logger.Infof("server starting on %s...", cfg.Admin.Address)
	if err := r.Run(cfg.Admin.Address); err != nil {
		logger.Errorf("server startup failed: %v", err)
		log.Fatal(err)
	}
}
