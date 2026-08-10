package server

import (
	"github.com/ThisIsHyum/OpenScheduleApi/docs"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/repository"
	"github.com/gofiber/contrib/v3/swaggerui"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	lg "github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/sirupsen/logrus"
)

func Register(app *fiber.App, repos repository.Repos, logger *logrus.Logger, adminToken string) {
	app.Use(cors.New())
	app.Use(lg.New())

	app.Use(swaggerui.New(swaggerui.Config{
		FileContent: docs.SwaggerJSON,
		Path:        "swagger",
		Title:       "Open Schedule API",
	}))

	NewCollegeHandler(app, repos, logger)
	NewCampusHandler(app, repos, logger)
	NewGroupHandler(app, repos, logger)
	NewScheduleHandler(app, repos, logger)
	NewParserHandler(app, repos, logger)
	NewAdminHandler(app, repos, adminToken, logger)
}
