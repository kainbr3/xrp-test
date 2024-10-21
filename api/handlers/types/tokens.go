package types

import (
	"crypto-braza-tokens-api/utils/validations"

	"github.com/gofiber/fiber/v2"
)

type SaveTokenRequest struct {
	Name       string `json:"name" example:"TOKEN-NAME"`
	Abbr       string `json:"abbr" example:"BTC" validate:"required,min=2"`
	Contract   string `json:"contract" example:"0x1234567890" validate:"required"`
	Precision  int    `json:"precision" example:"18" validate:"required"`
	Type       string `json:"type" example:"NATIVE" validate:"required"`
	Blockchain string `json:"blockchain" example:"66f6fe7eccc6398d39e981f9" validate:"required"`
	IsActive   bool   `json:"is_active" example:"true"`
}

func (t *SaveTokenRequest) IsValid() error {
	return validations.Validate(t)
}

func (t *SaveTokenRequest) FromBody(ctx *fiber.Ctx) error {
	return ctx.BodyParser(t)
}

type EditTokenRequest struct {
	Name       string `json:"name" example:"TOKEN-NAME"`
	Abbr       string `json:"abbr" example:"BTC" validate:"min=2"`
	Contract   string `json:"contract" example:"0x1234567890"`
	Precision  int    `json:"precision" example:"18"`
	Type       string `json:"type" example:"NATIVE"`
	Blockchain string `json:"blockchain" example:"66f6fe7eccc6398d39e981f9"`
	IsActive   bool   `json:"is_active" example:"true"`
}

func (t *EditTokenRequest) IsValid() error {
	return validations.Validate(t)
}

func (t *EditTokenRequest) FromBody(ctx *fiber.Ctx) error {
	return ctx.BodyParser(t)
}
