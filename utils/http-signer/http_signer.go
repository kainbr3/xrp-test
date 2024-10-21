package httpsigner

import (
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt"
)

type HttpSigner struct {
	privateKey *rsa.PrivateKey
	apiKey     string
	rnd        *rand.Rand
}

// NewHttpSigner creates a new instance of HttpSigner with the provided RSA private key and API key.
// It initializes the random number generator with a custom source.
//
// Parameters:
// - pk: A pointer to an RSA private key used for signing JWT tokens.
// - apiKey: A string representing the API key to be included in the token claims.
//
// Returns:
// - A pointer to an initialized HttpSigner instance.
func NewHttpSigner(pk *rsa.PrivateKey, apiKey string) *HttpSigner {
	var s secrets
	hs := new(HttpSigner)
	hs.privateKey = pk
	hs.apiKey = apiKey
	hs.rnd = rand.New(s)
	return hs
}

// CreateAndSignJWTToken generates a JWT token with the provided path and body JSON,
// signs it using the RSA private key, and returns the signed token as a string.
//
// Parameters:
// - path: The URI path to be included in the token claims.
// - bodyJSON: The JSON string of the request body to be hashed and included in the token claims.
//
// Returns:
// - A string representing the signed JWT token.
// - An error if there was an issue signing the token.
func (k *HttpSigner) CreateAndSignJWTToken(path string, bodyJSON string) (string, error) {

	token := &jwt.MapClaims{
		"uri":      path,
		"nonce":    k.rnd.Int63(),
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Second * 55).Unix(),
		"sub":      k.apiKey,
		"bodyHash": createHash(bodyJSON),
	}

	j := jwt.NewWithClaims(jwt.SigningMethodRS256, token)
	signedToken, err := j.SignedString(k.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to generate the signed token with error: %v", err)
	}

	return signedToken, nil
}

// createHash generates a SHA-256 hash of the provided data string and returns
// the resulting hash as a hexadecimal encoded string.
//
// Parameters:
// - data: The input string to be hashed.
//
// Returns:
// - A string representing the hexadecimal encoded SHA-256 hash of the input data.
func createHash(data string) string {
	h := sha256.New()
	h.Write([]byte(data))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}
