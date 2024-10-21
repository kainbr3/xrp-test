package xrpscan

import (
	"context"
	k "crypto-braza-tokens-api/utils/keys-values"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Something(t *testing.T) {
	k.Start()
	scan, _ := NewXrpScanClient()
	result, err := scan.GetTokenObligations(context.Background(), "rH5CJsqvNqZGxrMyGaqLEoMWRYcVTAPZMt")
	assert.NoError(t, err)
	assert.NotNil(t, result)

	result2, err := scan.GetAccountXrpBalance(context.Background(), "rP1rFtLizETzwySJQTRKzLk7F5ZH7NmPqv")
	assert.NoError(t, err)
	assert.NotNil(t, result2)

	result3, err := scan.GetAccountTokensBalances(context.Background(), "rP1rFtLizETzwySJQTRKzLk7F5ZH7NmPqv")
	assert.NoError(t, err)
	assert.NotNil(t, result3)
}
