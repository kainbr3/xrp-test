package types

import (
	"crypto-braza-tokens-api/utils/validations"

	"github.com/gofiber/fiber/v2"
)

type SaveFireblocksAccountRequest struct {
	VaultID   string `json:"vault_id" example:"12" validate:"required"`
	AssetID   string `json:"asset_id" example:"XRP" validate:"required"`
	WalletID  string `json:"wallet_id" example:"66fc49c62dd62529e6e879bb" validate:"required"`
	Name      string `json:"name" example:"Name" validate:"required"`
	Alias     string `json:"alias" example:"Alias" validate:"required"`
	Domain    string `json:"domain" example:"Domain" validate:"required,oneof=GET-BRAZA BRAZA-ON BRAZA-DESK"`
	AccFlags  int    `json:"acc_flags" example:"2" validate:"required"`
	PublicKey string `json:"public_key" example:"032eb952987ef159445955c87461b9cdd498d0396f8ec5dc52af328babc2ed518f" validate:"required"`
	IsActive  bool   `json:"is_active" example:"true"`
}

func (t *SaveFireblocksAccountRequest) IsValid() error {
	return validations.Validate(t)
}

func (t *SaveFireblocksAccountRequest) FromBody(ctx *fiber.Ctx) error {
	return ctx.BodyParser(t)
}

type EditFireblocksAccountRequest struct {
	VaultID   string `json:"vault_id" example:"17"`
	AssetID   string `json:"asset_id" example:"XRP_TEST"`
	WalletID  string `json:"wallet_id" example:"18"`
	Name      string `json:"name" example:"Name"`
	Alias     string `json:"alias" example:"Alias"`
	Domain    string `json:"domain" example:"Domain" validate:"oneof=GET-BRAZA BRAZA-ON BRAZA-DESK"`
	AccFlags  int    `json:"acc_flags" example:"1"`
	PublicKey string `json:"public_key" example:"123abc"`
	IsActive  bool   `json:"is_active" example:"true"`
}

func (t *EditFireblocksAccountRequest) IsValid() error {
	return validations.Validate(t)
}

func (t *EditFireblocksAccountRequest) FromBody(ctx *fiber.Ctx) error {
	return ctx.BodyParser(t)
}
