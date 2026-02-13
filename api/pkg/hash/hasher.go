package hash

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
)

// Hasher is the interface for password hashing
type Hasher interface {
	// Make creates a hash from a plain password
	Make(password string) (string, error)

	// Check verifies a password against a hash
	Check(password, hash string) bool

	// NeedsRehash checks if a hash needs to be re-hashed
	NeedsRehash(hash string) bool
}

// Algorithm represents a hashing algorithm
type Algorithm string

const (
	AlgorithmBcrypt Algorithm = "bcrypt"
	AlgorithmArgon2 Algorithm = "argon2id"
)

// DefaultAlgorithm is the default hashing algorithm
var DefaultAlgorithm Algorithm = AlgorithmBcrypt

// config holds the package configuration
var config = struct {
	bcryptCost   int
	argon2Config Argon2Config
}{
	bcryptCost: bcrypt.DefaultCost,
	argon2Config: Argon2Config{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	},
}

// Argon2Config holds Argon2 configuration
type Argon2Config struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

// --- Package-level Functions ---

// Make creates a hash using the default algorithm
func Make(password string) (string, error) {
	switch DefaultAlgorithm {
	case AlgorithmArgon2:
		return MakeArgon2(password)
	default:
		return MakeBcrypt(password)
	}
}

// Check verifies a password against a hash (auto-detects algorithm)
func Check(password, hash string) bool {
	if strings.HasPrefix(hash, "$argon2id$") {
		return CheckArgon2(password, hash)
	}
	return CheckBcrypt(password, hash)
}

// NeedsRehash checks if a hash needs to be upgraded
func NeedsRehash(hash string) bool {
	if strings.HasPrefix(hash, "$argon2id$") {
		return NeedsRehashArgon2(hash)
	}
	return NeedsRehashBcrypt(hash)
}

// --- Bcrypt Functions ---

// MakeBcrypt creates a bcrypt hash
func MakeBcrypt(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), config.bcryptCost)
	if err != nil {
		return "", fmt.Errorf("bcrypt hash failed: %w", err)
	}
	return string(bytes), nil
}

// MakeBcryptWithCost creates a bcrypt hash with custom cost
func MakeBcryptWithCost(password string, cost int) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", fmt.Errorf("bcrypt hash failed: %w", err)
	}
	return string(bytes), nil
}

// CheckBcrypt verifies a password against a bcrypt hash
func CheckBcrypt(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// NeedsRehashBcrypt checks if bcrypt hash needs rehashing
func NeedsRehashBcrypt(hash string) bool {
	cost, err := bcrypt.Cost([]byte(hash))
	if err != nil {
		return true
	}
	return cost < config.bcryptCost
}

// SetBcryptCost sets the default bcrypt cost
func SetBcryptCost(cost int) {
	if cost < bcrypt.MinCost {
		cost = bcrypt.MinCost
	}
	if cost > bcrypt.MaxCost {
		cost = bcrypt.MaxCost
	}
	config.bcryptCost = cost
}

// --- Argon2 Functions ---

// MakeArgon2 creates an argon2id hash
func MakeArgon2(password string) (string, error) {
	return MakeArgon2WithConfig(password, config.argon2Config)
}

// MakeArgon2WithConfig creates an argon2id hash with custom config
func MakeArgon2WithConfig(password string, cfg Argon2Config) (string, error) {
	salt := make([]byte, cfg.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("argon2 salt generation failed: %w", err)
	}

	hash := argon2.IDKey([]byte(password), salt, cfg.Iterations, cfg.Memory, cfg.Parallelism, cfg.KeyLength)

	// Encode to PHC string format
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, cfg.Memory, cfg.Iterations, cfg.Parallelism, b64Salt, b64Hash)

	return encoded, nil
}

// CheckArgon2 verifies a password against an argon2id hash
func CheckArgon2(password, encodedHash string) bool {
	cfg, salt, hash, err := decodeArgon2Hash(encodedHash)
	if err != nil {
		return false
	}

	computed := argon2.IDKey([]byte(password), salt, cfg.Iterations, cfg.Memory, cfg.Parallelism, cfg.KeyLength)

	return subtle.ConstantTimeCompare(hash, computed) == 1
}

// NeedsRehashArgon2 checks if argon2 hash needs rehashing
func NeedsRehashArgon2(encodedHash string) bool {
	cfg, _, _, err := decodeArgon2Hash(encodedHash)
	if err != nil {
		return true
	}

	return cfg.Memory != config.argon2Config.Memory ||
		cfg.Iterations != config.argon2Config.Iterations ||
		cfg.Parallelism != config.argon2Config.Parallelism
}

// SetArgon2Config sets the default Argon2 configuration
func SetArgon2Config(cfg Argon2Config) {
	config.argon2Config = cfg
}

// decodeArgon2Hash decodes a PHC format argon2id hash
func decodeArgon2Hash(encodedHash string) (Argon2Config, []byte, []byte, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return Argon2Config{}, nil, nil, fmt.Errorf("invalid argon2 hash format")
	}

	if parts[1] != "argon2id" {
		return Argon2Config{}, nil, nil, fmt.Errorf("unsupported argon2 variant")
	}

	var version int
	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		return Argon2Config{}, nil, nil, fmt.Errorf("invalid version")
	}

	var cfg Argon2Config
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &cfg.Memory, &cfg.Iterations, &cfg.Parallelism)
	if err != nil {
		return Argon2Config{}, nil, nil, fmt.Errorf("invalid parameters")
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return Argon2Config{}, nil, nil, fmt.Errorf("invalid salt")
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return Argon2Config{}, nil, nil, fmt.Errorf("invalid hash")
	}

	cfg.SaltLength = uint32(len(salt))
	cfg.KeyLength = uint32(len(hash))

	return cfg, salt, hash, nil
}

// --- Hasher Implementations ---

// BcryptHasher implements Hasher for bcrypt
type BcryptHasher struct {
	cost int
}

// NewBcryptHasher creates a new bcrypt hasher
func NewBcryptHasher(cost ...int) *BcryptHasher {
	c := config.bcryptCost
	if len(cost) > 0 {
		c = cost[0]
	}
	return &BcryptHasher{cost: c}
}

func (h *BcryptHasher) Make(password string) (string, error) {
	return MakeBcryptWithCost(password, h.cost)
}

func (h *BcryptHasher) Check(password, hash string) bool {
	return CheckBcrypt(password, hash)
}

func (h *BcryptHasher) NeedsRehash(hash string) bool {
	cost, err := bcrypt.Cost([]byte(hash))
	if err != nil {
		return true
	}
	return cost < h.cost
}

// Argon2Hasher implements Hasher for Argon2id
type Argon2Hasher struct {
	config Argon2Config
}

// NewArgon2Hasher creates a new Argon2 hasher
func NewArgon2Hasher(cfg ...Argon2Config) *Argon2Hasher {
	c := config.argon2Config
	if len(cfg) > 0 {
		c = cfg[0]
	}
	return &Argon2Hasher{config: c}
}

func (h *Argon2Hasher) Make(password string) (string, error) {
	return MakeArgon2WithConfig(password, h.config)
}

func (h *Argon2Hasher) Check(password, hash string) bool {
	return CheckArgon2(password, hash)
}

func (h *Argon2Hasher) NeedsRehash(hash string) bool {
	cfg, _, _, err := decodeArgon2Hash(hash)
	if err != nil {
		return true
	}
	return cfg.Memory != h.config.Memory ||
		cfg.Iterations != h.config.Iterations ||
		cfg.Parallelism != h.config.Parallelism
}
