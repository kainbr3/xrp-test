package fireblocks

// ACCOUNT INFO MODELS:
type GetPbkInfoByVaultAccResponse struct {
	Status         int    `json:"status"`
	Algorithm      string `json:"algorithm"`
	DerivationPath []int  `json:"derivationPath"`
	PublicKey      string `json:"publicKey"`
}

// VAULT ACCOUNT MODELS:
type VaultAccountResponse struct {
	Id            string                `json:"id"`
	Name          string                `json:"name"`
	HiddenOnUI    bool                  `json:"hiddenOnUI"`
	CustomerRefId string                `json:"customerRefId"`
	AutoFuel      bool                  `json:"autoFuel"`
	Assets        []*VaultAssetResponse `json:"assets"`
}

type VaultAssetResponse struct {
	Id                   string `json:"id"`
	Total                string `json:"total"`
	Balance              string `json:"balance"`
	Available            string `json:"available"`
	Pending              string `json:"pending"`
	LockedAmount         string `json:"lockedAmount"`
	TotalStackedCPU      string `json:"totalStackedCPU"`
	TotalStackedNetwork  string `json:"totalStackedNetwork"`
	SelfStackedCPU       string `json:"selfStackedCPU"`
	SelfStakedNetwork    string `json:"selfStakedNetwork"`
	PendingRefundCPU     string `json:"pendingRefundCPU"`
	PendingRefundNetwork string `json:"pendingRefundNetwork"`
	Frozen               string `json:"frozen"`
	Staked               string `json:"staked"`
}

// VAULT ACCOUNT ASSET ADDRESS MODELS:
type VaultAccountAssetAddressResponse struct {
	AssetId       string `json:"assetId"`        // The ID of the asset
	Address       string `json:"address"`        // Address of the asset in a Vault Account, for BTC/LTC the address is in Segwit (Bech32) format, for BCH cash format
	LegacyAddress string `json:"legacyAddress"`  // For BTC/LTC/BCH the legacy format address
	Description   string `json:"description"`    // Description of the address
	Tag           string `json:"tag"`            // Destination tag for XRP, used as memo for EOS/XLM, for Signet/SEN it is the Bank Transfer Description
	Type          string `json:"type"`           // Address type
	CustomerRefId string `json:" customerRefId"` // [optional] The ID for AML providers to associate the owner of funds with transactions

}

// BASE TRANSACTION MODELS:
type TargetVaultAccount struct {
	Type    string `json:"type"`
	ID      string `json:"id"`
	Name    string `json:"name,omitempty"`
	SubType string `json:"subType,omitempty"`
}

// CREATE RAW TRANSACTION MODELS:
type RawTxExtraParameters struct {
	RawMessageData *RawTxMessageData `json:"rawMessageData"`
}

type RawTxMessageData struct {
	Messages []*RawTxContent `json:"messages"`
}

type RawTxContent struct {
	Content string `json:"content"`
}

type RawTransactionRequest struct {
	AssetID         string                `json:"assetId"`
	ExtraParameters *RawTxExtraParameters `json:"extraParameters"`
	Note            string                `json:"note"`
	Operation       string                `json:"operation"`
	Source          *TargetVaultAccount   `json:"source"`
}

// TRANSACTION DETAILS MODELS:
type TxByIdSignature struct {
	R       string `json:"r"`
	S       string `json:"s"`
	V       int    `json:"v"`
	FullSig string `json:"fullSig"`
}

type TxByIdSignedMessages struct {
	DerivationPath []int            `json:"derivationPath"`
	Algorithm      string           `json:"algorithm"`
	PublicKey      string           `json:"publicKey"`
	Signature      *TxByIdSignature `json:"signature"`
	Content        string           `json:"content"`
}

type TxByIdMessages struct {
	Content string `json:"content"`
}

type TxByIdRawMessageData struct {
	Messages []*TxByIdMessages `json:"messages"`
}

type TxByIdExtraParameters struct {
	RawMessageData *TxByIdRawMessageData `json:"rawMessageData"`
}

type TransactionByIdResponse struct {
	ID                            string                  `json:"id"`
	AssetID                       string                  `json:"assetId"`
	Source                        *TargetVaultAccount     `json:"source"`
	Destination                   *TargetVaultAccount     `json:"destination"`
	RequestedAmount               any                     `json:"requestedAmount"`
	Amount                        any                     `json:"amount"`
	NetAmount                     int                     `json:"netAmount"`
	AmountUSD                     any                     `json:"amountUSD"`
	Fee                           int                     `json:"fee"`
	NetworkFee                    int                     `json:"networkFee"`
	CreatedAt                     int64                   `json:"createdAt"`
	LastUpdated                   int64                   `json:"lastUpdated"`
	Status                        string                  `json:"status"`
	TxHash                        string                  `json:"txHash"`
	SubStatus                     string                  `json:"subStatus"`
	SourceAddress                 string                  `json:"sourceAddress"`
	DestinationAddress            string                  `json:"destinationAddress"`
	DestinationAddressDescription string                  `json:"destinationAddressDescription"`
	DestinationTag                string                  `json:"destinationTag"`
	SignedBy                      []any                   `json:"signedBy"`
	CreatedBy                     string                  `json:"createdBy"`
	RejectedBy                    string                  `json:"rejectedBy"`
	AddressType                   string                  `json:"addressType"`
	Note                          string                  `json:"note"`
	ExchangeTxID                  string                  `json:"exchangeTxId"`
	FeeCurrency                   string                  `json:"feeCurrency"`
	Operation                     string                  `json:"operation"`
	AmountInfo                    any                     `json:"amountInfo"`
	FeeInfo                       any                     `json:"feeInfo"`
	SignedMessages                []*TxByIdSignedMessages `json:"signedMessages"`
	ExtraParameters               *TxByIdExtraParameters  `json:"extraParameters"`
	Destinations                  []any                   `json:"destinations"`
	BlockInfo                     any                     `json:"blockInfo"`
	AssetType                     string                  `json:"assetType"`
}

// TRANSACTIONS RESPONSE MODELS:
type SubmittedTransactionResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// CREATE INTERNAL TRANSACTION MODELS:
type InternalTransactionRequest struct {
	Operation    string              `json:"operation"`
	AssetID      string              `json:"assetId"`
	Amount       string              `json:"amount"`
	Source       *TargetVaultAccount `json:"source"`
	Destination  *TargetVaultAccount `json:"destination"`
	Note         string              `json:"note"`
	ExternalTxID string              `json:"externalTxId"`
}
