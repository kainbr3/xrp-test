package handlers

import (
	types "crypto-braza-tokens-api/api/handlers/types"
	l "crypto-braza-tokens-api/utils/logger"

	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// GetOperationTypesNames retrieve the list of operations types names
// @Summary Get the operations types names list
// @Description retrieve the list of operations types names
// @Tags OperationsTypes
// @ID get-operations-types-names
// @Produce json
// @Success 200 {object} fiber.Map
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/operations-types/list [get]
func (o OperationsHandler) GetOperationTypesNames(ctx *fiber.Ctx) error {
	result, err := o.Resources.OperationService.FindAllTypesMap(ctx.UserContext())
	if err != nil {
		if err.Error() != "no operation types found" {
			return InternalErrorWrapper(ctx, "operation type", err)
		}

		l.Logger.Info("handler: no operation types found", zap.Error(err))
		return ObjectResultWrapper(ctx, fiber.Map{"operation_types": []string{}})
	}

	return ObjectResultWrapper(ctx, result)
}

// GetOperationsTypes retrieve the list of operations types
// @Summary Get the operations types list
// @Description retrieve the list of operations types
// @Tags OperationsTypes
// @ID get-operations-types
// @Produce json
// @Success 200 {array} operation.OperationType
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/operations-types [get]
func (o OperationsHandler) GetOperationTypes(ctx *fiber.Ctx) error {
	result, err := o.Resources.OperationService.FindAllTypes(ctx.UserContext())
	if err != nil {
		if err.Error() != "no operation types found" {
			return InternalErrorWrapper(ctx, "operation type", err)
		}

		l.Logger.Info("handler: no operation types found", zap.Error(err))
		return ObjectResultWrapper(ctx, fiber.Map{"operation_types": []string{}})

	}

	return ObjectResultWrapper(ctx, result)
}

// GetOperationTypeById retrieve a operation type by id
// @Summary Get a operation type
// @Description retrieve a operation type by id
// @Tags OperationsTypes
// @ID get-operation-type-by-id
// @Produce json
// @Success 200 {object} operation.OperationType
// @Failure 400 {object} types.ErrorMessage
// @Failure 404 {object} types.ErrorMessage
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/operations-types/{id} [get]
func (o OperationsHandler) GetOperationTypeById(ctx *fiber.Ctx) error {
	err := ValidatePathParam(ctx, "id")
	if err != nil {
		return BadRequestWrapper(ctx, "operation type", err)
	}

	result, err := o.Resources.OperationService.FindOperationTypeById(ctx.UserContext(), ctx.Params("id"))
	if err != nil {
		return InternalErrorWrapper(ctx, "operation type", err)
	}

	if result == nil {
		l.Logger.Info("handler: no operation type found", zap.Error(err))
		return NotFoundWrapper(ctx)
	}

	return ObjectResultWrapper(ctx, result)
}

// PostOperationType create a new operation type
// @Summary Create a new operation type
// @Description create a new operation type
// @Tags OperationsTypes
// @ID post-operation-type
// @Accept json
// @Produce json
// @Param body body types.SaveOperationTypeRequest true "Operation Type"
// @Success 200 {object} types.Result
// @Failure 400 {object} types.ErrorMessage
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/operations-types [post]
func (o OperationsHandler) PostOperationType(ctx *fiber.Ctx) error {
	request := types.SaveOperationTypeRequest{}

	if err := request.FromBody(ctx); err != nil {
		return BadRequestWrapper(ctx, "operation type", err)
	}

	if err := request.IsValid(); err != nil {
		return BadRequestWrapper(ctx, "operation type", err)
	}

	result, err := o.Resources.OperationService.SaveOperationType(ctx.UserContext(), request.Name)
	if err != nil {
		return InternalErrorWrapper(ctx, "operation type", err)
	}

	return MessageResultWrapper(ctx, fmt.Sprintf("operation type %s was saved with ID %s", request.Name, result.Hex()))
}

// DeleteOperationType delete a operation type by id
// @Summary Delete a operation type
// @Description delete a operation type
// @Tags OperationsTypes
// @ID delete-operation-type
// @Param id path string true "Operation Type ID"
// @Produce json
// @Success 200 {object} types.Result
// @Failure 400 {object} types.ErrorMessage
// @Failure 404 {object} types.ErrorMessage
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/operations-types/{id} [delete]
func (o OperationsHandler) DeleteOperationType(ctx *fiber.Ctx) error {
	err := ValidatePathParam(ctx, "id")
	if err != nil {
		return BadRequestWrapper(ctx, "operation type", err)
	}

	result, err := o.Resources.OperationService.DeleteOperationType(ctx.UserContext(), ctx.Params("id"))
	if err != nil {
		return InternalErrorWrapper(ctx, "operation type", err)
	}

	return MessageResultWrapper(ctx, result)
}

// PatchOperationType update a operation type by id
// @Summary Update a operation type
// @Description update a operation type
// @Tags OperationsTypes
// @ID patch-operation-type
// @Accept json
// @Produce json
// @Param request body types.SaveOperationTypeRequest true "Operation Type"
// @Success 200 {object} types.Result
// @Failure 400 {object} types.ErrorMessage
// @Failure 404 {object} types.ErrorMessage
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/operations-types [patch]
func (o OperationsHandler) PatchOperationType(ctx *fiber.Ctx) error {
	request := types.SaveOperationTypeRequest{}

	if err := request.FromBody(ctx); err != nil {
		return BadRequestWrapper(ctx, "operation type", err)
	}

	if err := request.IsValid(); err != nil {
		return BadRequestWrapper(ctx, "operation type", err)
	}

	name := strings.ToUpper(request.Name)

	result, err := o.Resources.OperationService.EditOperationType(ctx.UserContext(), name)
	if err != nil {
		return InternalErrorWrapper(ctx, "operation type", err)
	}

	return MessageResultWrapper(ctx, fmt.Sprintf("operation type %s was saved with ID %s", name, result.Hex()))
}
