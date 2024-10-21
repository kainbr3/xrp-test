package types

import (
	"crypto-braza-tokens-api/utils/validations"

	"github.com/gofiber/fiber/v2"
)

type SaveBlockchainRequest struct {
	Name      string `json:"name" example:"BLOCKHAIN-NAME" validate:"required"`
	Abbr      string `json:"abbr" example:"BTC" validate:"required,min=2"`
	MainToken string `json:"main_token" example:"BTC" validate:"required,min=2"`
	IsActive  bool   `json:"is_active" example:"true"`
}

func (b *SaveBlockchainRequest) IsValid() error {
	return validations.Validate(b)
}

func (b *SaveBlockchainRequest) FromBody(ctx *fiber.Ctx) error {
	return ctx.BodyParser(b)
}

type EditBlockchainRequest struct {
	Name      string `json:"name" example:"BLOCKHAIN-NAME"`
	Abbr      string `json:"abbr" example:"BTC" validate:"min=2"`
	MainToken string `json:"main_token" example:"BTC" validate:"min=2"`
	IsActive  bool   `json:"is_active" example:"true"`
}

func (b *EditBlockchainRequest) IsValid() error {
	return validations.Validate(b)
}

func (b *EditBlockchainRequest) FromBody(ctx *fiber.Ctx) error {
	return ctx.BodyParser(b)
}
