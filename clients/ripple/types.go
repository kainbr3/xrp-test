package ripple

type XrpJsonRpcRequest struct {
	Method string `json:"method"`
	Params []any  `json:"params"`
}

type XrpAccountInfo struct {
	Result *XrpAccountResult `json:"result"`
}

type XrpAccountResult struct {
	AccountData        *XrpAccountData      `json:"account_data"`
	AccountFlags       *XrpAccountFlags     `json:"account_flags"`
	LedgerCurrentIndex int                  `json:"ledger_current_index"`
	QueueData          *XrpAccountQueueData `json:"queue_data"`
	Status             string               `json:"status"`
	Validated          bool                 `json:"validated"`
}

type XrpAccountData struct {
	Account           string `json:"Account"`
	Balance           string `json:"Balance"`
	Domain            string `json:"Domain"`
	Flags             int    `json:"Flags"`
	LedgerEntryType   string `json:"LedgerEntryType"`
	OwnerCount        int    `json:"OwnerCount"`
	PreviousTxnID     string `json:"PreviousTxnID"`
	PreviousTxnLgrSeq int    `json:"PreviousTxnLgrSeq"`
	Sequence          int    `json:"Sequence"`
	TickSize          int    `json:"TickSize"`
	Index             string `json:"index"`
}

type XrpAccountFlags struct {
	AllowTrustLineClawback       bool `json:"allowTrustLineClawback"`
	DefaultRipple                bool `json:"defaultRipple"`
	DepositAuth                  bool `json:"depositAuth"`
	DisableMasterKey             bool `json:"disableMasterKey"`
	DisallowIncomingCheck        bool `json:"disallowIncomingCheck"`
	DisallowIncomingNFTokenOffer bool `json:"disallowIncomingNFTokenOffer"`
	DisallowIncomingPayChan      bool `json:"disallowIncomingPayChan"`
	DisallowIncomingTrustline    bool `json:"disallowIncomingTrustline"`
	DisallowIncomingXRP          bool `json:"disallowIncomingXRP"`
	GlobalFreeze                 bool `json:"globalFreeze"`
	NoFreeze                     bool `json:"noFreeze"`
	PasswordSpent                bool `json:"passwordSpent"`
	RequireAuthorization         bool `json:"requireAuthorization"`
	RequireDestinationTag        bool `json:"requireDestinationTag"`
}

type XrpAccountQueueData struct {
	TxnCount int `json:"txn_count"`
}

type XrpPaymentTx struct {
	TransactionType    string              `json:"TransactionType"`
	Account            string              `json:"Account"`
	Destination        string              `json:"Destination"`
	DestinationTag     uint32              `json:"DestinationTag"`
	Amount             *XrpPaymentTxAmount `json:"Amount"`
	Flags              uint32              `json:"Flags"`
	Sequence           uint32              `json:"Sequence"`
	Fee                string              `json:"Fee"`
	LastLedgerSequence uint32              `json:"LastLedgerSequence"`
	SigningPubKey      string              `json:"SigningPubKey"`
}

type XrpPaymentTxAmount struct {
	Currency string `json:"currency"`
	Issuer   string `json:"issuer"`
	Value    string `json:"value"`
}

type SubmitTxResultTxAmount struct {
	Currency string `json:"currency"`
	Issuer   string `json:"issuer"`
	Value    string `json:"value"`
}

type SubmitTxResultTxJSON struct {
	Account            string                 `json:"Account"`
	Amount             SubmitTxResultTxAmount `json:"Amount"`
	Destination        string                 `json:"Destination"`
	DestinationTag     int                    `json:"DestinationTag"`
	Fee                string                 `json:"Fee"`
	Flags              int                    `json:"Flags"`
	LastLedgerSequence int                    `json:"LastLedgerSequence"`
	Sequence           int                    `json:"Sequence"`
	SigningPubKey      string                 `json:"SigningPubKey"`
	TransactionType    string                 `json:"TransactionType"`
	TxnSignature       string                 `json:"TxnSignature"`
	Hash               string                 `json:"hash"`
}

type SubmitTxResult struct {
	Accepted                 bool                 `json:"accepted"`
	AccountSequenceAvailable int                  `json:"account_sequence_available"`
	AccountSequenceNext      int                  `json:"account_sequence_next"`
	Applied                  bool                 `json:"applied"`
	Broadcast                bool                 `json:"broadcast"`
	EngineResult             string               `json:"engine_result"`
	EngineResultCode         int                  `json:"engine_result_code"`
	EngineResultMessage      string               `json:"engine_result_message"`
	Kept                     bool                 `json:"kept"`
	OpenLedgerCost           string               `json:"open_ledger_cost"`
	Queued                   bool                 `json:"queued"`
	Status                   string               `json:"status"`
	TxBlob                   string               `json:"tx_blob"`
	TxJSON                   SubmitTxResultTxJSON `json:"tx_json"`
	ValidatedLedgerIndex     int                  `json:"validated_ledger_index"`
}

type SubmitTxResultResponse struct {
	Result *SubmitTxResult `json:"result"`
}

type Line struct {
	Account      string `json:"account"`
	Balance      string `json:"balance"`
	Currency     string `json:"currency"`
	Limit        string `json:"limit"`
	LimitPeer    string `json:"limit_peer"`
	NoRipple     bool   `json:"no_ripple"`
	NoRipplePeer bool   `json:"no_ripple_peer"`
	QualityIn    int    `json:"quality_in"`
	QualityOut   int    `json:"quality_out"`
}

type Result struct {
	Account            string `json:"account"`
	LedgerCurrentIndex int    `json:"ledger_current_index"`
	Lines              []Line `json:"lines"`
	Status             string `json:"status"`
	Validated          bool   `json:"validated"`
}

type AccountLinesResponse struct {
	Result `json:"result"`
}
