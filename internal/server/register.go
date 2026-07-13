package server

import (
	"github.com/ThisIsHyum/OpenScheduleApi/internal/repository"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	lg "github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/sirupsen/logrus"
)

func Register(app *fiber.App, repos repository.Repos, logger *logrus.Logger, adminToken string) {
	app.Use(cors.New())
	app.Use(lg.New())

	NewCollegeHandler(app, repos.CollegeRepo, repos.CampusRepo, logger)
	NewCampusHandler(app, repos.CampusRepo, repos.StudentGroupRepo, repos.CollegeRepo, logger)
	NewGroupHandler(app, repos.StudentGroupRepo, repos.CampusRepo, repos.CollegeRepo, logger)
	NewScheduleHandler(app, repos.StudentGroupRepo, repos.LessonRepo, repos.CallRepo, repos.CollegeRepo, logger)
	NewParserHandler(app, repos.CallRepo, repos.StudentGroupRepo, repos.LessonRepo, repos.CampusRepo, repos.CollegeRepo, logger)
	NewAdminHandler(app, repos.CollegeRepo, repos.CampusRepo, repos.CreateTx, adminToken, logger)
}
