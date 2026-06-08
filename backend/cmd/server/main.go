package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/employment-center/campus-recruitment/config"
	"github.com/employment-center/campus-recruitment/internal/handler"
	"github.com/employment-center/campus-recruitment/internal/model"
	"github.com/employment-center/campus-recruitment/internal/router"
	"github.com/employment-center/campus-recruitment/internal/service"
	"github.com/employment-center/campus-recruitment/pkg/database"
	"github.com/employment-center/campus-recruitment/pkg/jwt"
	"github.com/employment-center/campus-recruitment/pkg/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config failed: %v", err)
	}

	if err := logger.Init(cfg.Log.Level); err != nil {
		log.Fatalf("init logger failed: %v", err)
	}
	defer logger.Sync()

	db, err := database.Init(&cfg.Database)
	if err != nil {
		log.Fatalf("init database failed: %v", err)
	}

	if cfg.Database.AutoMigrate {
		if err := autoMigrate(db); err != nil {
			log.Fatalf("auto migrate failed: %v", err)
		}
	} else {
		logger.Log.Info("skip auto migrate, using SQL init script schema")
	}

	jwtManager := jwt.NewManager(cfg.JWT.Secret, cfg.JWT.ExpireHours, cfg.JWT.RememberExpireHours)
	services := service.NewServices(db, cfg, jwtManager)
	h := handler.NewHandler(services)

	// 启动提醒调度器：每分钟检查一次待发送提醒
	go reminderScheduler(services.Reminder)

	engine := router.Setup(cfg, h, jwtManager)
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Log.Info(fmt.Sprintf("server starting on %s", addr))
	if err := engine.Run(addr); err != nil {
		log.Fatalf("server run failed: %v", err)
	}
}

func autoMigrate(db interface{ AutoMigrate(...interface{}) error }) error {
	return db.AutoMigrate(
		&model.User{},
		&model.UserPreference{},
		&model.CareerTalk{},
		&model.JobFair{},
		&model.CalendarEvent{},
		&model.ReminderLog{},
		&model.AuditLog{},
		&model.SyncLog{},
	)
}

// reminderScheduler 提醒调度器，每分钟检查并处理到时间的待发送提醒
func reminderScheduler(reminderSvc service.ReminderService) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	sugar := logger.Log.Sugar()
	for range ticker.C {
		count, err := reminderSvc.ProcessPending(context.Background())
		if err != nil {
			sugar.Errorf("reminder scheduler error: %v", err)
			continue
		}
		if count > 0 {
			sugar.Infof("reminder scheduler: processed %d reminders", count)
		}
	}
}
