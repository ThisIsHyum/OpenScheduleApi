package server

import (
	"errors"
	"net/url"

	"github.com/ThisIsHyum/OpenScheduleApi/internal/domain"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/dto"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/repository"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/service"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

const (
	getGroupsByCollegeId = "/colleges/:collegeId/groups"
	getGroupsByCampusId  = "/campuses/:campusId/groups"
	getGroup             = "/groups/:id"
)

type GroupHandler struct {
	groupService *service.StudentGroupService
	logger       *logrus.Logger
}

func NewGroupHandler(app *fiber.App,
	repos repository.Repos, logger *logrus.Logger) {
	groupService := service.NewStudentGroupService(repos.StudentGroupRepo, repos.CampusRepo, repos.CollegeRepo)
	handler := GroupHandler{logger: logger, groupService: groupService}
	app.Get(getGroupsByCampusId, handler.GetGroupsByCampusID)
	app.Get(getGroupsByCollegeId, handler.GetGroupsByCollegeID)
	app.Get(getGroup, handler.GetGroup)
}

// @Summary get groups by campus ID
// @Description get all groups by campus ID or a group by campus ID with specified name
// @Tags groups
// @Produce json
// @Param campusId path  int    true  "Campus ID"
// @Param name      query string false  "Group name"
// @Success 200 {array}  dto.StudentGroupResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /campuses/{campusId}/groups [get]
func (h GroupHandler) GetGroupsByCampusID(ctx fiber.Ctx) error {
	c := ctx.Context()
	name := ctx.Query("name")
	id := fiber.Params[uint](ctx, "campusId")
	if id == 0 {
		return dto.NewErrorResponse("invalid campusId", fiber.StatusBadRequest).Send(ctx)
	}

	groups, err := h.groupService.GetGroups(c, id, name)
	if errors.Is(err, domain.ErrNotFound) {
		return dto.NewErrorResponse("campus not found", fiber.StatusNotFound).Send(ctx)
	} else if err != nil {
		h.logger.WithError(err).Error("unable to get groups")
		return dto.NewErrorResponse("internal server error", fiber.StatusInternalServerError).Send(ctx)
	}
	return ctx.JSON(groups)
}

// @Summary get groups by college ID
// @Description get all groups by college ID or a group by college ID with specified name
// @Tags groups
// @Produce json
// @Param collegeId path  int    true  "College ID"
// @Param name      query string false  "Group name"
// @Success 200 {array}  dto.StudentGroupResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /colleges/{collegeId}/groups [get]
func (h GroupHandler) GetGroupsByCollegeID(ctx fiber.Ctx) error {
	name, err := url.QueryUnescape(ctx.Query("name"))
	if err != nil {
		return dto.NewErrorResponse("invalid name", fiber.StatusBadRequest).Send(ctx)
	}
	id := fiber.Params[uint](ctx, "collegeId")
	if id == 0 {
		return dto.NewErrorResponse("invalid collegeId", fiber.StatusBadRequest).Send(ctx)
	}

	groups, err := h.groupService.GetGroupsByCollegeID(ctx, id, name)
	if errors.Is(err, domain.ErrNotFound) {
		return dto.NewErrorResponse("college not found", fiber.StatusNotFound).Send(ctx)
	} else if err != nil {
		h.logger.WithError(err).Error("unable to get groups")
		return dto.NewErrorResponse("internal server error", fiber.StatusInternalServerError).Send(ctx)
	}
	return ctx.JSON(groups)
}

// @Summary get a group by ID
// @Description get a single group by its ID
// @Tags groups
// @Produce json
// @Param id path int true "Group ID"
// @Success 200 {object} dto.StudentGroupResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /groups/{id} [get]
func (h GroupHandler) GetGroup(ctx fiber.Ctx) error {
	id := fiber.Params[uint](ctx, "id")
	if id == 0 {
		return dto.NewErrorResponse("invalid id", fiber.StatusBadRequest).Send(ctx)
	}
	group, err := h.groupService.GetGroup(ctx, id)
	if errors.Is(err, domain.ErrNotFound) {
		return dto.NewErrorResponse("group not found", fiber.StatusNotFound).Send(ctx)
	} else if err != nil {
		h.logger.WithError(err).Error("unable to get group")
		return dto.NewErrorResponse("internal server error", fiber.StatusInternalServerError).Send(ctx)
	}
	return ctx.JSON(group)
}
