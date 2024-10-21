package handlers

import (
	types "crypto-braza-tokens-api/api/handlers/types"
	cfg "crypto-braza-tokens-api/configs"
	l "crypto-braza-tokens-api/utils/logger"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type WalletsHandler struct {
	Resources *cfg.Resources
}

// GetWallets retrieve the list of wallets
// @Summary Get the wallets list
// @Description retrieve the list of wallets
// @Tags Wallets
// @ID get-wallets
// @Produce json
// @Success 200 {array} wallet.Wallet
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/wallets [get]
func (w WalletsHandler) GetWallets(ctx *fiber.Ctx) error {
	result, err := w.Resources.WalletService.FindAll(ctx.UserContext())
	if err != nil {
		return InternalErrorWrapper(ctx, "wallet", err)
	}

	return ObjectResultWrapper(ctx, result)
}

// GetWallets retrieve a wallet by id
// @Summary Get a wallet
// @Description retrieve a wallet by id
// @Tags Wallets
// @ID get-wallet-by-id
// @Produce json
// @Success 200 {array} wallet.Wallet
// @Failure 400 {object} types.ErrorMessage
// @Failure 404 {object} types.ErrorMessage
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/wallets/{id} [get]
func (w WalletsHandler) GetWalletByID(ctx *fiber.Ctx) error {
	err := ValidatePathParam(ctx, "id")
	if err != nil {
		return BadRequestWrapper(ctx, "wallet", err)
	}

	result, err := w.Resources.WalletService.FindWalletById(ctx.UserContext(), ctx.Params("id"))
	if err != nil {
		return InternalErrorWrapper(ctx, "wallet", err)
	}

	return ObjectResultWrapper(ctx, result)
}

// GetWallets retrieve a wallet by address and blockchain
// @Summary Get a wallet
// @Description retrieve a wallet by address and blockchain
// @Tags Wallets
// @ID get-wallet-by-address-and-blockchain
// @Produce json
// @Success 200 {object} wallet.Wallet
// @Failure 400 {object} types.ErrorMessage
// @Failure 404 {object} types.ErrorMessage
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/wallets/address/{address}/blockchain_id/{blockchain_id} [get]
func (w WalletsHandler) GetWalletByAddressAndBlockchain(ctx *fiber.Ctx) error {
	ValidatePathParam(ctx, "address")
	ValidatePathParam(ctx, "blockchain_id")

	result, err := w.Resources.WalletService.FindWalletByAddressAndBlockchain(ctx.UserContext(), ctx.Params("address"), ctx.Params("blockchain_id"))
	if err != nil {
		return InternalErrorWrapper(ctx, "wallet", err)
	}

	return ObjectResultWrapper(ctx, result)
}

// GetWallets retrieve a wallet by blockchain, wallet type and domain
// @Summary Get a wallet
// @Description retrieve a wallet by blockchain, wallet type and domain
// @Tags Wallets
// @ID get-wallet-by-blockchain-wallet-type-and-domain
// @Produce json
// @Success 200 {object} wallet.Wallet
// @Failure 400 {object} types.ErrorMessage
// @Failure 404 {object} types.ErrorMessage
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/wallets/blockchain/{blockchain_id}/type/{wallet_type}/domain/{domain} [get]
func (w WalletsHandler) GetWalletByBlockchainWalletTypeAndDomain(ctx *fiber.Ctx) error {
	ValidatePathParam(ctx, "blockchain_id")
	ValidatePathParam(ctx, "wallet_type")
	ValidatePathParam(ctx, "domain")

	result, err := w.Resources.WalletService.FindWalletByBlockchainWalletTypeAndDomain(ctx.UserContext(), ctx.Params("blockchain_id"), ctx.Params("wallet_type"), ctx.Params("domain"))
	if err != nil {
		return InternalErrorWrapper(ctx, "wallet", err)
	}

	return ObjectResultWrapper(ctx, result)
}

// GetWallets retrieve the list of wallets by blockchain and domain
// @Summary Get the wallets list
// @Description retrieve the list of wallets by blockchain and domain
// @Tags Wallets
// @ID get-wallets-by-blockchain-and-domain
// @Produce json
// @Success 200 {array} wallet.Wallet
// @Failure 400 {object} types.ErrorMessage
// @Failure 404 {object} types.ErrorMessage
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/wallets/blockchain/{blockchain_id}/domain/{domain} [get]
func (w WalletsHandler) GetWalletsByBlockchainAndDomain(ctx *fiber.Ctx) error {
	ValidatePathParam(ctx, "blockchain_id")
	ValidatePathParam(ctx, "domain")

	result, err := w.Resources.WalletService.FindAllByBlockchainAndDomain(ctx.UserContext(), ctx.Params("blockchain_id"), ctx.Params("domain"))
	if err != nil {
		return InternalErrorWrapper(ctx, "wallet", err)
	}

	return ObjectResultWrapper(ctx, result)
}

// GetWalletsBalances retrieve the list of wallets balances
// @Summary Get the wallets balances list
// @Description retrieve the list of wallets balances
// @Tags Wallets
// @ID get-wallets-balances
// @Produce json
// @Success 200 {object} wallet.TokenBalance
// @Failure 400 {object} types.ErrorMessage
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/wallets/balances [get]
func (w WalletsHandler) GetWalletsBalances(ctx *fiber.Ctx) error {
	result, err := w.Resources.WalletService.GetAllBalances(ctx.UserContext())
	if err != nil {
		if err.Error() == "wallet balances are only available on MAINNET network" {
			l.Logger.Info("handler: no wallet balances found", zap.Error(err))
			return BadRequestWrapper(ctx, "wallet", err)
		}

		return InternalErrorWrapper(ctx, "wallet", err)
	}

	return ObjectResultWrapper(ctx, result)
}

// PostWallet create a new wallet
// @Summary Create a new wallet
// @Description create a new wallet
// @Tags Wallets
// @ID post-wallet
// @Produce json
// @Success 200 {object} wallet.Wallet
// @Failure 400 {object} types.ErrorMessage
// @Failure 409 {object} types.ErrorMessage
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/wallets [post]
func (w WalletsHandler) PostWallet(ctx *fiber.Ctx) error {
	request := types.SaveWalletRequest{}

	if err := request.FromBody(ctx); err != nil {
		return BadRequestWrapper(ctx, "wallet", err)
	}

	if err := request.IsValid(); err != nil {
		return BadRequestWrapper(ctx, "wallet", err)
	}

	result, err := w.Resources.WalletService.SaveWallet(ctx.UserContext(), request.Blockchain, request.Name, request.Address, request.Type, request.Domain, request.IsActive)
	if err != nil {
		return InternalErrorWrapper(ctx, "wallet", err)
	}

	return MessageResultWrapper(ctx, result.Hex())
}

// DeleteWallet delete a wallet
// @Summary Delete a wallet
// @Description delete a wallet
// @Tags Wallets
// @ID delete-wallet
// @Produce json
// @Param id path string true "Wallet ID"
// @Success 200 {object} wallet.Wallet
// @Failure 400 {object} types.ErrorMessage
// @Failure 404 {object} types.ErrorMessage
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/wallets/{id} [delete]
func (w WalletsHandler) DeleteWallet(ctx *fiber.Ctx) error {
	err := ValidatePathParam(ctx, "id")
	if err != nil {
		return BadRequestWrapper(ctx, "wallet", err)
	}

	result, err := w.Resources.WalletService.DeleteWallet(ctx.UserContext(), ctx.Params("id"))
	if err != nil {
		return InternalErrorWrapper(ctx, "wallet", err)
	}

	return MessageResultWrapper(ctx, result)
}

// PatchWallet update a wallet
// @Summary Update a wallet
// @Description update a wallet
// @Tags Wallets
// @ID patch-wallet
// @Accept json
// @Produce json
// @Param request body types.EditWalletRequest true "Wallet"
// @Success 200 {object} wallet.Wallet
// @Failure 400 {object} types.ErrorMessage
// @Failure 404 {object} types.ErrorMessage
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/wallets [patch]
func (w WalletsHandler) PatchWallet(ctx *fiber.Ctx) error {
	request := types.EditWalletRequest{}

	if err := request.FromBody(ctx); err != nil {
		return BadRequestWrapper(ctx, "wallet", err)
	}

	if err := request.IsValid(); err != nil {
		return BadRequestWrapper(ctx, "wallet", err)
	}

	result, err := w.Resources.WalletService.EditWallet(ctx.UserContext(), request.Blockchain, request.Name, request.Address, request.Type, request.Domain, request.IsActive)
	if err != nil {
		return InternalErrorWrapper(ctx, "wallet", err)
	}

	return MessageResultWrapper(ctx, result.Hex())
}
