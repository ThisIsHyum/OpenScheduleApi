package router

import (
	"github.com/ThisIsHyum/OpenScheduleApi/internal/database"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	lg "github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/sirupsen/logrus"
)

func Register(app *fiber.App, db *database.Db, logger *logrus.Logger, adminToken string) {
	app.Use(cors.New())
	app.Use(lg.New())

	NewCollegeHandler(app, db, logger)
	NewCampusHandler(app, db, logger)
	NewGroupHandler(app, db, logger)
	NewScheduleHandler(app, db, logger)
	NewParserHandler(app, db, logger)
	NewAdminHandler(app, db, logger, adminToken)
}
