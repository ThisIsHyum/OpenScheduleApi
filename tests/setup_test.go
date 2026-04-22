package tests

import (
	"testing"

	"github.com/ThisIsHyum/OpenScheduleApi/internal/database"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/server"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const TestAdminToken = "K3r98j28f4j3"

func SetupApp(db *gorm.DB, t *testing.T) *fiber.App {
	t.Helper()
	collegeRepo := database.NewCollegeDb(db)
	campusRepo := database.NewCampusDb(db)
	studentGroupRepo := database.NewGroupDb(db)
	callRepo := database.NewCallDb(db)
	lessonRepo := database.NewLessonDb(db)
	createTx := database.InitCreateTx(db)

	app := fiber.New()
	server.Register(app,
		collegeRepo, campusRepo,
		studentGroupRepo, callRepo, lessonRepo, createTx, logrus.New(), TestAdminToken)
	return app
}
