package router

import (
	"errors"

	"github.com/ThisIsHyum/OpenScheduleApi/internal/database"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/dto"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/models"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	getColleges = "/colleges"
	getCollege  = "/colleges/:collegeId"
)

type CollegeHandler struct {
	db     *database.Db
	logger *logrus.Logger
}

func NewCollegeHandler(app *fiber.App, db *database.Db, logger *logrus.Logger) {
	handler := CollegeHandler{db: db, logger: logger}
	app.Get(getColleges, handler.GetColleges)
	app.Get(getCollege, handler.GetCollege)
}

func (h CollegeHandler) GetColleges(ctx fiber.Ctx) error {
	name := ctx.Query("name")

	var colleges []models.College
	var err error
	if name != "" {
		colleges, err = h.db.GetCollegesByName(name)
	} else {
		colleges, err = h.db.GetColleges()
	}
	if err != nil {
		h.logger.WithField("error", err).Error("unable to get colleges")
		return dto.NewErrorResponse("internal server error", fiber.StatusInternalServerError).Send(ctx)
	}
	return ctx.JSON(colleges)
}

func (h CollegeHandler) GetCollege(ctx fiber.Ctx) error {
	id := fiber.Params[uint](ctx, "collegeId")
	if id == 0 {
		return dto.NewErrorResponse("invalid collegeId", fiber.StatusBadRequest).Send(ctx)
	}

	college, err := h.db.GetCollege(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.NewErrorResponse("college not found", fiber.StatusNotFound).Send(ctx)
	} else if err != nil {
		h.logger.WithField("error", err).Error("unable to get college")
		return dto.NewErrorResponse("internal server error", fiber.StatusInternalServerError).Send(ctx)
	}
	return ctx.JSON(college)
}
