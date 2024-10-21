package ripple

import (
	"context"
	kvs "crypto-braza-tokens-api/utils/keys-values"
	l "crypto-braza-tokens-api/utils/logger"
	"crypto-braza-tokens-api/utils/requests"
	"fmt"
	"strconv"

	"go.uber.org/zap"
)

var (
	BASE_FEE         string
	HASH_SIZE        int
	PREFIX_SIGNED    string
	PREFIX_UNSIGNED  string
	LEDGER_INCREMENT int
)

type RippleNodeClient struct {
	nodeApiUrl           string
	xrpScanApiUrl        string
	xrpScanExplorerUrl   string
	xrpLedgerExplorerUrl string
}

func NewRippleNodeClient() (*RippleNodeClient, error) {
	nodeApiUrl, err := kvs.Get("XRP_NODE_API_URL")
	if err != nil {
		l.Logger.Error("ripple client: error getting XRP_NODE_API_URL from KVS", zap.Error(err))
		return nil, err
	}

	xrpScanApiUrl, err := kvs.Get("XRP_SCAN_API_URL")
	if err != nil {
		l.Logger.Error("ripple client: error getting XRP_SCAN_API_URL from KVS", zap.Error(err))
		return nil, err
	}

	xrpScanExplorerUrl, err := kvs.Get("XRP_SCAN_EXPLORER_URL")
	if err != nil {
		l.Logger.Error("ripple client: error getting XRP_SCAN_EXPLORER_URL from KVS", zap.Error(err))
		return nil, err
	}

	xrpLedgerExplorerUrl, err := kvs.Get("XRP_LEDGER_EXPLORER_URL")
	if err != nil {
		l.Logger.Error("ripple client: error getting XRP_LEDGER_EXPLORER_URL from KVS", zap.Error(err))
		return nil, err
	}

	BASE_FEE, err = kvs.Get("XRP_BASE_FEE")
	if err != nil {
		l.Logger.Error("ripple client: error getting XRP_BASE_FEE from KVS", zap.Error(err))
		return nil, err
	}

	hashSizeStr, err := kvs.Get("XRP_HASH_SIZE")
	if err != nil {
		l.Logger.Error("ripple client: error getting XRP_HASH_SIZE from KVS", zap.Error(err))
		return nil, err
	}

	HASH_SIZE, err = strconv.Atoi(hashSizeStr)
	if err != nil {
		l.Logger.Error("ripple client: error converting hash size to int", zap.Error(err))
		return nil, fmt.Errorf("failed to convert hash size to int with error: %v", err)
	}

	ledgerIncremenetStr, err := kvs.Get("XRP_LEDGER_INCREMENT")
	if err != nil {
		l.Logger.Error("ripple client: error getting XRP_LEDGER_INCREMENT from KVS", zap.Error(err))
		return nil, err
	}

	LEDGER_INCREMENT, err = strconv.Atoi(ledgerIncremenetStr)
	if err != nil {
		l.Logger.Error("ripple client: error converting ledger incremene to int", zap.Error(err))
		return nil, fmt.Errorf("failed to convert ledger incremene to int with error: %v", err)
	}

	PREFIX_SIGNED, err = kvs.Get("XRP_PREFIX_SIGNED")
	if err != nil {
		l.Logger.Error("ripple client: error getting XRP_PREFIX_SIGNED from KVS", zap.Error(err))
		return nil, err
	}

	PREFIX_UNSIGNED, err = kvs.Get("XRP_PREFIX_UNSIGNED")
	if err != nil {
		l.Logger.Error("ripple client: error getting XRP_PREFIX_UNSIGNED from KVS", zap.Error(err))
		return nil, err
	}

	return &RippleNodeClient{nodeApiUrl, xrpScanApiUrl, xrpScanExplorerUrl, xrpLedgerExplorerUrl}, nil
}

func (r *RippleNodeClient) GetAccountInfo(ctx context.Context, address string) (*XrpAccountInfo, error) {
	request := &XrpJsonRpcRequest{
		Method: "account_info",
		Params: []any{
			map[string]any{
				"account":      address,
				"ledger_index": "current",
				"queue":        true,
			},
		},
	}

	parameters := map[string]any{"payload": request}

	result := &XrpAccountInfo{}

	err := requests.Execute(ctx, "POST", r.nodeApiUrl, &result, parameters)
	if err != nil {
		l.Logger.Error("ripple client: failed to retreive account info for address", zap.String("address", address), zap.Error(err))
		return nil, fmt.Errorf("failed to retreive account info for address: %s with error: %v", address, err)
	}

	return result, nil
}

func (r *RippleNodeClient) BuildRawTransactionRequest(txBlob string) *XrpJsonRpcRequest {
	return &XrpJsonRpcRequest{
		Method: "submit",
		Params: []any{
			map[string]any{"tx_blob": txBlob},
		},
	}
}

func (r *RippleNodeClient) BuildAccountLinesRequest(address string) *XrpJsonRpcRequest {
	return &XrpJsonRpcRequest{
		Method: "account_lines",
		Params: []any{
			map[string]any{"account": address},
		},
	}
}

func (r *RippleNodeClient) SubmitSignedTransaction(ctx context.Context, txBlob string, request *XrpJsonRpcRequest) (*SubmitTxResultResponse, error) {
	parameters := map[string]any{"payload": request}
	result := &SubmitTxResultResponse{}

	err := requests.Execute(ctx, "POST", r.nodeApiUrl, &result, parameters)
	if err != nil {
		l.Logger.Error("ripple client: failed to submit tx_blob", zap.String("tx_blob", txBlob), zap.Error(err))
		return nil, fmt.Errorf("failed to submit tx_blob \n%s \non ripple node with error: %v", txBlob, err)
	}

	return result, nil
}

func (r *RippleNodeClient) GetTransactionLink(txHash string) string {
	return r.xrpLedgerExplorerUrl + txHash
}

func (r *RippleNodeClient) GetAccountLines(ctx context.Context, address string) (*AccountLinesResponse, error) {
	request := r.BuildAccountLinesRequest(address)
	parameters := map[string]any{"payload": request}
	result := &AccountLinesResponse{}

	err := requests.Execute(ctx, "POST", r.nodeApiUrl, &result, parameters)
	if err != nil {
		l.Logger.Error("ripple client: failed to retreive account lines for address", zap.String("address", address), zap.Error(err))
		return nil, fmt.Errorf("failed to retreive account lines for address: %s with error: %v", address, err)
	}

	return result, nil
}
