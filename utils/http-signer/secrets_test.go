//go:build unit

package httpsigner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCases_Secrets_Unit(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{"Success generating a random uint64", testUint64},
		{"Success generating a random int63", testInt63},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

func testUint64(t *testing.T) {
	s := secrets{}
	value := s.Uint64()
	assert.NotZero(t, value, "Uint64 should generate a non-zero value")
}

func testInt63(t *testing.T) {
	s := secrets{}
	value := s.Int63()
	assert.True(t, value >= 0, "Int63 should generate a non-negative value")
	assert.True(t, value < (1<<63)-1, "Int63 should generate a value less than (1<<63)-1")
}
