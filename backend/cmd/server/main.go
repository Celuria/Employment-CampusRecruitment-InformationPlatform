package main

import (
	"fmt"
	"log"

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
	)
}
