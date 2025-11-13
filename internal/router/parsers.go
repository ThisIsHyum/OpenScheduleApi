package router

import (
	"github.com/ThisIsHyum/OpenScheduleApi/internal/database"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/dto"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/models"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

const (
	parser       = "/parser"
	getParser    = "/"
	updateCalls  = "/calls"
	addLessons   = "/lessons"
	updateGroups = "/groups"
)

type ParserHandler struct {
	db     *database.Db
	logger *logrus.Logger
}

func NewParserHandler(app *fiber.App, db *database.Db, logger *logrus.Logger) {
	handler := ParserHandler{db: db, logger: logger}
	mh := MiddlewareHandler{db: db, logger: logger}
	app.Group(parser, mh.ParserAuthMiddleware).
		Get(getParser, handler.GetParser).
		Post(updateCalls, handler.UpdateCalls).
		Post(addLessons, handler.AddLessons).
		Post(updateGroups, handler.UpdateGroups)
}

func (h ParserHandler) UpdateGroups(ctx fiber.Ctx) error {
	var requestBody dto.UpdateGroupsRequest
	if err := ctx.Bind().Body(&requestBody); err != nil {
		return dto.NewErrorResponse("invalid request body", fiber.StatusBadRequest).Send(ctx)
	}
	if requestBody.CampusID == 0 || len(requestBody.StudentGroupNames) == 0 {
		return dto.NewErrorResponse("invalid request data", fiber.StatusBadRequest).Send(ctx)
	}
	if err := h.db.UpdateGroups(requestBody.CampusID, requestBody.StudentGroupNames); err != nil {
		h.logger.WithError(err).Error("unable to update groups")
		return dto.NewErrorResponse("internal server error", fiber.StatusInternalServerError).Send(ctx)
	}
	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h ParserHandler) UpdateCalls(ctx fiber.Ctx) error {
	collegeId := fiber.Locals[uint](ctx, "collegeId")
	var calls []models.Call
	if err := ctx.Bind().Body(&calls); err != nil {
		return dto.NewErrorResponse("invalid request body", fiber.StatusBadRequest).Send(ctx)
	}
	if len(calls) == 0 {
		return dto.NewErrorResponse("empty calls list", fiber.StatusBadRequest).Send(ctx)
	}

	if err := h.db.UpdateCalls(collegeId, calls); err != nil {
		h.logger.WithError(err).Error("unable to update calls")
		return dto.NewErrorResponse("internal server error", fiber.StatusInternalServerError).Send(ctx)
	}
	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h ParserHandler) AddLessons(ctx fiber.Ctx) error {
	var lessons []models.Lesson
	if err := ctx.Bind().Body(&lessons); err != nil {
		return dto.NewErrorResponse("invalid request body", fiber.StatusBadRequest).Send(ctx)
	}

	if len(lessons) == 0 {
		return dto.NewErrorResponse("empty lessons list", fiber.StatusBadRequest).Send(ctx)
	}

	if err := h.db.AddLessons(lessons); err != nil {
		h.logger.WithError(err).Error("unable to add lessons")
		return dto.NewErrorResponse("internal server error", fiber.StatusInternalServerError).Send(ctx)
	}
	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h ParserHandler) GetParser(ctx fiber.Ctx) error {
	collegeId := fiber.Locals[uint](ctx, "collegeId")
	parserId := fiber.Locals[uint](ctx, "parserId")
	return ctx.JSON(dto.GetParserResponse{
		CollegeID: collegeId,
		ParserID:  parserId,
	})
}
