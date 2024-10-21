package handlers

import (
	types "crypto-braza-tokens-api/api/handlers/types"
	l "crypto-braza-tokens-api/utils/logger"

	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// GetOperationsDomains retrieve the list of operations domains names
// @Summary Get the operations domains names list
// @Description retrieve the list of operations domains names
// @Tags OperationsDomains
// @ID get-operations-domains-names
// @Produce json
// @Success 200 {object} fiber.Map
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/operations-domains/list [get]
func (o OperationsHandler) GetOperationDomainsNames(ctx *fiber.Ctx) error {
	result, err := o.Resources.OperationService.GetOperationDomainsNames(ctx.UserContext())
	if err != nil {
		if err.Error() != "no operation domains found" {
			return InternalErrorWrapper(ctx, "operation domain", err)
		}

		l.Logger.Info("handler: no operation domains found", zap.Error(err))
		return ObjectResultWrapper(ctx, fiber.Map{"operation_domains": []string{}})
	}

	return ObjectResultWrapper(ctx, result)
}

// GetOperationsDomains retrieve the list of operations domains
// @Summary Get the operations domains list
// @Description retrieve the list of operations domains
// @Tags OperationsDomains
// @ID get-operations-domains
// @Produce json
// @Success 200 {array} operation.OperationDomain
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/operations-domains [get]
func (o OperationsHandler) GetOperationDomains(ctx *fiber.Ctx) error {
	result, err := o.Resources.OperationService.FindAllDomains(ctx.UserContext())
	if err != nil {
		if err.Error() != "no operation domains found" {
			return InternalErrorWrapper(ctx, "operation domain", err)
		}

		l.Logger.Info("handler: no operation domains found", zap.Error(err))
		return ObjectResultWrapper(ctx, fiber.Map{"operation_domains": []string{}})
	}

	return ObjectResultWrapper(ctx, result)
}

// GetOperationDomainById retrieve a operation domain by id
// @Summary Get a operation domain
// @Description retrieve a operation domain by id
// @Tags OperationsDomains
// @ID get-operation-domain-by-id
// @Produce json
// @Success 200 {object} operation.OperationDomain
// @Failure 400 {object} types.ErrorMessage
// @Failure 404 {object} types.ErrorMessage
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/operations-domains/{id} [get]
func (o OperationsHandler) GetOperationDomainById(ctx *fiber.Ctx) error {
	err := ValidatePathParam(ctx, "id")
	if err != nil {
		return BadRequestWrapper(ctx, "operation domain", err)
	}

	result, err := o.Resources.OperationService.FindOperationDomainById(ctx.UserContext(), ctx.Params("id"))
	if err != nil {
		return InternalErrorWrapper(ctx, "operation domain", err)
	}

	if result == nil {
		l.Logger.Info("handler: no operation domain found", zap.Error(err))
		return NotFoundWrapper(ctx)
	}

	return ObjectResultWrapper(ctx, result)
}

// PostOperationDomain create a new operation domain
// @Summary Create a new operation domain
// @Description create a new operation domain
// @Tags OperationsDomains
// @ID post-operation-domain
// @Accept json
// @Produce json
// @Param body body types.SaveOperationDomainRequest true "Operation Domain"
// @Success 200 {object} types.Result
// @Failure 400 {object} types.ErrorMessage
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/operations-domains [post]
func (o OperationsHandler) PostOperationDomain(ctx *fiber.Ctx) error {
	request := types.SaveOperationDomainRequest{}

	if err := request.FromBody(ctx); err != nil {
		return BadRequestWrapper(ctx, "operation domain", err)
	}

	if err := request.IsValid(); err != nil {
		return BadRequestWrapper(ctx, "operation domain", err)
	}

	result, err := o.Resources.OperationService.SaveOperationDomain(ctx.UserContext(), request.Name, request.IsActive)
	if err != nil {
		return InternalErrorWrapper(ctx, "operation domain", err)
	}

	return MessageResultWrapper(ctx, result.Hex())
}

// DeleteOperationDomain delete a operation domain by id
// @Summary Delete a operation domain
// @Description delete a operation domain
// @Tags OperationsDomains
// @ID delete-operation-domain
// @Param id path string true "Operation Domain ID"
// @Produce json
// @Success 200 {object} types.Result
// @Failure 400 {object} types.ErrorMessage
// @Failure 404 {object} types.ErrorMessage
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/operations-domains/{id} [delete]
func (o OperationsHandler) DeleteOperationDomain(ctx *fiber.Ctx) error {
	err := ValidatePathParam(ctx, "id")
	if err != nil {
		return BadRequestWrapper(ctx, "operation domain", err)
	}

	result, err := o.Resources.OperationService.DeleteOperationDomain(ctx.UserContext(), ctx.Params("id"))
	if err != nil {
		return InternalErrorWrapper(ctx, "operation domain", err)
	}

	return MessageResultWrapper(ctx, result)
}

// PatchOperationDomain update a operation domain by id
// @Summary Update a operation domain
// @Description update a operation domain
// @Tags OperationsDomains
// @ID patch-operation-domain
// @Accept json
// @Produce json
// @Param request body types.EditOperationDomainRequest true "Operation Domain"
// @Success 200 {object} types.Result
// @Failure 400 {object} types.ErrorMessage
// @Failure 404 {object} types.ErrorMessage
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/operations-domains [patch]
func (o OperationsHandler) PatchOperationDomain(ctx *fiber.Ctx) error {
	request := types.EditOperationDomainRequest{}

	if err := request.FromBody(ctx); err != nil {
		return BadRequestWrapper(ctx, "operation domain", err)
	}

	if err := request.IsValid(); err != nil {
		return BadRequestWrapper(ctx, "operation domain", err)
	}

	result, err := o.Resources.OperationService.EditOperationDomain(ctx.UserContext(), request.Name, request.IsActive)
	if err != nil {
		return InternalErrorWrapper(ctx, "operation domain", err)
	}

	return MessageResultWrapper(ctx, fmt.Sprintf("operation domain %s was saved with ID %s", request.Name, result.Hex()))
}
