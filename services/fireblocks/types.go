package fireblocks

import (
	"time"

	"github.com/shopspring/decimal"
)

type Base struct {
	ID        string    `json:"id" bson:"_id"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Blockchain struct {
	Base
	Name      string `json:"name"`
	Abbr      string `json:"abbr"`
	MainToken string `json:"main_token"`
}

type Wallet struct {
	Base
	Name         string      `json:"name"`
	Address      string      `json:"address"`
	Type         string      `json:"type"`
	BlockchainID string      `json:"blockchain_id"`
	Blockchain   *Blockchain `json:"blockchain,omitempty"`
}

type FireblocksAccount struct {
	Base
	VaultID           string  `json:"vault_id"`
	AssetID           string  `json:"asset_id"`
	WalletID          string  `json:"wallet_id"`
	Name              string  `json:"name"`
	Alias             string  `json:"alias"`
	PublicKey         string  `json:"public_key"`
	PublicKeyFallback string  `json:"public_key_fallback"`
	Wallet            *Wallet `json:"wallet,omitempty"`
}

type Asset struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	Precision  int             `json:"precision"`
	BalanceStr string          `json:"balance_str"`
	Balance    decimal.Decimal `json:"balance"`
}

type FireblocksAccountAssets struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Assets []*Asset `json:"assets"`
}
