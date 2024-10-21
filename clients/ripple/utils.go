package ripple

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"strings"
)

// ConvertStringToHex converts a string to its hexadecimal representation
func ConvertStringToHex(input string) string {
	return hex.EncodeToString([]byte(input))
}

// ConvertHexToString converts a hexadecimal representation to its string representation
func ConvertHexToString(hexStr string) (string, error) {
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return "", fmt.Errorf("failed to decode hex string: %v", err)
	}
	return string(bytes), nil
}

// ParseStringToHex returns the original string representation of the input hex string
// if its length is greater than 3 characters, otherwise returns the input string.
// If the length is greater than 3 characters, the result is the original string without trailing zeros.
func ParseStringToHex(input string) string {
	if len(input) > 3 {
		hexString := ConvertStringToHex(input)
		upperHexString := strings.ToUpper(hexString)
		paddedHexString := upperHexString + strings.Repeat("0", 40-len(upperHexString))
		return paddedHexString
	}
	return input
}

// ParseHexToString converts a hexadecimal representation to its abbreviated string representation
// by taking the first two characters and the last two characters of the original string.
func ParseHexToString(hexStr string) (string, error) {
	originalString, err := ConvertHexToString(hexStr)
	if err != nil {
		return "", fmt.Errorf("failed to decode hex string: %v", err)
	}

	if len(originalString) < 4 {
		return originalString, nil
	}

	abbr := originalString[:2] + originalString[len(originalString)-2:]
	return abbr, nil
}

// Sha512Half computes the SHA-512 hash of the input hex string and returns the first n characters in uppercase according to the hash size
func Sha512Half(hashSize int, hexStr string) (string, error) {
	// Decode the hex string to bytes
	data, err := hex.DecodeString(hexStr)
	if err != nil {
		return "", fmt.Errorf("failed to decode hex string: %v", err)
	}

	// Compute the SHA-512 hash
	hash := sha512.Sum512(data)

	// Convert the hash to a hex string
	hashHex := hex.EncodeToString(hash[:])

	// Return the first n characters in uppercase according to the hash size
	return strings.ToUpper(hashHex[:hashSize]), nil
}

// ConcactPrefixWithTxBlob concats the prefix with the txBlobHex
func ConcactPrefixWithTxBlob(prefix, txBlobHex string) string {
	var sb strings.Builder
	sb.WriteString(prefix)
	sb.WriteString(txBlobHex)
	return sb.String()
}

// EncodeDER encodes r and s values into DER format
func EncodeDER(rHex, sHex string) (string, error) {
	rBytes, err := hex.DecodeString(rHex)
	if err != nil {
		return "", fmt.Errorf("failed to decode rHex: %v", err)
	}

	sBytes, err := hex.DecodeString(sHex)
	if err != nil {
		return "", fmt.Errorf("failed to decode sHex: %v", err)
	}

	// Ensure the r value is correctly encoded with the 00 byte
	if rBytes[0]&0x80 != 0 {
		rBytes = append([]byte{0x00}, rBytes...)
	}

	// Ensure the s value is correctly encoded with the 00 byte
	if sBytes[0]&0x80 != 0 {
		sBytes = append([]byte{0x00}, sBytes...)
	}

	// Construct the DER encoded signature
	der := append([]byte{0x30, byte(4 + len(rBytes) + len(sBytes))}, append([]byte{0x02, byte(len(rBytes))}, rBytes...)...)
	der = append(der, append([]byte{0x02, byte(len(sBytes))}, sBytes...)...)

	return hex.EncodeToString(der), nil
}
