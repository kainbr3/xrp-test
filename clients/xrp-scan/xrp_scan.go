package xrpscan

import (
	"context"
	kvs "crypto-braza-tokens-api/utils/keys-values"
	"crypto-braza-tokens-api/utils/requests"
	"fmt"
)

type XrpScanClient struct {
	apiUrl string
}

func NewXrpScanClient() (*XrpScanClient, error) {
	apiUrl, err := kvs.Get("XRP_SCAN_API_URL")
	if err != nil {
		return nil, err
	}

	return &XrpScanClient{apiUrl}, nil
}

func (c *XrpScanClient) GetTokenObligations(ctx context.Context, tokenAddress string) (*TokenObligationsResponse, error) {
	endpoint := fmt.Sprintf("%s/api/v1/account/%s/obligations", c.apiUrl, tokenAddress)

	result := []*TokenObligationsResponse{}

	err := requests.Execute(ctx, "GET", endpoint, &result, nil)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no obligations found for token address %s", tokenAddress)
	}

	return result[0], nil
}

func (c *XrpScanClient) GetAccountXrpBalance(ctx context.Context, tokenAddress string) (*AccountInfo, error) {
	endpoint := fmt.Sprintf("%s/api/v1/account/%s", c.apiUrl, tokenAddress)

	result := &AccountInfo{}

	err := requests.Execute(ctx, "GET", endpoint, &result, nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *XrpScanClient) GetAccountTokensBalances(ctx context.Context, tokenAddress string) ([]*AccountInfoTokenBalance, error) {
	endpoint := fmt.Sprintf("%s/api/v1/account/%s/assets", c.apiUrl, tokenAddress)

	result := []*AccountInfoTokenBalance{}

	err := requests.Execute(ctx, "GET", endpoint, &result, nil)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no token balances found for address %s", tokenAddress)
	}

	return result, nil
}
