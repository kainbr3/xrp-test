package fireblocks

import (
	"context"
	"crypto-braza-tokens-api/utils/files"
	sgn "crypto-braza-tokens-api/utils/http-signer"
	kvs "crypto-braza-tokens-api/utils/keys-values"
	l "crypto-braza-tokens-api/utils/logger"
	"crypto-braza-tokens-api/utils/requests"
	"encoding/json"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

type FireblocksClient struct {
	apiUrl string
	apiKey string
	signer *sgn.HttpSigner
}

func NewFireblocksClient() (*FireblocksClient, error) {
	fb := &FireblocksClient{
		apiKey: os.Getenv("FIREBLOCKS_API_KEY"),
	}
	var err error

	fb.apiUrl, err = kvs.Get("FIREBLOCKS_API_URL")
	if err != nil {
		l.Logger.Error("fireblocks client: failed to get fireblocks api url from kv store", zap.Error(err))
		return nil, err
	}

	secret := os.Getenv("FIREBLOCKS_API_SECRET")
	secretBytes, err := files.ReadFile(secret)
	if err != nil {
		l.Logger.Error("fireblocks client: failed to read secret file", zap.Error(err))
		return nil, fmt.Errorf("failed to read secret file with error: %v", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(secretBytes)
	if err != nil {
		l.Logger.Error("fireblocks client: failed to convert secret into pem encoded private key", zap.Error(err))
		return nil, fmt.Errorf("failed to convert secret into pem encoded private key with error: %v", err)
	}

	fb.signer = sgn.NewHttpSigner(privateKey, fb.apiKey)

	return fb, nil
}

func (f *FireblocksClient) GetPublicKeyInfoFromVaultAccount(ctx context.Context, vaultAccountId, asset_id string, change, addressIndex int) (*GetPbkInfoByVaultAccResponse, error) {
	path := fmt.Sprintf("/v1/vault/accounts/%s/%s/%d/%d/public_key_info?compressed=true", vaultAccountId, asset_id, change, addressIndex)
	endpoint := fmt.Sprintf("%s%s", f.apiUrl, path)

	parameters, err := f.createSignedRequest(path, nil)
	if err != nil {
		l.Logger.Error("fireblocks client: fireblocks client: failed to create signed request", zap.Error(err))
		return nil, err
	}

	result := &GetPbkInfoByVaultAccResponse{}

	err = requests.Execute(ctx, "GET", endpoint, &result, parameters)
	if err != nil {
		l.Logger.Error("fireblocks client: failed to execute request", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (f *FireblocksClient) GetVaultAccount(ctx context.Context, vaultAccountID string) (*VaultAccountResponse, error) {
	path := fmt.Sprintf("/v1/vault/accounts/%s", vaultAccountID)
	endpoint := fmt.Sprintf("%s%s", f.apiUrl, path)

	parameters, err := f.createSignedRequest(path, nil)
	if err != nil {
		l.Logger.Error("fireblocks client: failed to create signed request", zap.Error(err))
		return nil, err
	}

	result := &VaultAccountResponse{}

	err = requests.Execute(ctx, "GET", endpoint, &result, parameters)
	if err != nil {
		l.Logger.Error("fireblocks client: failed to execute request", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (f *FireblocksClient) GetVaultAccountAddresses(ctx context.Context, vaultAccountID, assetID string) ([]*VaultAccountAssetAddressResponse, error) {
	path := fmt.Sprintf("/v1/vault/accounts/%s/%s/addresses", vaultAccountID, assetID)
	endpoint := fmt.Sprintf("%s%s", f.apiUrl, path)

	parameters, err := f.createSignedRequest(path, nil)
	if err != nil {
		l.Logger.Error("fireblocks client: failed to create signed request", zap.Error(err))
		return nil, err
	}

	result := []*VaultAccountAssetAddressResponse{}

	err = requests.Execute(ctx, "GET", endpoint, &result, parameters)
	if err != nil {
		l.Logger.Error("fireblocks client: failed to execute request", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (f *FireblocksClient) GetVaultAccountAssetBalance(ctx context.Context, vaultAccountID, assetID string) (*VaultAssetResponse, error) {
	path := fmt.Sprintf("/v1/vault/accounts/%s/%s", vaultAccountID, assetID)
	endpoint := fmt.Sprintf("%s%s", f.apiUrl, path)

	parameters, err := f.createSignedRequest(path, nil)
	if err != nil {
		l.Logger.Error("fireblocks client: failed to create signed request", zap.Error(err))
		return nil, err
	}

	result := &VaultAssetResponse{}

	err = requests.Execute(ctx, "GET", endpoint, &result, parameters)
	if err != nil {
		l.Logger.Error("fireblocks client: failed to execute request", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (f *FireblocksClient) SubmitTransaction(ctx context.Context, payload any) (*SubmittedTransactionResponse, error) {
	path := "/v1/transactions"
	endpoint := fmt.Sprintf("%s%s", f.apiUrl, path)

	parameters, err := f.createSignedRequest(path, payload)
	if err != nil {
		l.Logger.Error("fireblocks client: failed to create signed request", zap.Error(err))
		return nil, err
	}

	result := &SubmittedTransactionResponse{}

	err = requests.Execute(ctx, "POST", endpoint, &result, parameters)
	if err != nil {
		l.Logger.Error("fireblocks client: failed to execute request", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (f *FireblocksClient) GetTransactionByID(ctx context.Context, transactionID string) (*TransactionByIdResponse, error) {
	path := fmt.Sprintf("/v1/transactions/%s", transactionID)
	endpoint := fmt.Sprintf("%s%s", f.apiUrl, path)

	parameters, err := f.createSignedRequest(path, nil)
	if err != nil {
		l.Logger.Error("fireblocks client: failed to create signed request", zap.Error(err))
		return nil, err
	}

	result := &TransactionByIdResponse{}

	err = requests.Execute(ctx, "GET", endpoint, &result, parameters)
	if err != nil {
		l.Logger.Error("fireblocks client: failed to execute request", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (f *FireblocksClient) BuildRawTransactionRequest(ctx context.Context, vaultAccountID, assetID, note, rawMessageContent string) *RawTransactionRequest {
	payload := &RawTransactionRequest{
		Operation: OPERATION_RAW,
		AssetID:   assetID,
		Note:      note,
		Source: &TargetVaultAccount{
			Type: TARGET_VAULT_ACCOUNT,
			ID:   vaultAccountID,
		},
		ExtraParameters: &RawTxExtraParameters{
			RawMessageData: &RawTxMessageData{
				Messages: []*RawTxContent{
					{
						Content: rawMessageContent,
					},
				},
			},
		},
	}

	return payload
}

func (f *FireblocksClient) BuildInternalTransactionRequest(ctx context.Context, sourceVaultId, destinaitionVaultId, assetId, amount, note, externalTxId string) *InternalTransactionRequest {
	payload := &InternalTransactionRequest{
		Operation: OPERATION_TRANSFER,
		AssetID:   assetId,
		Amount:    amount,
		Source: &TargetVaultAccount{
			Type: TARGET_VAULT_ACCOUNT,
			ID:   sourceVaultId,
		},
		Destination: &TargetVaultAccount{
			Type: TARGET_VAULT_ACCOUNT,
			ID:   destinaitionVaultId,
		},
		Note:         note,
		ExternalTxID: externalTxId,
	}

	return payload
}

func (f *FireblocksClient) createSignedRequest(path string, payload any) (map[string]any, error) {
	stringPayload := ""

	if payload != nil {
		marshalledPayload, err := json.Marshal(payload)
		if err != nil {
			l.Logger.Error("fireblocks client: failed to process json payload", zap.Error(err))
			return nil, fmt.Errorf("error processing json payload: %v", err)
		}

		stringPayload = string(marshalledPayload)
	}

	jwtToken, err := f.signer.CreateAndSignJWTToken(path, stringPayload)
	if err != nil {
		l.Logger.Error("fireblocks client: failed to sign with JWT token", zap.Error(err))
		return nil, fmt.Errorf("failed to sign with JWT token with err: %v", err)
	}

	parameters := map[string]any{
		"headers": map[string]string{
			"X-API-Key":     f.apiKey,
			"Authorization": fmt.Sprintf("Bearer %s", jwtToken),
		},
	}

	if payload != nil {
		parameters["payload"] = payload
	}

	return parameters, nil
}
