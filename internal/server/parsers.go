package server

import (
	"time"

	"github.com/ThisIsHyum/OpenScheduleApi/internal/domain"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/dto"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/repository"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/service"
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
	service *service.ParserService
	logger  *logrus.Logger
}

func NewParserHandler(app *fiber.App,
	repos repository.Repos, logger *logrus.Logger) {
	parserService := service.NewParserService(
		repos.CallRepo, repos.StudentGroupRepo,
		repos.LessonRepo, repos.CampusRepo, repos.CollegeRepo)

	handler := ParserHandler{logger: logger, service: parserService}
	mh := MiddlewareHandler{parserService: parserService, logger: logger}
	app.Group(parser, mh.ParserAuthMiddleware).
		Get(getParser, handler.GetParser).
		Post(updateCalls, handler.UpdateCalls).
		Post(addLessons, handler.AddLessons).
		Post(updateGroups, handler.UpdateGroups)
}

// @Summary get parser information
// @Description get information about the college associated with the parser token
// @Tags parser
// @Security ParserBearerAuth
// @Produce json
// @Success 200 {object} dto.GetParserResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /parser [get]
func (h ParserHandler) GetParser(ctx fiber.Ctx) error {
	collegeId := fiber.Locals[uint](ctx, "collegeId")
	return ctx.JSON(dto.GetParserResponse{CollegeID: collegeId})
}

// @Summary update calls of parser college
// @Description update calls of parser college
// @Tags parser
// @Security ParserBearerAuth
// @Accept json
// @Produce json
// @Param calls body []dto.Call true "Calls"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /parser/calls [post]
func (h ParserHandler) UpdateCalls(ctx fiber.Ctx) error {
	c := ctx.Context()
	collegeId := fiber.Locals[uint](ctx, "collegeId")
	var calls []dto.Call
	if err := ctx.Bind().Body(&calls); err != nil {
		return dto.NewErrorResponse("invalid request body", fiber.StatusBadRequest).Send(ctx)
	}

	if len(calls) == 0 {
		return dto.NewErrorResponse("empty calls list", fiber.StatusBadRequest).Send(ctx)
	}

	domainCalls := make([]domain.Call, 0, len(calls))
	for _, call := range calls {
		domainCalls = append(domainCalls, domain.Call{
			ID:      call.CallID,
			Weekday: call.Weekday,
			Begins:  time.Time(call.Begins),
			Ends:    time.Time(call.Ends),
			Order:   call.Order,
		})
	}
	if err := h.service.UpdateCalls(c, collegeId, domainCalls); err != nil {
		h.logger.WithError(err).Error("unable to update calls")
		return dto.NewErrorResponse("internal server error", fiber.StatusInternalServerError).Send(ctx)
	}
	return ctx.SendStatus(fiber.StatusNoContent)
}

// @Summary add lessons to parser college
// @Description add lessons to parser college
// @Tags parser
// @Security ParserBearerAuth
// @Accept json
// @Produce json
// @Param lessons body []dto.Lesson true "Lessons"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /parser/lessons [post]
func (h ParserHandler) AddLessons(ctx fiber.Ctx) error {
	c := ctx.Context()

	var lessons []dto.Lesson
	if err := ctx.Bind().Body(&lessons); err != nil {
		return dto.NewErrorResponse("invalid request body", fiber.StatusBadRequest).Send(ctx)
	}

	if len(lessons) == 0 {
		return dto.NewErrorResponse("empty lessons list", fiber.StatusBadRequest).Send(ctx)
	}

	domainLessons := make([]domain.Lesson, 0, len(lessons))
	for _, lesson := range lessons {
		domainLessons = append(domainLessons, domain.Lesson{
			ID:             lesson.LessonID,
			Title:          lesson.Title,
			Cabinet:        lesson.Cabinet,
			Date:           lesson.Date.Time,
			Teacher:        lesson.Teacher,
			Order:          lesson.Order,
			StudentGroupID: lesson.StudentGroupID,
		})
	}

	if err := h.service.AddLessons(c, domainLessons); err != nil {
		h.logger.WithError(err).Error("unable to add lessons")
		return dto.NewErrorResponse("internal server error", fiber.StatusInternalServerError).Send(ctx)
	}
	return ctx.SendStatus(fiber.StatusNoContent)
}

// @Summary update groups of parser college
// @Description update groups of parser college
// @Tags parser
// @Security ParserBearerAuth
// @Accept json
// @Produce json
// @Param updateGroupsRequest body dto.UpdateGroupsRequest true "Update groups request"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /parser/groups [post]
func (h ParserHandler) UpdateGroups(ctx fiber.Ctx) error {
	c := ctx.Context()
	var requestBody dto.UpdateGroupsRequest
	if err := ctx.Bind().Body(&requestBody); err != nil {
		return dto.NewErrorResponse("invalid request body", fiber.StatusBadRequest).Send(ctx)
	}

	if requestBody.CampusID == 0 || len(requestBody.StudentGroupNames) == 0 {
		return dto.NewErrorResponse("invalid request data", fiber.StatusBadRequest).Send(ctx)
	}
	if err := h.service.UpdateGroups(c, requestBody.CampusID, requestBody.StudentGroupNames); err != nil {
		h.logger.WithError(err).Error("unable to update groups")
		return dto.NewErrorResponse("internal server error", fiber.StatusInternalServerError).Send(ctx)
	}
	return ctx.SendStatus(fiber.StatusNoContent)
}
