package account

import (
	"github.com/google/uuid"
)

// RandomIDGenerator is an IDGenerator using random
// This struct isn't tested because it only calls primitives from vendor uuid dependency
type RandomIDGenerator struct {
}

// NewRandomIDGenerator creates a new RandomIDGenerator
func NewRandomIDGenerator() RandomIDGenerator {
	return RandomIDGenerator{}
}

// Next return the next random UUID
func (r RandomIDGenerator) Next() uuid.UUID {
	return uuid.New()
}
