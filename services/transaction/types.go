package transaction

import (
	"time"
)

type TransactionType struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type InternalTransactionParams struct {
	Domain          string `json:"domain"`
	BlockchinID     string `json:"blockchain_id"`
	WalletID        string `json:"wallet_id"`
	FbAccVaultID    string `json:"fb_acc_vault_id"`
	FbAccAssetID    string `json:"fb_acc_asset_id"`
	FbAccName       string `json:"fb_acc_name"`
	FbAccAlias      string `json:"fb_acc_alias"`
	Flags           int    `json:"wallet_acc_flags"`
	WalletPublicKey string `json:"wallet_public_key"`
	WalletName      string `json:"wallet_name"`
	Address         string `json:"wallet_address"`
	WalletType      string `json:"wallet_type"`
}

type Transaction struct {
}
