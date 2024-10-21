package xrpscan

import "time"

type TokenObligationsResponse struct {
	Currency string `json:"currency"`
	Value    string `json:"value"`
}

type Settings struct {
	RequireAuthorization bool   `json:"requireAuthorization"`
	Domain               string `json:"domain"`
}

type ParsedFlags struct {
	LsfRequireAuth bool `json:"lsfRequireAuth"`
}

type AccountInfo struct {
	Sequence                                  int         `json:"sequence"`
	XrpBalance                                string      `json:"xrpBalance"`
	OwnerCount                                int         `json:"ownerCount"`
	PreviousAffectingTransactionID            string      `json:"previousAffectingTransactionID"`
	PreviousAffectingTransactionLedgerVersion int         `json:"previousAffectingTransactionLedgerVersion"`
	Account                                   string      `json:"Account"`
	Balance                                   string      `json:"Balance"`
	Domain                                    string      `json:"Domain"`
	Flags                                     int         `json:"Flags"`
	LedgerEntryType                           string      `json:"LedgerEntryType"`
	OwnerCount0                               int         `json:"OwnerCount"`
	PreviousTxnID                             string      `json:"PreviousTxnID"`
	PreviousTxnLgrSeq                         int         `json:"PreviousTxnLgrSeq"`
	Sequence0                                 int         `json:"Sequence"`
	Index                                     string      `json:"index"`
	Settings                                  Settings    `json:"settings"`
	ParsedFlags                               ParsedFlags `json:"ParsedFlags"`
	Account0                                  string      `json:"account"`
	Parent                                    string      `json:"parent"`
	InitialBalance                            int         `json:"initial_balance"`
	Inception                                 time.Time   `json:"inception"`
	LedgerIndex                               int         `json:"ledger_index"`
	TxHash                                    string      `json:"tx_hash"`
	AccountName                               any         `json:"accountName"`
	ParentName                                any         `json:"parentName"`
	Advisory                                  any         `json:"advisory"`
}

type AccountInfoTokenBalance struct {
	Counterparty string `json:"counterparty"`
	Currency     string `json:"currency"`
	Value        string `json:"value"`
}
