package handlers

import (
	types "crypto-braza-tokens-api/api/handlers/types"
	cfg "crypto-braza-tokens-api/configs"

	"strings"

	"github.com/gofiber/fiber/v2"
)

type TokensHandler struct {
	Resources *cfg.Resources
}

// GetTokens retrieve the list of supported tokens
// @Summary Get the tokens list
// @Description retrieve the list of supported tokens
// @Tags Tokens
// @ID get-tokens
// @Produce json
// @Success 200 {array} tokens.Token
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/tokens [get]
func (t TokensHandler) GetTokens(ctx *fiber.Ctx) error {
	result, err := t.Resources.TokenService.FindAll(ctx.UserContext())
	if err != nil {
		return InternalErrorWrapper(ctx, "token", err)
	}

	return ObjectResultWrapper(ctx, result)
}

// GetTokens retrieve a supported token by id
// @Summary Get a token
// @Description retrieve a supported token by id
// @Tags Tokens
// @ID get-token-by-id
// @Produce json
// @Success 200 {object} tokens.Token
// @Failure 404 {object} types.ErrorMessage
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/tokens/{id} [get]
func (t TokensHandler) GetTokenByID(ctx *fiber.Ctx) error {
	err := ValidatePathParam(ctx, "id")
	if err != nil {
		return BadRequestWrapper(ctx, "token", err)
	}

	result, err := t.Resources.TokenService.FindTokenByID(ctx.UserContext(), ctx.Params("id"))
	if err != nil {
		return InternalErrorWrapper(ctx, "token", err)
	}

	return ObjectResultWrapper(ctx, result)
}

// PostToken create a new token
// @Summary Create a new token
// @Description create a new token
// @Tags Tokens
// @ID post-token
// @Accept json
// @Produce json
// @Param token body types.SaveTokenRequest true "Token object that needs to be added"
// @Success 200 {object} types.Result
// @Failure 400 {object} types.ErrorMessage
// @Failure 409 {object} types.ErrorMessage
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/tokens [post]
func (t TokensHandler) PostToken(ctx *fiber.Ctx) error {
	request := types.SaveTokenRequest{}

	if err := request.FromBody(ctx); err != nil {
		return BadRequestWrapper(ctx, "token", err)
	}

	if err := request.IsValid(); err != nil {
		return BadRequestWrapper(ctx, "token", err)
	}

	abbr := strings.ToUpper(request.Abbr)

	result, err := t.Resources.TokenService.SaveToken(ctx.UserContext(), request.Blockchain, request.Name, abbr, request.Contract, request.Type, request.Precision, request.IsActive)
	if err != nil {
		return InternalErrorWrapper(ctx, "token", err)
	}

	return MessageResultWrapper(ctx, result.Hex())
}

// DeleteToken delete a token
// @Summary Delete a token
// @Description delete a token
// @Tags Tokens
// @ID delete-token
// @Param id path string true "Token ID"
// @Produce json
// @Success 200 {object} types.Result
// @Failure 404 {object} types.ErrorMessage
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/tokens/{id} [delete]
func (t TokensHandler) DeleteToken(ctx *fiber.Ctx) error {
	result, err := t.Resources.TokenService.DeleteToken(ctx.UserContext(), ctx.Params("id"))
	if err != nil {
		return InternalErrorWrapper(ctx, "token", err)
	}

	return MessageResultWrapper(ctx, result)
}

// PatchToken update a token
// @Summary Update a token
// @Description update a token
// @Tags Tokens
// @ID patch-token
// @Accept json
// @Produce json
// @Param request body types.EditTokenRequest true "Token"
// @Success 200 {object} types.Result
// @Failure 400 {object} types.ErrorMessage
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/tokens [patch]
func (t TokensHandler) PatchToken(ctx *fiber.Ctx) error {
	request := types.EditTokenRequest{}

	if err := request.FromBody(ctx); err != nil {
		return BadRequestWrapper(ctx, "token", err)
	}

	if err := request.IsValid(); err != nil {
		return BadRequestWrapper(ctx, "token", err)
	}

	result, err := t.Resources.TokenService.EditToken(ctx.UserContext(), request.Blockchain, request.Name, request.Abbr, request.Contract, request.Type, request.Precision, request.IsActive)
	if err != nil {
		return InternalErrorWrapper(ctx, "token", err)
	}

	return MessageResultWrapper(ctx, result.Hex())
}
