package operation

import (
	xrpn "crypto-braza-tokens-api/clients/ripple"
	"encoding/json"
)

func buildRippleRawTransactionPayload(
	walletFromAddress, walletToAddress,
	tokenAbbr, issuerAddress,
	amount, publicKey string,
	flags, sequence, ledgerCurrentIndex int,
) map[string]any {

	// builds the base payload for the RAW transaction
	return map[string]any{
		"TransactionType": "Payment",
		"Account":         walletFromAddress,
		"Destination":     walletToAddress,
		"DestinationTag":  1,
		"Amount": map[string]any{
			"currency": xrpn.ParseStringToHex(tokenAbbr),
			"issuer":   issuerAddress,
			"value":    amount,
		},
		"Flags":              flags,
		"Sequence":           sequence,
		"Fee":                xrpn.BASE_FEE,
		"LastLedgerSequence": ledgerCurrentIndex + xrpn.LEDGER_INCREMENT,
		"SigningPubKey":      publicKey,
	}
}

func parseStructToJson(data any) string {
	// converts a struct to a JSON string
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}
