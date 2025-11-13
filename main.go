package main

import (
	"fmt"

	"github.com/ThisIsHyum/OpenScheduleApi/internal/config"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/database"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/logger"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/router"
	"github.com/gofiber/fiber/v3"
)

func main() {
	logger := logger.New()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.WithError(err).Fatal("unable to load config")
	}

	db, err := database.NewDb(cfg)
	if err != nil {
		logger.WithError(err).Fatal("unable to connect database")
	}

	app := fiber.New()
	router.Register(app, db, logger, cfg.AdminToken)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Infof("Running server on %s \n", addr)
	if err := app.Listen(addr, fiber.ListenConfig{
		DisableStartupMessage: true,
	}); err != nil {
		logger.WithError(err).Fatal("unable to run server")
	}
}
