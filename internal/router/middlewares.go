package router

import (
	"errors"
	"strings"

	"github.com/ThisIsHyum/OpenScheduleApi/internal/database"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/dto"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type MiddlewareHandler struct {
	db         *database.Db
	logger     *logrus.Logger
	adminToken string
}

func (h MiddlewareHandler) AdminAuthMiddleware(ctx fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return dto.NewErrorResponse("missing authorization header",
			fiber.StatusUnauthorized).Send(ctx)
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return dto.NewErrorResponse("invalid authorization header format",
			fiber.StatusUnauthorized).Send(ctx)
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	if h.adminToken != token {
		return dto.NewErrorResponse("invalid authorization token",
			fiber.StatusUnauthorized).Send(ctx)
	}
	return ctx.Next()
}

func (h MiddlewareHandler) ParserAuthMiddleware(ctx fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return dto.NewErrorResponse("missing authorization header",
			fiber.StatusUnauthorized).Send(ctx)
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return dto.NewErrorResponse("invalid authorization header format",
			fiber.StatusUnauthorized).Send(ctx)
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	parser, err := h.db.GetParserByToken(token)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.NewErrorResponse("wrong authorization token",
			fiber.StatusUnauthorized).Send(ctx)
	} else if err != nil {
		h.logger.WithError(err).Error("unable to get parser")
		return dto.NewErrorResponse("internal server error", fiber.StatusInternalServerError).Send(ctx)
	}

	fiber.Locals(ctx, "collegeId", parser.College.CollegeID)
	fiber.Locals(ctx, "parserId", parser.ParserID)
	return ctx.Next()
}
