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
	getColleges = "/colleges"
	getCollege  = "/colleges/:collegeId"
)

type CollegeHandler struct {
	service *service.CollegeService
	logger  *logrus.Logger
}

func NewCollegeHandler(app *fiber.App, repos repository.Repos, logger *logrus.Logger) {
	collegeService := service.NewCollegeService(repos.CollegeRepo, repos.CampusRepo)
	handler := CollegeHandler{service: collegeService, logger: logger}
	app.Get(getColleges, handler.GetColleges)
	app.Get(getCollege, handler.GetCollege)
}

func (h CollegeHandler) GetColleges(ctx fiber.Ctx) error {
	c := ctx.Context()
	name := ctx.Query("name")
	colleges, err := h.service.GetColleges(c, name)
	if err != nil {
		h.logger.WithField("error", err).Error("unable to get colleges")
		return dto.NewErrorResponse("internal server error", fiber.StatusInternalServerError).Send(ctx)
	}
	return ctx.JSON(colleges)
}

func (h CollegeHandler) GetCollege(ctx fiber.Ctx) error {
	c := ctx.Context()
	id := fiber.Params[uint](ctx, "collegeId")
	if id == 0 {
		return dto.NewErrorResponse("invalid collegeId", fiber.StatusBadRequest).Send(ctx)
	}

	college, err := h.service.GetCollege(c, id)
	if errors.Is(err, domain.ErrNotFound) {
		return dto.NewErrorResponse("college not found", fiber.StatusNotFound).Send(ctx)
	} else if err != nil {
		h.logger.WithField("error", err).Error("unable to get college")
		return dto.NewErrorResponse("internal server error", fiber.StatusInternalServerError).Send(ctx)
	}
	return ctx.JSON(college)
}
