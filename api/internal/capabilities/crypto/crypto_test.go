package crypto

import (
	"testing"
)

func TestAESEncryptor(t *testing.T) {
	encryptor := NewAESEncryptorFromString("my-secret-key")

	plaintext := "Hello, World! 你好世界"

	// Test encryption and decryption
	ciphertext, err := encryptor.EncryptString(plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	t.Logf("Ciphertext: %s", ciphertext)

	decrypted, err := encryptor.DecryptString(ciphertext)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("Decrypted = %q, want %q", decrypted, plaintext)
	}
}

func TestSHA256(t *testing.T) {
	hash := SHA256Hex("hello")
	expected := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"

	if hash != expected {
		t.Errorf("SHA256Hex = %s, want %s", hash, expected)
	}
}

func TestHMACSHA256(t *testing.T) {
	signature := HMACSHA256Hex("message", "secret")
	if len(signature) != 64 {
		t.Errorf("HMAC length = %d, want 64", len(signature))
	}
	t.Logf("HMAC: %s", signature)
}

func TestVerifyHMAC(t *testing.T) {
	key := []byte("secret")
	data := []byte("message")
	signature := HMACSHA256(data, key)

	if !VerifyHMAC(data, signature, key) {
		t.Error("HMAC verification failed")
	}

	// Wrong data should fail
	if VerifyHMAC([]byte("wrong"), signature, key) {
		t.Error("HMAC verification should fail for wrong data")
	}
}

func TestPasswordHashing(t *testing.T) {
	password := "myP@ssw0rd!"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	t.Logf("Password hash: %s", hash)

	// Verify correct password
	if !VerifyPassword(password, hash) {
		t.Error("VerifyPassword should return true for correct password")
	}

	// Verify wrong password
	if VerifyPassword("wrong", hash) {
		t.Error("VerifyPassword should return false for wrong password")
	}
}

func TestGenerateKey(t *testing.T) {
	key, err := GenerateKeyHex(32)
	if err != nil {
		t.Fatalf("GenerateKeyHex failed: %v", err)
	}

	if len(key) != 64 { // 32 bytes = 64 hex chars
		t.Errorf("Key length = %d, want 64", len(key))
	}

	t.Logf("Generated key: %s", key)
}
