package wallet

import (
	"time"
)

type Base struct {
	ID        string    `json:"id"`
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
	Domain       string      `json:"domain"`
	Blockchain   *Blockchain `json:"blockchain"`
}

type TokenBalance struct {
	Name     string `json:"name"`
	Address  string `json:"address"`
	Contract string `json:"contract"`
	Amount   string `json:"amount"`
}
