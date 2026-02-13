// Package crypto provides cryptographic capabilities.
// It supports encryption, decryption, hashing, and signing operations.
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/bcrypt"
)

// =============================================================================
// Errors
// =============================================================================

var (
	ErrInvalidKey        = errors.New("crypto: invalid key length")
	ErrInvalidCiphertext = errors.New("crypto: invalid ciphertext")
	ErrDecryptionFailed  = errors.New("crypto: decryption failed")
)

// =============================================================================
// AES Encryption
// =============================================================================

// Encryptor defines the interface for encryption operations.
type Encryptor interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
}

// AESEncryptor provides AES-GCM encryption.
type AESEncryptor struct {
	key []byte
}

// NewAESEncryptor creates a new AES encryptor.
// Key must be 16, 24, or 32 bytes for AES-128, AES-192, or AES-256.
func NewAESEncryptor(key []byte) (*AESEncryptor, error) {
	keyLen := len(key)
	if keyLen != 16 && keyLen != 24 && keyLen != 32 {
		return nil, ErrInvalidKey
	}
	return &AESEncryptor{key: key}, nil
}

// NewAESEncryptorFromString creates an AES encryptor from a string key.
// The key will be hashed to 32 bytes for AES-256.
func NewAESEncryptorFromString(key string) *AESEncryptor {
	hash := sha256.Sum256([]byte(key))
	return &AESEncryptor{key: hash[:]}
}

// Encrypt encrypts plaintext using AES-GCM.
func (e *AESEncryptor) Encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, fmt.Errorf("crypto: failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("crypto: failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("crypto: failed to generate nonce: %w", err)
	}

	// Prepend nonce to ciphertext
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// Decrypt decrypts ciphertext using AES-GCM.
func (e *AESEncryptor) Decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, fmt.Errorf("crypto: failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("crypto: failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, ErrInvalidCiphertext
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, ErrDecryptionFailed
	}

	return plaintext, nil
}

// EncryptString encrypts a string and returns base64-encoded ciphertext.
func (e *AESEncryptor) EncryptString(plaintext string) (string, error) {
	ciphertext, err := e.Encrypt([]byte(plaintext))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptString decrypts a base64-encoded ciphertext and returns the plaintext.
func (e *AESEncryptor) DecryptString(ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("crypto: failed to decode base64: %w", err)
	}

	plaintext, err := e.Decrypt(data)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// =============================================================================
// Hashing
// =============================================================================

// Hasher defines the interface for hashing operations.
type Hasher interface {
	Hash(data []byte) []byte
	HashString(data string) string
	Verify(data, hash []byte) bool
}

// SHA256Hash computes SHA-256 hash.
func SHA256Hash(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

// SHA256Hex computes SHA-256 hash and returns hex string.
func SHA256Hex(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// SHA512Hash computes SHA-512 hash.
func SHA512Hash(data []byte) []byte {
	hash := sha512.Sum512(data)
	return hash[:]
}

// SHA512Hex computes SHA-512 hash and returns hex string.
func SHA512Hex(data string) string {
	hash := sha512.Sum512([]byte(data))
	return hex.EncodeToString(hash[:])
}

// =============================================================================
// HMAC
// =============================================================================

// HMACSHA256 computes HMAC-SHA256.
func HMACSHA256(data, key []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

// HMACSHA256Hex computes HMAC-SHA256 and returns hex string.
func HMACSHA256Hex(data, key string) string {
	return hex.EncodeToString(HMACSHA256([]byte(data), []byte(key)))
}

// VerifyHMAC verifies HMAC signature.
func VerifyHMAC(data, signature, key []byte) bool {
	expected := HMACSHA256(data, key)
	return hmac.Equal(expected, signature)
}

// =============================================================================
// Password Hashing (bcrypt)
// =============================================================================

// HashPassword hashes a password using bcrypt.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("crypto: failed to hash password: %w", err)
	}
	return string(hash), nil
}

// VerifyPassword verifies a password against a bcrypt hash.
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// HashPasswordWithCost hashes a password with a specific bcrypt cost.
func HashPasswordWithCost(password string, cost int) (string, error) {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", fmt.Errorf("crypto: failed to hash password: %w", err)
	}
	return string(hash), nil
}

// =============================================================================
// Utility Functions
// =============================================================================

// GenerateKey generates a random key of the specified length.
func GenerateKey(length int) ([]byte, error) {
	key := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fmt.Errorf("crypto: failed to generate key: %w", err)
	}
	return key, nil
}

// GenerateKeyHex generates a random key and returns it as hex string.
func GenerateKeyHex(length int) (string, error) {
	key, err := GenerateKey(length)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(key), nil
}

// ConstantTimeCompare compares two byte slices in constant time.
func ConstantTimeCompare(a, b []byte) bool {
	return hmac.Equal(a, b)
}

// =============================================================================
// Convenience: Default Encryptor
// =============================================================================

var defaultEncryptor *AESEncryptor

// SetDefaultKey sets the default encryption key.
func SetDefaultKey(key string) {
	defaultEncryptor = NewAESEncryptorFromString(key)
}

// Encrypt encrypts using the default encryptor.
func Encrypt(plaintext string) (string, error) {
	if defaultEncryptor == nil {
		return "", errors.New("crypto: default key not set, call SetDefaultKey first")
	}
	return defaultEncryptor.EncryptString(plaintext)
}

// Decrypt decrypts using the default encryptor.
func Decrypt(ciphertext string) (string, error) {
	if defaultEncryptor == nil {
		return "", errors.New("crypto: default key not set, call SetDefaultKey first")
	}
	return defaultEncryptor.DecryptString(ciphertext)
}
