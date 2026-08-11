package main

import (
	"fmt"

	"github.com/ThisIsHyum/OpenScheduleApi/internal/config"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/database"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/logger"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/server"
	"github.com/gofiber/fiber/v3"
)

// @title Open Schedule API
// @version 0.3.1
// @description API for managing college schedules
// @license.name GPL-3.0
// @license.url https://www.gnu.org/licenses/gpl-3.0.html
// @securityDefinitions.apikey AdminBearerAuth
// @in header
// @name Authorization
// @securityDefinitions.apikey ParserBearerAuth
// @in header
// @name Authorization

func main() {
	logger := logger.New()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.WithError(err).Fatal("unable to load config")
	}

	db, err := database.NewDb(&cfg.Db)
	if err != nil {
		logger.WithError(err).Fatal("unable to connect database")
	}

	repos := database.NewRepos(db)

	app := fiber.New()
	server.Register(app, repos, logger, cfg.AdminToken)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Infof("Running server on %s", addr)
	if err := app.Listen(addr, fiber.ListenConfig{
		DisableStartupMessage: true,
	}); err != nil {
		logger.WithError(err).Fatal("unable to run server")
	}
}
