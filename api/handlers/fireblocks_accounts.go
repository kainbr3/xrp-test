package handlers

import (
	types "crypto-braza-tokens-api/api/handlers/types"
	cfg "crypto-braza-tokens-api/configs"

	"github.com/gofiber/fiber/v2"
)

type FireblocksAccountsHandler struct {
	Resources *cfg.Resources
}

// GetFireblocksAccounts retrieve the list of fireblocks accounts
// @Summary Get the fireblocks accounts list
// @Description retrieve the list of fireblocks accounts
// @Tags FireblocksAccounts
// @ID get-fireblocks-accounts
// @Produce json
// @Success 200 {array} fireblocks.FireblocksAccount
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/fireblocks-accounts [get]
func (f FireblocksAccountsHandler) GetFireblocksAccounts(ctx *fiber.Ctx) error {
	result, err := f.Resources.FireblocksService.FindAllFireblocksAccounts(ctx.UserContext())
	if err != nil {
		return InternalErrorWrapper(ctx, "fireblocks accounts", err)
	}

	return ObjectResultWrapper(ctx, result)
}

// GetFireblocksAccounts retrieve a fireblocks account by id
// @Summary Get a fireblocks account
// @Description retrieve a fireblocks account by id
// @Tags FireblocksAccounts
// @ID get-fireblocks-accounts-by-id
// @Produce json
// @Success 200 {object} fireblocks.FireblocksAccount
// @Failure 400 {object} types.ErrorMessage
// @Failure 404 {object} types.ErrorMessage
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/fireblocks-accounts/{id} [get]
func (f FireblocksAccountsHandler) GetFireblocksAccountById(ctx *fiber.Ctx) error {
	err := ValidatePathParam(ctx, "id")
	if err != nil {
		return BadRequestWrapper(ctx, "fireblocks accounts", err)
	}

	result, err := f.Resources.FireblocksService.FindFireblocksAccountByID(ctx.UserContext(), ctx.Params("id"))
	if err != nil {
		return InternalErrorWrapper(ctx, "fireblocks accounts", err)
	}

	return ObjectResultWrapper(ctx, result)
}

// GetFireblocksAccounts retrieve a fireblocks account by vault id
// @Summary Get a fireblocks account
// @Description retrieve a fireblocks account by vault id
// @Tags FireblocksAccounts
// @ID get-fireblocks-accounts-by-vault-id
// @Produce json
// @Success 200 {object} fireblocks.FireblocksAccount
// @Failure 400 {object} types.ErrorMessage
// @Failure 404 {object} types.ErrorMessage
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/fireblocks-accounts/vault/{vault_id} [get]
func (f FireblocksAccountsHandler) GetFireblocksAccountByVaultId(ctx *fiber.Ctx) error {
	ValidatePathParam(ctx, "vault_id")

	result, err := f.Resources.FireblocksService.FindFireblocksAccountByID(ctx.UserContext(), ctx.Params("vault_id"))
	if err != nil {
		return InternalErrorWrapper(ctx, "fireblocks accounts", err)
	}

	return ObjectResultWrapper(ctx, result)
}

// PostFireblocksAccounts create a new fireblocks account
// @Summary Create a new fireblocks account
// @Description create a new fireblocks account
// @Tags FireblocksAccounts
// @ID post-fireblocks-accounts
// @Accept json
// @Produce json
// @Param request body types.SaveFireblocksAccountRequest true "Fireblocks Account"
// @Success 200 {object} types.Result
// @Failure 400 {object} types.ErrorMessage
// @Failure 409 {object} types.ErrorMessage
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/fireblocks-accounts [post]
func (f FireblocksAccountsHandler) PostFireblocksAccounts(ctx *fiber.Ctx) error {
	request := types.SaveFireblocksAccountRequest{}

	if err := request.FromBody(ctx); err != nil {
		return BadRequestWrapper(ctx, "fireblocks accounts", err)
	}

	if err := request.IsValid(); err != nil {
		return BadRequestWrapper(ctx, "fireblocks accounts", err)
	}

	result, err := f.Resources.FireblocksService.SaveFireblocksAccount(ctx.UserContext(), request.VaultID, request.AssetID, request.WalletID, request.Name, request.Alias, request.Domain, request.PublicKey, request.AccFlags, request.IsActive)
	if err != nil {
		if err.Error() == "service: fireblocks account already exists" {
			return ConflictWrapper(ctx, err.Error())
		}

		return InternalErrorWrapper(ctx, "fireblocks accounts", err)
	}

	return MessageResultWrapper(ctx, result.Hex())
}

// DeleteFireblocksAccount delete a fireblocks account
// @Summary Delete a fireblocks account
// @Description delete a fireblocks account
// @Tags FireblocksAccounts
// @ID delete-fireblocks-accounts
// @Param id path string true "Fireblocks Account ID"
// @Produce json
// @Success 200 {object} types.Result
// @Failure 400 {object} types.ErrorMessage
// @Failure 404 {object} types.ErrorMessage
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/fireblocks-accounts/{id} [delete]
func (f FireblocksAccountsHandler) DeleteFireblocksAccount(ctx *fiber.Ctx) error {
	err := ValidatePathParam(ctx, "id")
	if err != nil {
		return BadRequestWrapper(ctx, "fireblocks accounts", err)
	}

	result, err := f.Resources.FireblocksService.DeleteFireblocks(ctx.UserContext(), ctx.Params("id"))
	if err != nil {
		return InternalErrorWrapper(ctx, "fireblocks accounts", err)
	}

	return MessageResultWrapper(ctx, result)
}

// PatchFireblocksAccounts update a fireblocks account
// @Summary Update a fireblocks account
// @Description update a fireblocks account
// @Tags FireblocksAccounts
// @ID patch-fireblocks-accounts
// @Accept json
// @Produce json
// @Param request body types.EditFireblocksAccountRequest true "Fireblocks Account"
// @Success 200 {object} types.Result
// @Failure 400 {object} types.ErrorMessage
// @Failure 404 {object} types.ErrorMessage
// @Failure 500 {object} types.ErrorMessage
func (f FireblocksAccountsHandler) PatchFireblocksAccounts(ctx *fiber.Ctx) error {
	request := types.EditFireblocksAccountRequest{}

	if err := request.FromBody(ctx); err != nil {
		return BadRequestWrapper(ctx, "fireblocks accounts", err)
	}

	if err := request.IsValid(); err != nil {
		return BadRequestWrapper(ctx, "fireblocks accounts", err)
	}

	result, err := f.Resources.FireblocksService.EditFireblocksAccount(ctx.UserContext(), request.VaultID, request.AssetID, request.WalletID, request.Name, request.Alias, request.Domain, request.PublicKey, request.AccFlags, request.IsActive)
	if err != nil {
		return InternalErrorWrapper(ctx, "fireblocks accounts", err)
	}

	return MessageResultWrapper(ctx, result.Hex())
}
