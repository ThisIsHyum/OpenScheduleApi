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
	getCampuses = "/colleges/:collegeId/campuses"
	getCampus   = "/campuses/:campusId"
)

type CampusHandler struct {
	db     *database.Db
	logger *logrus.Logger
}

func NewCampusHandler(app *fiber.App, db *database.Db, logger *logrus.Logger) {
	handler := CampusHandler{db: db, logger: logger}
	app.Get(getCampuses, handler.GetCampuses)
	app.Get(getCampus, handler.GetCampus)
}

func (h CampusHandler) GetCampuses(ctx fiber.Ctx) error {
	var name = ctx.Query("name")
	var id = fiber.Params[uint](ctx, "collegeId")
	if id == 0 {
		return dto.NewErrorResponse("invalid collegeId", fiber.StatusBadRequest).Send(ctx)
	}

	var campuses []models.Campus
	var err error

	if _, err := h.db.GetCollege(id); errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.NewErrorResponse("college not found", fiber.StatusNotFound).Send(ctx)
	} else if err != nil {
		h.logger.WithError(err).Error("unable to get college")
		return dto.NewErrorResponse("internal server error", fiber.StatusInternalServerError).Send(ctx)
	}
	if name != "" {
		campuses, err = h.db.GetCampusesByName(id, name)
	} else {
		campuses, err = h.db.GetCampusesByCollegeID(id)
	}
	if err != nil {
		h.logger.WithError(err).Error("unable to get campuses")
		return dto.NewErrorResponse("internal server error", fiber.StatusInternalServerError).Send(ctx)
	}
	return ctx.JSON(campuses)
}

func (h CampusHandler) GetCampus(ctx fiber.Ctx) error {
	id := fiber.Params[uint](ctx, "campusId")
	if id == 0 {
		return dto.NewErrorResponse("invalid campusId", fiber.StatusBadRequest).Send(ctx)
	}
	campus, err := h.db.GetCampusByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.NewErrorResponse("campus not found", fiber.StatusNotFound).Send(ctx)
	} else if err != nil {
		h.logger.WithError(err).Error("unable to get campus")
		return dto.NewErrorResponse("internal server error", fiber.StatusInternalServerError).Send(ctx)
	}
	return ctx.JSON(campus)
}
