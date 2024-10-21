package blockchains

import "time"

type Blockchain struct {
	ID        string    `json:"id" bson:"_id"`
	Name      string    `json:"name"`
	Abbr      string    `json:"abbr"`
	MainToken string    `json:"main_token"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
