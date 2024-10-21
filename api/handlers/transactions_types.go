package handlers

import (
	types "crypto-braza-tokens-api/api/handlers/types"
	cfg "crypto-braza-tokens-api/configs"
	l "crypto-braza-tokens-api/utils/logger"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type TransactionsTypesHandler struct {
	Resources *cfg.Resources
}

// GetTransactionsTypesNames retrieve the list of transactions types names
// @Summary Get the transactions types names list
// @Description retrieve the list of transactions types names
// @Tags TransactionsTypes
// @ID get-transactions-types-names
// @Produce json
// @Success 200 {object} fiber.Map
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/transactions-types/list [get]
func (t TransactionsTypesHandler) GetTransactionsTypesNames(ctx *fiber.Ctx) error {
	result, err := t.Resources.TransactionService.FindAllTypesMap(ctx.UserContext())
	if err != nil {
		if err.Error() != "no transaction types found" {
			return InternalErrorWrapper(ctx, "transaction type", err)
		}

		l.Logger.Info("handler: no transaction types found", zap.Error(err))
		return ObjectResultWrapper(ctx, fiber.Map{"transaction_types": []string{}})
	}

	return ObjectResultWrapper(ctx, result)
}

// GetTransactionsTypes retrieve the list of transactions types
// @Summary Get the transactions types list
// @Description retrieve the list of transactions types
// @Tags TransactionsTypes
// @ID get-transactions-types
// @Produce json
// @Success 200 {array} transaction.TransactionType
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/transactions-types [get]
func (t TransactionsTypesHandler) GetTransactionsTypes(ctx *fiber.Ctx) error {
	result, err := t.Resources.TransactionService.FindAllTypes(ctx.UserContext())
	if err != nil {
		if err.Error() != "no transaction types found" {
			return InternalErrorWrapper(ctx, "transaction type", err)
		}

		l.Logger.Info("handler: no transaction types found", zap.Error(err))
		return ObjectResultWrapper(ctx, fiber.Map{"transaction_types": []string{}})
	}

	return ObjectResultWrapper(ctx, result)
}

// GetTransactionTypeById retrieve a transaction type by id
// @Summary Get a transaction type
// @Description retrieve a transaction type by id
// @Tags TransactionsTypes
// @ID get-transaction-type-by-id
// @Produce json
// @Success 200 {object} transaction.TransactionType
// @Failure 400 {object} types.ErrorMessage
// @Failure 404 {object} types.ErrorMessage
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/transactions-types/{id} [get]
func (t TransactionsTypesHandler) GetTransactionTypeById(ctx *fiber.Ctx) error {
	err := ValidatePathParam(ctx, "id")
	if err != nil {
		return BadRequestWrapper(ctx, "transaction type", err)
	}

	result, err := t.Resources.TransactionService.FindTransactionTypeById(ctx.UserContext(), ctx.Params("id"))
	if err != nil {
		return InternalErrorWrapper(ctx, "transaction type", err)
	}

	if result == nil {
		l.Logger.Info("handler: no transaction type found", zap.Error(err))
		return NotFoundWrapper(ctx)
	}

	return ObjectResultWrapper(ctx, result)
}

// PostTransactionType create a new transaction type
// @Summary Create a new transaction type
// @Description create a new transaction type
// @Tags TransactionsTypes
// @ID post-transaction-type
// @Accept json
// @Produce json
// @Param body body types.SaveTransactionTypeRequest true "transaction Type"
// @Success 200 {object} types.Result
// @Failure 400 {object} types.ErrorMessage
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/transactions-types [post]
func (t TransactionsTypesHandler) PostTransactionType(ctx *fiber.Ctx) error {
	request := &types.SaveTransactionTypeRequest{}

	if err := request.FromBody(ctx); err != nil {
		return InternalErrorWrapper(ctx, "transaction type", err)
	}

	if err := request.IsValid(); err != nil {
		return BadRequestWrapper(ctx, "transaction type", err)
	}

	result, err := t.Resources.OperationService.SaveOperationType(ctx.UserContext(), request.Name)
	if err != nil {
		return InternalErrorWrapper(ctx, "transaction type", err)
	}

	return MessageResultWrapper(ctx, fmt.Sprintf("transaction type %s was saved with ID %s", request.Name, result.Hex()))
}

// DeleteTransactionType delete a transaction type by id
// @Summary Delete a transaction type
// @Description delete a transaction type
// @Tags TransactionsTypes
// @ID delete-transaction-type
// @Param id path string true "transaction Type ID"
// @Produce json
// @Success 200 {object} types.Result
// @Failure 400 {object} types.ErrorMessage
// @Failure 404 {object} types.ErrorMessage
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/transactions-types/{id} [delete]
func (t TransactionsTypesHandler) DeleteTransactionType(ctx *fiber.Ctx) error {
	err := ValidatePathParam(ctx, "id")
	if err != nil {
		return BadRequestWrapper(ctx, "transaction type", err)
	}

	result, err := t.Resources.TransactionService.DeleteTransactionType(ctx.UserContext(), ctx.Params("id"))
	if err != nil {
		return InternalErrorWrapper(ctx, "transaction type", err)
	}

	return MessageResultWrapper(ctx, result)
}

// PatchTransactionType update a transaction type by id
// @Summary Update a transaction type
// @Description update a transaction type
// @Tags TransactionsTypes
// @ID patch-transaction-type
// @Accept json
// @Produce json
// @Param request body types.SavetransactionTypeRequest true "transaction Type"
// @Success 200 {object} types.Result
// @Failure 400 {object} types.ErrorMessage
// @Failure 404 {object} types.ErrorMessage
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/transactions-types/{id} [patch]
func (t TransactionsTypesHandler) PatchTransactionType(ctx *fiber.Ctx) error {
	err := ValidatePathParam(ctx, "id")
	if err != nil {
		return BadRequestWrapper(ctx, "transaction type", err)
	}

	request := types.SaveTransactionTypeRequest{}

	if err := request.FromBody(ctx); err != nil {
		return BadRequestWrapper(ctx, "transaction type", err)
	}

	if err := request.IsValid(); err != nil {
		return BadRequestWrapper(ctx, "transaction type", err)
	}

	name := strings.ToUpper(request.Name)

	result, err := t.Resources.TransactionService.EditTransactionType(ctx.UserContext(), ctx.Params("id"), name)
	if err != nil {
		return InternalErrorWrapper(ctx, "transaction type", err)
	}

	return MessageResultWrapper(ctx, fmt.Sprintf("transaction type %s was saved with ID %s", name, result.Hex()))
}
