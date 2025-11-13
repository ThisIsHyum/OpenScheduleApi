package router

import (
	"errors"
	"net/url"

	"github.com/ThisIsHyum/OpenScheduleApi/internal/database"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/dto"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/models"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	getGroupsByCollegeId = "/colleges/:collegeId/groups"
	getGroupsByCampusId  = "/campuses/:campusId/groups"
	getGroup             = "/groups/:groupId"
)

type GroupHandler struct {
	db     *database.Db
	logger *logrus.Logger
}

func NewGroupHandler(app *fiber.App, db *database.Db, logger *logrus.Logger) {
	handler := GroupHandler{db: db, logger: logger}
	app.Get(getGroupsByCampusId, handler.GetGroupsByCampusID)
	app.Get(getGroupsByCollegeId, handler.GetGroupsByCollegeID)
	app.Get(getGroup, handler.GetGroup)
}
func (h GroupHandler) GetGroupsByCampusID(ctx fiber.Ctx) error {
	name := ctx.Query("name")
	id := fiber.Params[uint](ctx, "campusId")
	if id == 0 {
		return dto.NewErrorResponse("invalid campusId", fiber.StatusBadRequest).Send(ctx)
	}

	var groups []models.StudentGroup
	var err error

	if _, err := h.db.GetCampusByID(id); errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.NewErrorResponse("campus not found", fiber.StatusNotFound).Send(ctx)
	} else if err != nil {
		h.logger.WithError(err).Error("unable to get campus")
		return dto.NewErrorResponse("internal server error", fiber.StatusInternalServerError).Send(ctx)
	}
	if name != "" {
		groups, err = h.db.GetGroupsByName(id, name)
	} else {
		groups, err = h.db.GetGroupsByCampusID(id)
	}
	if err != nil {
		h.logger.WithError(err).Error("unable to get groups")
		return dto.NewErrorResponse("internal server error", fiber.StatusInternalServerError).Send(ctx)
	}
	return ctx.JSON(groups)
}

func (h GroupHandler) GetGroupsByCollegeID(ctx fiber.Ctx) error {
	name, err := url.QueryUnescape(ctx.Query("name"))
	if err != nil {
		return dto.NewErrorResponse("invalid name", fiber.StatusBadRequest).Send(ctx)
	}
	id := fiber.Params[uint](ctx, "collegeId")
	if id == 0 {
		return dto.NewErrorResponse("invalid collegeId", fiber.StatusBadRequest).Send(ctx)
	}

	var groups []models.StudentGroup

	if name != "" {
		groups, err = h.db.GetGroupsByName(id, name)
	} else {
		groups, err = h.db.GetGroupsByCollegeID(id)
	}
	if err != nil {
		h.logger.WithError(err).Error("unable to get groups")
		return dto.NewErrorResponse("internal server error", fiber.StatusInternalServerError).Send(ctx)
	}
	return ctx.JSON(groups)
}

func (h GroupHandler) GetGroup(ctx fiber.Ctx) error {
	id := fiber.Params[uint](ctx, "groupId")
	if id == 0 {
		return dto.NewErrorResponse("invalid groupId", fiber.StatusBadRequest).Send(ctx)
	}
	group, err := h.db.GetGroupByID(id)
	if err != nil {
		h.logger.WithError(err).Error("unable to get group")
		return dto.NewErrorResponse("internal server error", fiber.StatusInternalServerError).Send(ctx)
	} else if group.StudentGroupID == 0 {
		return dto.NewErrorResponse("group not found", fiber.StatusNotFound).Send(ctx)
	}
	return ctx.JSON(group)
}
