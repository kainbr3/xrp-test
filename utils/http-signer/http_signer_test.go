//go:build unit

package httpsigner

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func generateMockPrivateKey() *rsa.PrivateKey {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	return privateKey
}

func TestCases_HttpSigner_Unit(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{"Success creating a new HttpSigner", testNewHttpSigner},
		{"Success creating and signing JWT token", testCreateAndSignJWTToken},
		{"Success creating hash", testCreateHash},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

func testNewHttpSigner(t *testing.T) {
	privateKey := generateMockPrivateKey()
	apiKey := "testApiKey"
	signer := NewHttpSigner(privateKey, apiKey)

	assert.NotNil(t, signer)
	assert.Equal(t, privateKey, signer.privateKey)
	assert.Equal(t, apiKey, signer.apiKey)
	assert.NotNil(t, signer.rnd)
}

func testCreateAndSignJWTToken(t *testing.T) {
	privateKey := generateMockPrivateKey()
	apiKey := "testApiKey"
	signer := NewHttpSigner(privateKey, apiKey)

	path := "/test/path"
	bodyJSON := `{"key": "value"}`
	token, err := signer.CreateAndSignJWTToken(path, bodyJSON)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func testCreateHash(t *testing.T) {
	data := "testData"
	expectedHash := sha256.Sum256([]byte(data))
	expectedHashString := hex.EncodeToString(expectedHash[:])
	hash := createHash(data)

	assert.Equal(t, expectedHashString, hash)
}
