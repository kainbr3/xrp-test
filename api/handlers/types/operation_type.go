package types

import (
	"crypto-braza-tokens-api/utils/validations"

	"github.com/gofiber/fiber/v2"
)

type SaveOperationTypeRequest struct {
	Name string `json:"name" example:"TYPE-NAME" validate:"required"`
}

func (t *SaveOperationTypeRequest) IsValid() error {
	return validations.Validate(t)
}

func (t *SaveOperationTypeRequest) FromBody(ctx *fiber.Ctx) error {
	return ctx.BodyParser(t)
}
