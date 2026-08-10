package server

import (
	"errors"

	"github.com/ThisIsHyum/OpenScheduleApi/internal/domain"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/dto"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/repository"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/service"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

const (
	getCampuses = "/colleges/:collegeId/campuses"
	getCampus   = "/campuses/:id"
)

type CampusHandler struct {
	service *service.CampusService
	logger  *logrus.Logger
}

func NewCampusHandler(app *fiber.App, repos repository.Repos, logger *logrus.Logger) {
	campusService := service.NewCampusService(
		repos.CampusRepo, repos.StudentGroupRepo, repos.CollegeRepo)
	handler := CampusHandler{service: campusService, logger: logger}
	app.Get(getCampuses, handler.GetCampuses)
	app.Get(getCampus, handler.GetCampus)
}

// @Summary get campuses by college ID
// @Description get all campuses by college ID or a campus by college ID with specified name
// @Tags campuses
// @Produce json
// @Param collegeId path  int    true  "College ID"
// @Param name      query string false  "Campus name"
// @Success 200 {array}  dto.CampusResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /colleges/{collegeId}/campuses [get]
func (h CampusHandler) GetCampuses(ctx fiber.Ctx) error {
	c := ctx.Context()

	name := ctx.Query("name")
	id := fiber.Params[uint](ctx, "collegeId")
	if id == 0 {
		return dto.NewErrorResponse("invalid collegeId", fiber.StatusBadRequest).Send(ctx)
	}

	campuses, err := h.service.GetCampusesByCollegeID(c, id, name)
	if errors.Is(err, domain.ErrNotFound) {
		return dto.NewErrorResponse("college not found", fiber.StatusNotFound).Send(ctx)
	}
	if err != nil {
		h.logger.WithError(err).Error("unable to get campuses")
		return dto.NewErrorResponse("internal server error", fiber.StatusInternalServerError).Send(ctx)
	}
	return ctx.JSON(campuses)
}

// @Summary get a campus by ID
// @Description get a single campus by its ID
// @Tags campuses
// @Produce json
// @Param id path int true "Campus ID"
// @Success 200 {object}  dto.CampusResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /campuses/{id} [get]
func (h CampusHandler) GetCampus(ctx fiber.Ctx) error {
	c := ctx.Context()
	id := fiber.Params[uint](ctx, "id")
	if id == 0 {
		return dto.NewErrorResponse("invalid id", fiber.StatusBadRequest).Send(ctx)
	}
	campus, err := h.service.GetCampusByID(c, id)
	if errors.Is(err, domain.ErrNotFound) {
		return dto.NewErrorResponse("campus not found", fiber.StatusNotFound).Send(ctx)
	} else if err != nil {
		h.logger.WithError(err).Error("unable to get campus")
		return dto.NewErrorResponse("internal server error", fiber.StatusInternalServerError).Send(ctx)
	}
	return ctx.JSON(campus)
}
