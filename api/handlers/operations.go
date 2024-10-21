package handlers

import (
	"crypto-braza-tokens-api/api/handlers/types"
	cfg "crypto-braza-tokens-api/configs"
	r "crypto-braza-tokens-api/repositories"
	"errors"
	"fmt"
	"sync"

	"github.com/gofiber/fiber/v2"
)

var (
	operationMutex     sync.Mutex
	isOperationRunning bool
)

type OperationsHandler struct {
	Resources *cfg.Resources
}

// GetOperations retrieve the list of operations
// @Summary Get the operations list
// @Description retrieve the list of operations
// @Tags Operations
// @ID get-operations
// @Produce json
// @Param filter_param query string false "Filter parameter"
// @Param filter_value query string false "Filter value"
// @Param sort_field query string false "Sort field"
// @Param sort_order query string false "Sort order"
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Success 200 {object} repositories.PaginatedOperations
// @Failure 400 {object} types.ErrorMessage
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/operations [get]
func (o OperationsHandler) GetOperations(ctx *fiber.Ctx) error {
	params := &r.QueryParams{
		FilterParam: ctx.Query("filter_param", ""),
		FilterValue: ctx.Query("filter_value", ""),
		SortField:   ctx.Query("sort_field", "updated_at"),
		SortOrder:   ctx.Query("sort_order", "desc"),
		Page:        ctx.QueryInt("page", 1),
		Limit:       ctx.QueryInt("limit", 10),
	}

	result, err := o.Resources.OperationService.GetPaginatedOperations(ctx.UserContext(), params)
	if err != nil {
		if err.Error() == "no operations found" {
			return ctx.Status(fiber.StatusOK).JSON(&r.PaginatedOperations{})
		}
		return InternalErrorWrapper(ctx, "operation", err)
	}

	return ObjectResultWrapper(ctx, result)
}

// GetOperationById retrieve an operation by id
// @Summary Get an operation
// @Description retrieve an operation by id
// @Tags Operations
// @ID get-operation-by-id
// @Produce json
// @Param id path string true "Operation ID"
// @Success 200 {object} operation.OperationWithLogs
// @Failure 400 {object} types.ErrorMessage
// @Failure 404 {object} types.ErrorMessage
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/operations/{id} [get]
func (o OperationsHandler) GetOperationById(ctx *fiber.Ctx) error {
	err := ValidatePathParam(ctx, "id")
	if err != nil {
		return BadRequestWrapper(ctx, "blockchain", err)
	}

	result, err := o.Resources.OperationService.GetOperationById(ctx.UserContext(), ctx.Params("id"))
	if err != nil {
		return InternalErrorWrapper(ctx, "operation", err)
	}

	if result == nil {
		return BadRequestWrapper(ctx, "operation", errors.New("operation not found"))
	}

	return ObjectResultWrapper(ctx, result)
}

// PostOperation create a new operation
// @Summary Create a new operation
// @Description create a new operation
// @Tags Operations
// @ID post-operation
// @Accept json
// @Produce json
// @Param operation body types.OperationRequest true "Operation object"
// @Success 200 {object} types.OperationResponse
// @Failure 400 {object} types.ErrorMessage
// @Failure 409 {object} types.ErrorMessage
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/operations [post]
func (o OperationsHandler) PostOperation(ctx *fiber.Ctx) error {
	request := types.OperationRequest{}

	if err := request.FromBody(ctx); err != nil {
		return InternalErrorWrapper(ctx, "operation", err)
	}

	if err := request.IsValid(); err != nil {
		return BadRequestWrapper(ctx, "blockchain", err)
	}

	if err := o.Resources.OperationService.ValidateParams(ctx.UserContext(), request.Type, request.Domain, request.TokenId, request.BlockchainId); err != nil {
		return BadRequestWrapper(ctx, "blockchain", err)
	}

	// Lock the mutex after all validation checks
	operationMutex.Lock()
	defer operationMutex.Unlock()

	// Check if an operation is already running
	if isOperationRunning {
		return ctx.Status(fiber.StatusLocked).JSON(fiber.Map{"error": "Another operation is currently being executed. Please try again later."})
	}

	// Set the flag to indicate that an operation is running
	isOperationRunning = true

	// Create a channel to receive the result of the operation
	resultChan := make(chan types.ExecuteOperationResult)

	// Define the callback function
	callback := func() {
		// Reset the flag and unlock the mutex when the operation is done
		operationMutex.Lock()
		isOperationRunning = false
		operationMutex.Unlock()
	}

	// Execute the operation in a separate goroutine
	go func() {
		// Execute the operation and send the result to the channel
		operationId, err := o.Resources.OperationService.ExecuteOperation(ctx.UserContext(), request.Type, request.Domain, request.TokenId, request.BlockchainId, request.Amount, request.Operator, callback)
		resultChan <- types.ExecuteOperationResult{OperationId: operationId, Error: err}
		close(resultChan)
	}()

	// Wait for the result of the operation
	executeOpResult := <-resultChan
	if executeOpResult.Error != nil {
		return BadRequestWrapper(ctx, "operation", executeOpResult.Error)
	}

	return ctx.Status(fiber.StatusOK).JSON(&types.OperationResponse{Success: true, Message: fmt.Sprintf("operation %s was accepted to be processed on blockchain", executeOpResult.OperationId)})
}
