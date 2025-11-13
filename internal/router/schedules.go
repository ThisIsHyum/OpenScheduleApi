package router

import (
	"time"

	"github.com/ThisIsHyum/OpenScheduleApi/internal/database"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/dto"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/models"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/weekdays"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

const schedules = "/groups/:groupId/schedules"

type ScheduleHandler struct {
	db     *database.Db
	logger *logrus.Logger
}

func NewScheduleHandler(app *fiber.App, db *database.Db, logger *logrus.Logger) {
	handler := ScheduleHandler{db: db, logger: logger}
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
		switch day {
		case "today":
			return h.getScheduleForDate(ctx, time.Now(), id)
		case "tomorrow":
			return h.getScheduleForDate(ctx, time.Now().AddDate(0, 0, 1), id)
		default:
			return dto.NewErrorResponse("invalid day (expected 'today' or 'tomorrow')", fiber.StatusBadRequest).Send(ctx)
		}
	}

	if dateStr != "" {
		date, err := time.Parse("02-01-2006", dateStr) // dd-mm-yyyy format
		if err != nil {
			return dto.NewErrorResponse("invalid date format (expected dd-mm-yyyy)", fiber.StatusBadRequest).Send(ctx)
		}
		return h.getScheduleForDate(ctx, date, id)
	}

	if weekday != "" {
		w, ok := weekdays.ParseWeekday(weekday)
		if !ok {
			return dto.NewErrorResponse("invalid weekday", fiber.StatusBadRequest).Send(ctx)
		}
		switch week {
		case "previous":
			return h.GetScheduleByWeekday(ctx, w, -1, id)
		case "", "current":
			return h.GetScheduleByWeekday(ctx, w, 0, id)
		case "next":
			return h.GetScheduleByWeekday(ctx, w, 1, id)
		default:
			return dto.NewErrorResponse("invalid week (expected previous, current, next)", fiber.StatusBadRequest).Send(ctx)
		}
	}
	if week != "" {
		switch week {
		case "previous":
			return h.getScheduleForWeek(ctx, -1, id)
		case "current":
			return h.getScheduleForWeek(ctx, 0, id)
		case "next":
			return h.getScheduleForWeek(ctx, 1, id)
		default:
			return dto.NewErrorResponse("invalid week (expected previous, current, next)", fiber.StatusBadRequest).Send(ctx)
		}
	}
	return dto.NewErrorResponse("no valid parameters provided", fiber.StatusBadRequest).Send(ctx)
}

func (h ScheduleHandler) GetScheduleByWeekday(ctx fiber.Ctx, weekday time.Weekday, offset int, id uint) error {
	date := weekdays.GetDateByWeekday(weekday, offset)
	return h.getScheduleForDate(ctx, date, id)
}

func (h ScheduleHandler) getScheduleForDate(ctx fiber.Ctx, date time.Time, id uint) error {
	lessons, err := h.db.GetLessonsForDate(id, date)
	if err != nil {
		h.logger.WithError(err).Error("unable to get lessons")
		return dto.NewErrorResponse("internal server error", fiber.StatusInternalServerError).Send(ctx)
	}
	if len(lessons) == 0 {
		return ctx.JSON([]models.Schedule{})
	}

	calls, err := h.getCallsByGroupID(id)
	if err != nil {
		h.logger.WithError(err).Error("unable to get calls")
		return dto.NewErrorResponse("internal server error", fiber.StatusInternalServerError).Send(ctx)
	}

	schedules := []models.Schedule{models.NewSchedule(lessons, calls)}
	return ctx.JSON(schedules)
}

func (h ScheduleHandler) getScheduleForWeek(ctx fiber.Ctx, offset int, id uint) error {
	startDate, endDate := weekdays.WeekBounds(offset)
	lessons, err := h.db.GetLessonsForDates(id, startDate, endDate)
	if err != nil {
		h.logger.WithError(err).Error("unable to get lessons")
		return dto.NewErrorResponse("internal server error", fiber.StatusInternalServerError).Send(ctx)
	}

	calls, err := h.getCallsByGroupID(id)
	if err != nil {
		h.logger.WithError(err).Error("unable to get calls")
		return dto.NewErrorResponse("internal server error", fiber.StatusInternalServerError).Send(ctx)
	}

	schedules := models.NewSchedules(lessons, calls)
	return ctx.JSON(schedules)
}

func (h ScheduleHandler) getCallsByGroupID(id uint) ([]models.Call, error) {
	collegeID, err := h.db.GetCollegeIDByGroupID(id)
	if err != nil {
		return nil, err
	}
	calls, err := h.db.GetCalls(collegeID)
	return calls, nil
}
