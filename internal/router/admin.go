package router

import (
	"github.com/ThisIsHyum/OpenScheduleApi/internal/database"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/dto"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/token"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

const (
	admin        = "/admin"
	createParser = "/parser"
	deleteParser = "/parser/:parserId"
)

type AdminHandler struct {
	db     *database.Db
	logger *logrus.Logger
}

func NewAdminHandler(app *fiber.App, db *database.Db, logger *logrus.Logger, adminToken string) {
	handler := AdminHandler{db: db, logger: logger}
	mh := MiddlewareHandler{db: db, logger: logger, adminToken: adminToken}
	app.Group(admin, mh.AdminAuthMiddleware).
		Post(createParser, handler.NewParser).
		Delete(deleteParser, handler.DeleteParser)
}
func (h AdminHandler) NewParser(ctx fiber.Ctx) error {
	requestBody := dto.NewParserRequest{}

	if err := ctx.Bind().Body(&requestBody); err != nil {
		return dto.NewErrorResponse("invalid request body", fiber.StatusBadRequest).Send(ctx)
	}

	if colleges, err := h.db.GetCollegesByName(requestBody.CollegeName); err != nil {
		h.logger.WithField("error", err).Error("unable to get colleges")
		return dto.NewErrorResponse("internal server error", fiber.StatusInternalServerError).Send(ctx)
	} else if len(colleges) != 0 {
		return dto.NewErrorResponse("college already exists", fiber.StatusConflict).Send(ctx)
	}

	token, err := token.GenerateToken()
	if err != nil {
		h.logger.WithField("error", err).Error("unable to create token")
		return dto.NewErrorResponse("internal server error", fiber.StatusInternalServerError).Send(ctx)
	}

	if err = h.db.NewParser(token, requestBody.CollegeName, requestBody.CampusNames); err != nil {
		h.logger.WithError(err).Error("unable to create parser")
		return dto.NewErrorResponse("internal server error", fiber.StatusInternalServerError).Send(ctx)
	}
	return ctx.Status(fiber.StatusCreated).JSON(dto.NewParserResponse{Token: token})
}

func (h AdminHandler) DeleteParser(ctx fiber.Ctx) error {
	id := fiber.Params[uint](ctx, "parserId")
	if id == 0 {
		return dto.NewErrorResponse("invalid parserId", fiber.StatusBadRequest).Send(ctx)
	}

	deleted, err := h.db.DeleteParser(id)
	if !deleted {
		return dto.NewErrorResponse("parser not found", fiber.StatusNotFound).Send(ctx)
	} else if err != nil {
		h.logger.WithError(err).Error("unable to delete parser")
		return dto.NewErrorResponse("internal server error", fiber.StatusInternalServerError).Send(ctx)
	}
	return ctx.SendStatus(fiber.StatusNoContent)
}
