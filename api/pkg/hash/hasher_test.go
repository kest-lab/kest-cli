package hash

import (
	"testing"
)

func TestBcrypt_MakeAndCheck(t *testing.T) {
	password := "secret123"

	hash, err := MakeBcrypt(password)
	if err != nil {
		t.Fatalf("MakeBcrypt failed: %v", err)
	}

	if !CheckBcrypt(password, hash) {
		t.Error("CheckBcrypt should return true for correct password")
	}

	if CheckBcrypt("wrong", hash) {
		t.Error("CheckBcrypt should return false for wrong password")
	}
}

func TestBcrypt_NeedsRehash(t *testing.T) {
	// Create hash with low cost
	hash, _ := MakeBcryptWithCost("password", 4)

	// Set higher default cost
	SetBcryptCost(10)

	if !NeedsRehashBcrypt(hash) {
		t.Error("NeedsRehashBcrypt should return true for low cost hash")
	}

	// Reset
	SetBcryptCost(10)
}

func TestArgon2_MakeAndCheck(t *testing.T) {
	password := "secret123"

	hash, err := MakeArgon2(password)
	if err != nil {
		t.Fatalf("MakeArgon2 failed: %v", err)
	}

	if !CheckArgon2(password, hash) {
		t.Error("CheckArgon2 should return true for correct password")
	}

	if CheckArgon2("wrong", hash) {
		t.Error("CheckArgon2 should return false for wrong password")
	}
}

func TestArgon2_HashFormat(t *testing.T) {
	hash, _ := MakeArgon2("password")

	// Should start with $argon2id$
	if hash[:10] != "$argon2id$" {
		t.Errorf("Hash should start with $argon2id$, got: %s", hash[:10])
	}
}

func TestCheck_AutoDetect(t *testing.T) {
	password := "secret123"

	// Test bcrypt
	bcryptHash, _ := MakeBcrypt(password)
	if !Check(password, bcryptHash) {
		t.Error("Check should work with bcrypt hash")
	}

	// Test argon2
	argon2Hash, _ := MakeArgon2(password)
	if !Check(password, argon2Hash) {
		t.Error("Check should work with argon2 hash")
	}
}

func TestMake_DefaultAlgorithm(t *testing.T) {
	password := "secret123"

	// Test with bcrypt (default)
	DefaultAlgorithm = AlgorithmBcrypt
	hash1, _ := Make(password)
	if !Check(password, hash1) {
		t.Error("Make should work with bcrypt")
	}

	// Test with argon2
	DefaultAlgorithm = AlgorithmArgon2
	hash2, _ := Make(password)
	if !Check(password, hash2) {
		t.Error("Make should work with argon2")
	}

	// Reset
	DefaultAlgorithm = AlgorithmBcrypt
}

func TestBcryptHasher(t *testing.T) {
	hasher := NewBcryptHasher(10)
	password := "secret123"

	hash, err := hasher.Make(password)
	if err != nil {
		t.Fatalf("BcryptHasher.Make failed: %v", err)
	}

	if !hasher.Check(password, hash) {
		t.Error("BcryptHasher.Check should return true")
	}
}

func TestArgon2Hasher(t *testing.T) {
	hasher := NewArgon2Hasher()
	password := "secret123"

	hash, err := hasher.Make(password)
	if err != nil {
		t.Fatalf("Argon2Hasher.Make failed: %v", err)
	}

	if !hasher.Check(password, hash) {
		t.Error("Argon2Hasher.Check should return true")
	}
}

func BenchmarkBcrypt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MakeBcrypt("password")
	}
}

func BenchmarkArgon2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MakeArgon2("password")
	}
}
