package httpsigner

import (
	l "crypto-braza-tokens-api/utils/logger"
	crand "crypto/rand"

	"encoding/binary"

	"go.uber.org/zap"
)

type secrets struct{}

// Seed is a placeholder function for setting a seed value.
// Currently, it does not perform any operations.
func (s secrets) Seed(seed int64) {}

// Uint64 generates a random uint64 value using crypto/rand.
// It logs an error if the random number generation fails.
func (s secrets) Uint64() (r uint64) {
	err := binary.Read(crand.Reader, binary.BigEndian, &r)
	if err != nil {
		l.Logger.Error("utils: failed to read random number", zap.Error(err))
	}
	return r
}

// Int63 generates a random int64 value in the range [0, 1<<63).
// It uses the Uint64 method and masks the highest bit to ensure the value is non-negative.
func (s secrets) Int63() int64 {
	return int64(s.Uint64() & ^uint64(1<<63))
}
