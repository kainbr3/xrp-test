package types

import (
	"crypto-braza-tokens-api/utils/validations"

	"github.com/gofiber/fiber/v2"
)

type SaveTransactionTypeRequest struct {
	Name string `json:"name" example:"TYPE-NAME" validate:"required"`
}

func (s *SaveTransactionTypeRequest) IsValid() error {
	return validations.Validate(s)
}

func (s *SaveTransactionTypeRequest) FromBody(ctx *fiber.Ctx) error {
	return ctx.BodyParser(s)
}
