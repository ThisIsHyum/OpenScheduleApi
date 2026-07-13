package server

import (
	"errors"
	"strings"

	"github.com/ThisIsHyum/OpenScheduleApi/internal/domain"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/dto"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/service"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

const localCollegeID = "collegeId"

type MiddlewareHandler struct {
	parserService *service.ParserService
	logger        *logrus.Logger
	adminToken    string
}

func (h MiddlewareHandler) AdminAuthMiddleware(ctx fiber.Ctx) error {
	token, err := h.extractBearerToken(ctx)
	if err != nil {
		return dto.NewErrorResponse(err.Error(), fiber.StatusUnauthorized).Send(ctx)
	}

	if h.adminToken != token {
		return dto.NewErrorResponse("invalid authorization token",
			fiber.StatusUnauthorized).Send(ctx)
	}
	return ctx.Next()
}

func (h MiddlewareHandler) ParserAuthMiddleware(ctx fiber.Ctx) error {
	c := ctx.Context()

	token, err := h.extractBearerToken(ctx)
	if err != nil {
		return dto.NewErrorResponse(err.Error(), fiber.StatusUnauthorized).Send(ctx)
	}
	parser, err := h.parserService.GetByToken(c, token)

	if errors.Is(err, domain.ErrNotFound) {
		return dto.NewErrorResponse("wrong authorization token",
			fiber.StatusUnauthorized).Send(ctx)
	} else if err != nil {
		h.logger.WithError(err).Error("unable to get parser")
		return dto.NewErrorResponse("internal server error", fiber.StatusInternalServerError).Send(ctx)
	}

	fiber.Locals(ctx, localCollegeID, parser.ID)
	return ctx.Next()
}

func (h MiddlewareHandler) extractBearerToken(ctx fiber.Ctx) (string, error) {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("missing authorization header")
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("missing authorization header")
	}

	return strings.TrimPrefix(authHeader, "Bearer "), nil
}
