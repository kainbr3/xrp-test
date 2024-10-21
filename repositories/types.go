package repositories

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	m "go.mongodb.org/mongo-driver/mongo"
)

type Blockchain struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Abbr      string             `bson:"abbr" json:"abbr"`
	MainToken string             `bson:"main_token" json:"main_token"`
	IsActive  bool               `bson:"is_active" json:"is_active"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type Token struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	Blockchain string             `bson:"blockchain" json:"blockchain"`
	Name       string             `bson:"name" json:"name"`
	Abbr       string             `bson:"abbr" json:"abbr"`
	Contract   string             `bson:"contract" json:"contract"`
	Address    string             `bson:"address" json:"address"`
	Precision  int                `bson:"precision" json:"precision"`
	Type       string             `bson:"type" json:"type"`
	IsActive   bool               `bson:"is_active" json:"is_active"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

type Wallet struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	Blockchain string             `bson:"blockchain" json:"blockchain"`
	Name       string             `bson:"name" json:"name"`
	Address    string             `bson:"address" json:"address"`
	Type       string             `bson:"type" json:"type"`
	Domain     string             `bson:"domain" json:"domain"`
	IsActive   bool               `bson:"is_active" json:"is_active"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

type FireblocksAccount struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	WalletID  string             `bson:"wallet_id" json:"wallet_id"`
	VaultID   string             `bson:"vault_id" json:"vault_id"`
	AssetID   string             `bson:"asset_id" json:"asset_id"`
	Name      string             `bson:"name" json:"name"`
	Alias     string             `bson:"alias" json:"alias"`
	Domain    string             `bson:"domain" json:"domain"`
	Flags     int                `bson:"acc_flags" json:"acc_flags"`
	PublicKey string             `bson:"public_key" json:"public_key"`
	IsActive  bool               `bson:"is_active" json:"is_active"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type Operation struct {
	ID               primitive.ObjectID `bson:"_id" json:"id"`
	Type             string             `bson:"type" json:"type"`
	Domain           string             `bson:"domain" json:"domain"`
	Amount           string             `bson:"amount" json:"amount"`
	Operator         string             `bson:"operator" json:"operator"`
	FireblocksStatus string             `bson:"fireblocks_status" json:"fireblocks_status"`
	BlockchainStatus string             `bson:"blockchain_status" json:"blockchain_status"`
	FireblocksId     string             `bson:"fireblocks_id" json:"fireblocks_id"`
	TransactionHash  string             `bson:"transaction_hash" json:"transaction_hash"`
	TransactionLink  string             `bson:"transaction_link" json:"transaction_link"`
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updated_at"`
}

type OperationType struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `bson:"name" json:"name"`
	IsActive  bool               `bson:"is_active" json:"is_active"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type OperationDomain struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `bson:"name" json:"name"`
	IsActive  bool               `bson:"is_active" json:"is_active"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type OperationLog struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	Event        string             `bson:"event" json:"name"`
	Description  string             `bson:"description" json:"description"`
	OperationID  string             `bson:"operation_id" json:"operation_id"`
	FireblocksID string             `bson:"fireblocks_id" json:"fireblocks_id"`
	Payload      any                `bson:"payload" json:"payload"`
	Response     any                `bson:"response" json:"response"`
	Error        any                `bson:"error" json:"error"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}

type QueryParams struct {
	FilterParam string `json:"filter_param"`
	FilterValue string `json:"filter_value"`
	SortField   string `json:"sort_field"`
	SortOrder   string `json:"sort_order"`
	Page        int    `json:"page"`
	Limit       int    `json:"limit"`
}

type PaginatedResult struct {
	TotalCount   int `json:"total_count"`
	TotalPages   int `json:"total_pages"`
	CurrentPage  int `json:"current_page"`
	NextPage     int `json:"next_page"`
	PreviousPage int `json:"previous_page"`
	Data         any `json:"data"`
}

type PaginatedOperations struct {
	TotalCount   int          `json:"total_count"`
	TotalPages   int          `json:"total_pages"`
	CurrentPage  int          `json:"current_page"`
	NextPage     int          `json:"next_page"`
	PreviousPage int          `json:"previous_page"`
	Data         []*Operation `json:"data"`
}

type Transaction struct {
	ID              primitive.ObjectID `bson:"_id" json:"id"`
	Type            string             `bson:"type" json:"type"`
	Domain          string             `bson:"domain" json:"domain"`
	Amount          string             `bson:"amount" json:"amount"`
	Operator        string             `bson:"operator" json:"operator"`
	Status          string             `bson:"status" json:"status"`
	ExternalId      string             `bson:"external_id" json:"external_id"`
	FireblocksId    string             `bson:"fireblocks_id" json:"fireblocks_id"`
	TransactionHash string             `bson:"transaction_hash" json:"transaction_hash"`
	TransactionLink string             `bson:"transaction_link" json:"transaction_link"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
}

var _ IDocument = (*TransactionType)(nil)

type TransactionType struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `bson:"name" json:"name"`
	IsActive  bool               `bson:"is_active" json:"is_active"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

func (t *TransactionType) GetCollection(repo *Repository) *m.Collection {
	return repo.database.Collection("transactions_types")
}

func (t *TransactionType) GetName() string {
	return "transactions types"
}
