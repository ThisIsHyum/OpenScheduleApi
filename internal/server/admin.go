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
	admin        = "/admin"
	createParser = "/parser"
	deleteParser = "/parser/:parserId"
)

type AdminHandler struct {
	adminService *service.AdminService
	logger       *logrus.Logger
}

func NewAdminHandler(app *fiber.App,
	repos repository.Repos, adminToken string, logger *logrus.Logger) {
	adminService := service.NewAdminService(repos.CollegeRepo, repos.CampusRepo, repos.CreateTx)
	handler := AdminHandler{adminService: adminService, logger: logger}
	mh := MiddlewareHandler{logger: logger, adminToken: adminToken}
	app.Group(admin, mh.AdminAuthMiddleware).
		Post(createParser, handler.NewParser).
		Delete(deleteParser, handler.DeleteParser)
}
func (h AdminHandler) NewParser(ctx fiber.Ctx) error {
	requestBody := dto.NewParserRequest{}

	if err := ctx.Bind().Body(&requestBody); err != nil {
		return dto.NewErrorResponse("invalid request body", fiber.StatusBadRequest).Send(ctx)
	}

	token, err := h.adminService.NewParser(ctx, requestBody.CollegeName, requestBody.CampusNames)
	if errors.Is(err, domain.ErrConflict) {
		return dto.NewErrorResponse("college already exists", fiber.StatusConflict).Send(ctx)
	} else if err != nil {
		h.logger.WithField("error", err).Error("unable to get colleges")
		return dto.NewErrorResponse("internal server error", fiber.StatusInternalServerError).Send(ctx)
	}

	return ctx.Status(fiber.StatusCreated).JSON(dto.NewParserResponse{Token: token})
}

func (h AdminHandler) DeleteParser(ctx fiber.Ctx) error {
	c := ctx.Context()
	id := fiber.Params[uint](ctx, "parserId")
	if id == 0 {
		return dto.NewErrorResponse("invalid parserId", fiber.StatusBadRequest).Send(ctx)
	}

	if err := h.adminService.DeleteParser(c, id); errors.Is(err, domain.ErrNotFound) {
		return dto.NewErrorResponse("parser not found", fiber.StatusNotFound).Send(ctx)
	} else if err != nil {
		h.logger.WithError(err).Error("unable to delete parser")
		return dto.NewErrorResponse("internal server error", fiber.StatusInternalServerError).Send(ctx)
	}
	return ctx.SendStatus(fiber.StatusNoContent)
}
