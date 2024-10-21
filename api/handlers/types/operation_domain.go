package types

import (
	"crypto-braza-tokens-api/utils/validations"

	"github.com/gofiber/fiber/v2"
)

type SaveOperationDomainRequest struct {
	Name     string `json:"name" example:"DOMAIN-NAME" validate:"required"`
	IsActive bool   `json:"is_active" example:"true"`
}

func (d *SaveOperationDomainRequest) IsValid() error {
	return validations.Validate(d)
}

func (d *SaveOperationDomainRequest) FromBody(ctx *fiber.Ctx) error {
	return ctx.BodyParser(d)
}

type EditOperationDomainRequest struct {
	Name     string `json:"name" example:"DOMAIN-NAME"`
	IsActive bool   `json:"is_active" example:"true"`
}

func (d *EditOperationDomainRequest) IsValid() error {
	return validations.Validate(d)
}

func (d *EditOperationDomainRequest) FromBody(ctx *fiber.Ctx) error {
	return ctx.BodyParser(d)
}
