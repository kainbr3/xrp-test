package types

import (
	"crypto-braza-tokens-api/utils/validations"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type OperationResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"operation was accepted to be processed on blockchain"`
}

type OperationRequest struct {
	Type         string `json:"type" example:"MINT,BURN" validate:"required,oneof=MINT BURN"`
	BlockchainId string `json:"blockchain_id" example:"66f6fe7eccc6398d39e981f9" validate:"required"`
	TokenId      string `json:"token_id" example:"66f74acbba6b56108cb3e80a" validate:"required"`
	Amount       string `json:"amount" example:"2.75" validate:"required"`
	Domain       string `json:"domain" example:"GET-BRAZA" validate:"required,oneof=GET-BRAZA BRAZA-ON BRAZA-DESK"`
	Operator     string `json:"operator" example:"123e4567-e89b-12d3-a456-426614174000" validate:"required"`
}

// IsValid validates the OperationRequest fields
func (o *OperationRequest) IsValid() error {
	amountFloat, err := strconv.ParseFloat(o.Amount, 64)
	if err != nil {
		return fmt.Errorf("invalid amount: %v", err)
	}

	if amountFloat < 1 {
		return fmt.Errorf("the amount must be at least 1")
	}

	return validations.Validate(o)
}

// FromBody parses the request body into the OperationRequest struct
func (o *OperationRequest) FromBody(ctx *fiber.Ctx) error {
	return ctx.BodyParser(o)
}
