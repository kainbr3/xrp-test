package tokens

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

type Token struct {
	Base
	Name         string      `json:"name"`
	Abbr         string      `json:"abbr"`
	Contract     string      `json:"contract"`
	Precision    int         `json:"precision"`
	Type         string      `json:"type"`
	BlockchainID string      `json:"blockchain_id"`
	Blockchain   *Blockchain `json:"blockchain"`
}
