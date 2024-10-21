package types

import (
	"crypto-braza-tokens-api/utils/validations"

	"github.com/gofiber/fiber/v2"
)

type SaveWalletRequest struct {
	Name       string `json:"name" example:"WALLET-NAME" validate:"required"`
	Address    string `json:"address" example:"rfWmf1YZLfcaHVZioBBSUuRLHgMMSfBkBd" validate:"required"`
	Blockchain string `json:"blockchain" example:"66f6fe7eccc6398d39e981f9" validate:"required"`
	Type       string `json:"type" example:"NATIVE" validate:"required"`
	Domain     string `json:"domain" example:"DOMAIN-NAME" validate:"required"`
	IsActive   bool   `json:"is_active" example:"true"`
}

func (t *SaveWalletRequest) IsValid() error {
	return validations.Validate(t)
}

func (t *SaveWalletRequest) FromBody(ctx *fiber.Ctx) error {
	return ctx.BodyParser(t)
}

type EditWalletRequest struct {
	Name       string `json:"name" example:"WALLET-NAME"`
	Address    string `json:"address" example:"rfWmf1YZLfcaHVZioBBSUuRLHgMMSfBkBd"`
	Blockchain string `json:"blockchain" example:"66f6fe7eccc6398d39e981f9"`
	Type       string `json:"type" example:"NATIVE"`
	Domain     string `json:"domain" example:"DOMAIN-NAME"`
	IsActive   bool   `json:"is_active" example:"true"`
}

func (t *EditWalletRequest) IsValid() error {
	return validations.Validate(t)
}

func (t *EditWalletRequest) FromBody(ctx *fiber.Ctx) error {
	return ctx.BodyParser(t)
}
