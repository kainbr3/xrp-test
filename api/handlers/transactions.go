package handlers

import (
	types "crypto-braza-tokens-api/api/handlers/types"
	cfg "crypto-braza-tokens-api/configs"
	l "crypto-braza-tokens-api/utils/logger"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type TransactionsHandler struct {
	Resources *cfg.Resources
}

func (t TransactionsHandler) GetTransactions(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "ok"})
}

func (t TransactionsHandler) PostTransactions(ctx *fiber.Ctx) error {
	request := &types.InternalTransferRequest{}

	if err := request.FromBody(ctx); err != nil {
		return InternalErrorWrapper(ctx, "transaction", err)
	}

	if err := request.IsValid(); err != nil {
		return BadRequestWrapper(ctx, "transaction", err)
	}

	result, err := t.Resources.TransactionService.ExecuteInternalTransaction(ctx.UserContext(), request.Domain, request.Type, request.BlockchainId, request.AssetId, request.Amount, request.ExternalId)
	if err != nil {
		return InternalErrorWrapper(ctx, "transaction", err)
	}

	// operator := ctx.Get("client-id")
	// if operator == "" {
	// 	return BadRequestWrapper(ctx, "transaction", errors.New("required head missing: client-id"))
	// }

	return ctx.Status(fiber.StatusOK).JSON(result)
}

func (t TransactionsHandler) PostWebhook(ctx *fiber.Ctx) error {
	l.Logger.Info("fireblocks callback received", zap.ByteString("body", ctx.Body()))
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "ok"})
}
