// Package idgen provides ID generation capabilities.
// It supports multiple ID generation strategies including Snowflake, UUID, and short IDs.
package idgen

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Generator defines the interface for ID generation.
type Generator interface {
	// Generate generates a new unique ID.
	Generate() string

	// GenerateN generates n unique IDs.
	GenerateN(n int) []string
}

// =============================================================================
// UUID Generator
// =============================================================================

// UUIDGenerator generates UUID v4 IDs.
type UUIDGenerator struct{}

// NewUUIDGenerator creates a new UUID generator.
func NewUUIDGenerator() *UUIDGenerator {
	return &UUIDGenerator{}
}

// Generate generates a new UUID v4.
func (g *UUIDGenerator) Generate() string {
	return uuid.New().String()
}

// GenerateN generates n UUIDs.
func (g *UUIDGenerator) GenerateN(n int) []string {
	ids := make([]string, n)
	for i := 0; i < n; i++ {
		ids[i] = g.Generate()
	}
	return ids
}

// =============================================================================
// Snowflake Generator
// =============================================================================

const (
	// Snowflake bit allocation
	timestampBits = 41
	nodeBits      = 10
	sequenceBits  = 12

	maxNodeID   = (1 << nodeBits) - 1
	maxSequence = (1 << sequenceBits) - 1

	// Custom epoch: 2020-01-01 00:00:00 UTC
	epoch = 1577836800000
)

// SnowflakeGenerator generates Snowflake IDs.
type SnowflakeGenerator struct {
	mu       sync.Mutex
	nodeID   int64
	sequence int64
	lastTime int64
}

// NewSnowflakeGenerator creates a new Snowflake generator.
func NewSnowflakeGenerator(nodeID int64) (*SnowflakeGenerator, error) {
	if nodeID < 0 || nodeID > maxNodeID {
		return nil, fmt.Errorf("node ID must be between 0 and %d", maxNodeID)
	}
	return &SnowflakeGenerator{
		nodeID:   nodeID,
		sequence: 0,
	}, nil
}

// Generate generates a new Snowflake ID.
func (g *SnowflakeGenerator) Generate() string {
	g.mu.Lock()
	defer g.mu.Unlock()

	now := time.Now().UnixMilli()

	if now == g.lastTime {
		g.sequence = (g.sequence + 1) & maxSequence
		if g.sequence == 0 {
			// Wait for next millisecond
			for now <= g.lastTime {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		g.sequence = 0
	}

	g.lastTime = now

	id := ((now - epoch) << (nodeBits + sequenceBits)) |
		(g.nodeID << sequenceBits) |
		g.sequence

	return fmt.Sprintf("%d", id)
}

// GenerateN generates n Snowflake IDs.
func (g *SnowflakeGenerator) GenerateN(n int) []string {
	ids := make([]string, n)
	for i := 0; i < n; i++ {
		ids[i] = g.Generate()
	}
	return ids
}

// GenerateInt64 generates a Snowflake ID as int64.
func (g *SnowflakeGenerator) GenerateInt64() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()

	now := time.Now().UnixMilli()

	if now == g.lastTime {
		g.sequence = (g.sequence + 1) & maxSequence
		if g.sequence == 0 {
			for now <= g.lastTime {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		g.sequence = 0
	}

	g.lastTime = now

	return ((now - epoch) << (nodeBits + sequenceBits)) |
		(g.nodeID << sequenceBits) |
		g.sequence
}

// =============================================================================
// Short ID Generator
// =============================================================================

// ShortIDGenerator generates short, URL-safe IDs.
type ShortIDGenerator struct {
	length int
}

// NewShortIDGenerator creates a new short ID generator.
func NewShortIDGenerator(length int) *ShortIDGenerator {
	if length < 4 {
		length = 8
	}
	return &ShortIDGenerator{length: length}
}

// Generate generates a new short ID.
func (g *ShortIDGenerator) Generate() string {
	bytes := make([]byte, g.length)
	_, _ = rand.Read(bytes)
	return base64.RawURLEncoding.EncodeToString(bytes)[:g.length]
}

// GenerateN generates n short IDs.
func (g *ShortIDGenerator) GenerateN(n int) []string {
	ids := make([]string, n)
	for i := 0; i < n; i++ {
		ids[i] = g.Generate()
	}
	return ids
}

// =============================================================================
// Nano ID Generator
// =============================================================================

const nanoAlphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// NanoIDGenerator generates NanoID-style IDs.
type NanoIDGenerator struct {
	length   int
	alphabet string
}

// NewNanoIDGenerator creates a new NanoID generator.
func NewNanoIDGenerator(length int) *NanoIDGenerator {
	if length < 1 {
		length = 21
	}
	return &NanoIDGenerator{
		length:   length,
		alphabet: nanoAlphabet,
	}
}

// Generate generates a new NanoID.
func (g *NanoIDGenerator) Generate() string {
	bytes := make([]byte, g.length)
	_, _ = rand.Read(bytes)

	result := make([]byte, g.length)
	for i := 0; i < g.length; i++ {
		result[i] = g.alphabet[bytes[i]%byte(len(g.alphabet))]
	}
	return string(result)
}

// GenerateN generates n NanoIDs.
func (g *NanoIDGenerator) GenerateN(n int) []string {
	ids := make([]string, n)
	for i := 0; i < n; i++ {
		ids[i] = g.Generate()
	}
	return ids
}

// =============================================================================
// Convenience Functions
// =============================================================================

var (
	defaultUUID      = NewUUIDGenerator()
	defaultShortID   = NewShortIDGenerator(8)
	defaultNanoID    = NewNanoIDGenerator(21)
	defaultSnowflake *SnowflakeGenerator
	snowflakeOnce    sync.Once
)

// UUID generates a new UUID v4.
func UUID() string {
	return defaultUUID.Generate()
}

// ShortID generates a short URL-safe ID.
func ShortID() string {
	return defaultShortID.Generate()
}

// NanoID generates a NanoID.
func NanoID() string {
	return defaultNanoID.Generate()
}

// Snowflake generates a Snowflake ID.
func Snowflake() string {
	snowflakeOnce.Do(func() {
		defaultSnowflake, _ = NewSnowflakeGenerator(1)
	})
	return defaultSnowflake.Generate()
}

// SnowflakeInt64 generates a Snowflake ID as int64.
func SnowflakeInt64() int64 {
	snowflakeOnce.Do(func() {
		defaultSnowflake, _ = NewSnowflakeGenerator(1)
	})
	return defaultSnowflake.GenerateInt64()
}

// RandomHex generates a random hex string.
func RandomHex(length int) string {
	bytes := make([]byte, (length+1)/2)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)[:length]
}

// RandomBytes generates random bytes.
func RandomBytes(length int) []byte {
	bytes := make([]byte, length)
	_, _ = rand.Read(bytes)
	return bytes
}
