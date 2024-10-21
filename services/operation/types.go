package operation

import (
	r "crypto-braza-tokens-api/repositories"
	"time"
)

type OperationWithLogs struct {
	r.Operation
	Logs []*r.OperationLog `json:"logs"`
}

type Base struct {
	ID        string    `json:"id"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OperationType struct {
	Base
	Name string `json:"name"`
}

type OperationDomain struct {
	Base
	Name string `json:"name"`
}
