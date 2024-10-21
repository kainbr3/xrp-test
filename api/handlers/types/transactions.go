package types

import (
	"crypto-braza-tokens-api/utils/validations"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type InternalTransferRequest struct {
	Type             string `json:"type" example:"ON-RAMP,OFF-RAMP" validate:"required,oneof=ON-RAMP OFF-RAMP"`
	SourceVault      string `json:"source_vault_id" example:"17" validate:"required"`
	DestinationVault string `json:"destination_vault_id" example:"18" validate:"required"`
	AssetId          string `json:"asset_id" example:"66f74acbba6b56108cb3e80a" validate:"required"`
	Amount           string `json:"amount" example:"2.75" validate:"required"`
}

func (i *InternalTransferRequest) IsValid() error {
	amountFloat, err := strconv.ParseFloat(i.Amount, 64)
	if err != nil {
		return fmt.Errorf("invalid amount: %v", err)
	}

	if amountFloat < 1 {
		return fmt.Errorf("the amount must be at least 1")
	}

	return validations.Validate(i)
}

// FromBody parses the request body into the OperationRequest struct
func (i *InternalTransferRequest) FromBody(ctx *fiber.Ctx) error {
	return ctx.BodyParser(i)
}
