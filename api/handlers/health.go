package handlers

import (
	types "crypto-braza-tokens-api/api/handlers/types"
	cfg "crypto-braza-tokens-api/configs"
	c "crypto-braza-tokens-api/constants"
	l "crypto-braza-tokens-api/utils/logger"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type HealthHandler struct {
	Resources *cfg.Resources
}

func (h HealthHandler) Health(ctx *fiber.Ctx) error {
	return h.HealthLiveness(ctx)
}

func (h HealthHandler) HealthLiveness(ctx *fiber.Ctx) error {
	result := &types.HealthStatus{
		App:    c.APP_NAME,
		Status: h.Resources.HealthService.Liveness(),
	}

	return ctx.Status(fiber.StatusOK).JSON(result)
}

func (h HealthHandler) HealthReadiness(ctx *fiber.Ctx) error {
	status, err := h.Resources.HealthService.Readiness(ctx.UserContext())
	if err != nil {
		l.Logger.Error("error retrieving readiness status from health service", zap.Error(err))
		InternalErrorWrapper(ctx, "blockchain", err)
	}

	result := &types.HealthStatus{
		App:    c.APP_NAME,
		Status: status,
	}

	return ctx.Status(fiber.StatusOK).JSON(result)
}
