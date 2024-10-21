package handlers

import (
	types "crypto-braza-tokens-api/api/handlers/types"
	cfg "crypto-braza-tokens-api/configs"
	l "crypto-braza-tokens-api/utils/logger"

	"strings"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type BlockchainHandler struct {
	Resources *cfg.Resources
}

// GetBlockchains retrieve the list of supported blockchains
// @Summary Get the blockchains list
// @Description retrieve the list of supported blockchains
// @Tags Blockchains
// @ID get-blockchains
// @Produce json
// @Success 200 {array} blockchains.Blockchain
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/blockchains [get]
func (b BlockchainHandler) GetBlockchains(ctx *fiber.Ctx) error {
	result, err := b.Resources.BlockchainService.FindAll(ctx.UserContext())
	if err != nil {
		l.Logger.Error("handler: error finding blockchains", zap.Error(err))
		return InternalErrorWrapper(ctx, "blockchain", err)
	}

	return ObjectResultWrapper(ctx, result)
}

// GetBlockchains retrieve a supported blockchain by id
// @Summary Get a blockchain
// @Description retrieve a supported blockchain by id
// @Tags Blockchains
// @ID get-blockchain-by-id
// @Produce json
// @Param id path string true "Blockchain ID"
// @Success 200 {array} blockchains.Blockchain
// @Failure 404 {object} types.ErrorMessage
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/blockchains/{id} [get]
func (b BlockchainHandler) GetBlockchainByID(ctx *fiber.Ctx) error {
	err := ValidatePathParam(ctx, "id")
	if err != nil {
		return BadRequestWrapper(ctx, "blockchain", err)
	}

	result, err := b.Resources.BlockchainService.FindBlockchainByID(ctx.UserContext(), ctx.Params("id"))
	if err != nil {
		l.Logger.Error("handler: error finding blockchain", zap.Error(err))
		return InternalErrorWrapper(ctx, "blockchain", err)
	}

	return ObjectResultWrapper(ctx, result)
}

// GetBlockchainTokens retrieve the list of tokens for a blockchain
// @Summary Get the blockchain tokens list
// @Description retrieve the list of tokens for a blockchain
// @Tags Blockchains
// @ID get-blockchains-tokens
// @Produce json
// @Param id path string true "Blockchain ID"
// @Success 200 {array} tokens.Token
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/blockchains/{id}/tokens [get]
func (b BlockchainHandler) GetBlockchainTokens(ctx *fiber.Ctx) error {
	result, err := b.Resources.TokenService.FindTokensByBlockchainID(ctx.UserContext(), ctx.Params("id"))
	if err != nil {
		l.Logger.Error("handler: error finding blockchain tokens", zap.Error(err))
		return InternalErrorWrapper(ctx, "blockchain", err)
	}

	return ObjectResultWrapper(ctx, result)
}

// PostBlockchain create a new blockchain
// @Summary Create a new blockchain
// @Description create a new blockchain
// @Tags Blockchains
// @ID post-blockchains
// @Accept json
// @Produce json
// @Param request body types.SaveBlockchainRequest true "Blockchain data"
// @Success 200 {object} types.Result
// @Failure 400 {object} types.ErrorMessage
// @Failure 409 {object} types.ErrorMessage
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/blockchains [post]
func (b BlockchainHandler) PostBlockchain(ctx *fiber.Ctx) error {
	request := types.SaveBlockchainRequest{}

	if err := request.FromBody(ctx); err != nil {
		return BadRequestWrapper(ctx, "blockchain", err)
	}

	if err := request.IsValid(); err != nil {
		return BadRequestWrapper(ctx, "blockchain", err)
	}

	abbr := strings.ToUpper(request.Abbr)
	mainToken := strings.ToUpper(request.MainToken)

	blockchain, err := b.Resources.BlockchainService.SaveBlockchain(ctx.UserContext(), request.Name, abbr, mainToken, request.IsActive)
	if err != nil {
		if err.Error() == "blockchain already exists" {

			return ConflictWrapper(ctx, err.Error())
		}

		return InternalErrorWrapper(ctx, "blockchain", err)
	}

	return MessageResultWrapper(ctx, blockchain.Hex())
}

// DeleteBlockchain delete a blockchain
// @Summary Delete a blockchain
// @Description delete a blockchain
// @Tags Blockchains
// @ID delete-blockchains
// @Param id path string true "Blockchain ID"
// @Produce json
// @Success 200 {object} types.Result
// @Failure 400 {object} types.ErrorMessage
// @Failure 404 {object} types.ErrorMessage
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/blockchains/{id} [delete]
func (b BlockchainHandler) DeleteBlockchain(ctx *fiber.Ctx) error {
	err := ValidatePathParam(ctx, "id")
	if err != nil {
		return BadRequestWrapper(ctx, "blockchain", err)
	}

	result, err := b.Resources.BlockchainService.DeleteBlockchain(ctx.UserContext(), ctx.Params("id"))
	if err != nil {
		return InternalErrorWrapper(ctx, "blockchain", err)
	}

	return MessageResultWrapper(ctx, result)
}

// PatchBlockchain update a blockchain
// @Summary Update a blockchain
// @Description update a blockchain
// @Tags Blockchains
// @ID patch-blockchains
// @Accept json
// @Produce json
// @Param request body types.EditBlockchainRequest true "Blockchain data"
// @Success 200 {object} types.Result
// @Failure 400 {object} types.ErrorMessage
// @Failure 404 {object} types.ErrorMessage
// @Failure 500 {object} types.ErrorMessage
// @Router /api/v1/blockchains [patch]
func (b BlockchainHandler) PatchBlockchain(ctx *fiber.Ctx) error {
	request := types.EditBlockchainRequest{}

	if err := request.FromBody(ctx); err != nil {
		return BadRequestWrapper(ctx, "blockchain", err)
	}

	if err := request.IsValid(); err != nil {
		return BadRequestWrapper(ctx, "blockchain", err)
	}

	abbr := strings.ToUpper(request.Abbr)
	mainToken := strings.ToUpper(request.MainToken)

	blockchain, err := b.Resources.BlockchainService.EditBlockchain(ctx.UserContext(), request.Name, abbr, mainToken, request.IsActive)
	if err != nil {
		return InternalErrorWrapper(ctx, "blockchain", err)
	}

	return MessageResultWrapper(ctx, blockchain.Hex())
}
