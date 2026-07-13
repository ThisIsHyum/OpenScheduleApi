package server

import (
	"errors"
	"time"

	"github.com/ThisIsHyum/OpenScheduleApi/internal/domain"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/dto"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/repository"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/service"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

const schedules = "/groups/:groupId/schedules"

type ScheduleHandler struct {
	scheduleService *service.ScheduleService
	logger          *logrus.Logger
}

func NewScheduleHandler(app *fiber.App,
	repos repository.Repos,
	logger *logrus.Logger) {
	handler := ScheduleHandler{
		scheduleService: service.NewScheduleService(
			repos.StudentGroupRepo, repos.LessonRepo, repos.CallRepo, repos.CollegeRepo),
		logger: logger}
	app.Get(schedules, handler.GetSchedules)
}

func (h ScheduleHandler) GetSchedules(ctx fiber.Ctx) error {
	dateStr := ctx.Query("date")
	week := ctx.Query("week")
	weekday := ctx.Query("weekday")
	day := ctx.Query("day")

	id := fiber.Params[uint](ctx, "groupId")
	if id == 0 {
		return dto.NewErrorResponse("invalid groupId", fiber.StatusBadRequest).Send(ctx)
	}

	if dateStr != "" && (week != "" || weekday != "" || day != "") {
		return dto.NewErrorResponse(
			"parameter 'date' cannot be combined with 'week', 'weekday' or 'day'",
			fiber.StatusConflict).Send(ctx)
	} else if day != "" && (week != "" || weekday != "" || dateStr != "") {
		return dto.NewErrorResponse(
			"parameter 'day' cannot be combined with 'week', 'weekday' or date",
			fiber.StatusConflict).Send(ctx)
	}

	if day != "" {
		resp, err := h.scheduleService.GetScheduleByDay(ctx, id, day)
		if errors.Is(err, domain.ErrInvalidData) {
			return dto.NewErrorResponse("invalid day (expected 'today' or 'tomorrow')", fiber.StatusBadRequest).Send(ctx)
		} else if err != nil {
			return dto.NewErrorResponse("internal server error", fiber.StatusInternalServerError).Send(ctx)
		}
		return ctx.JSON([]dto.ScheduleResponse{resp})
	}

	if dateStr != "" {
		date, err := time.Parse("02-01-2006", dateStr) // dd-mm-yyyy format
		if err != nil {
			return dto.NewErrorResponse("invalid date format (expected dd-mm-yyyy)", fiber.StatusBadRequest).Send(ctx)
		}
		resp, err := h.scheduleService.GetScheduleForDate(ctx, id, date)
		if err != nil {
			return dto.NewErrorResponse("internal server error", fiber.StatusInternalServerError).Send(ctx)
		}
		return ctx.JSON([]dto.ScheduleResponse{resp})
	}

	if weekday != "" {
		resp, err := h.scheduleService.GetScheduleByWeekday(ctx, id, weekday, week)
		if errors.Is(err, domain.ErrInvalidData) {
			return dto.NewErrorResponse(err.Error(), fiber.StatusBadRequest).Send(ctx)
		} else if err != nil {
			return dto.NewErrorResponse("internal server error", fiber.StatusInternalServerError).Send(ctx)
		}
		return ctx.JSON([]dto.ScheduleResponse{resp})
	}
	if week != "" {
		resp, err := h.scheduleService.GetSchedulesByWeek(ctx, id, week)
		if errors.Is(err, domain.ErrInvalidData) {
			return dto.NewErrorResponse(err.Error(), fiber.StatusBadRequest).Send(ctx)
		} else if err != nil {
			return dto.NewErrorResponse("internal server error", fiber.StatusInternalServerError).Send(ctx)
		}
		return ctx.JSON(resp)
	}
	return dto.NewErrorResponse("no valid parameters provided", fiber.StatusBadRequest).Send(ctx)
}
