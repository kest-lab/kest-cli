package project

import (
	"testing"
)

func TestParseDSN(t *testing.T) {
	tests := []struct {
		name     string
		rawURL   string
		wantErr  bool
		wantKey  string
		wantHost string
		wantID   string
		wantPort int
	}{
		{
			name:     "basic https DSN",
			rawURL:   "https://abc123@trac.example.com/42",
			wantErr:  false,
			wantKey:  "abc123",
			wantHost: "trac.example.com",
			wantID:   "42",
			wantPort: 443,
		},
		{
			name:     "http DSN with port",
			rawURL:   "http://testkey@localhost:8025/1",
			wantErr:  false,
			wantKey:  "testkey",
			wantHost: "localhost",
			wantID:   "1",
			wantPort: 8025,
		},
		{
			name:     "DSN with secret key",
			rawURL:   "https://public:secret@sentry.io/123",
			wantErr:  false,
			wantKey:  "public",
			wantHost: "sentry.io",
			wantID:   "123",
			wantPort: 443,
		},
		{
			name:     "DSN with path",
			rawURL:   "https://key@example.com/prefix/path/99",
			wantErr:  false,
			wantKey:  "key",
			wantHost: "example.com",
			wantID:   "99",
			wantPort: 443,
		},
		{
			name:    "invalid scheme",
			rawURL:  "ftp://key@host/1",
			wantErr: true,
		},
		{
			name:    "missing public key",
			rawURL:  "https://sentry.io/1",
			wantErr: true,
		},
		{
			name:    "missing project id",
			rawURL:  "https://key@sentry.io/",
			wantErr: true,
		},
		{
			name:    "invalid URL",
			rawURL:  "not a valid url",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dsn, err := ParseDSN(tt.rawURL)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseDSN() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("ParseDSN() error = %v", err)
				return
			}
			if dsn.PublicKey != tt.wantKey {
				t.Errorf("PublicKey = %q, want %q", dsn.PublicKey, tt.wantKey)
			}
			if dsn.Host != tt.wantHost {
				t.Errorf("Host = %q, want %q", dsn.Host, tt.wantHost)
			}
			if dsn.ProjectID != tt.wantID {
				t.Errorf("ProjectID = %q, want %q", dsn.ProjectID, tt.wantID)
			}
			if dsn.Port != tt.wantPort {
				t.Errorf("Port = %d, want %d", dsn.Port, tt.wantPort)
			}
		})
	}
}

func TestDSNString(t *testing.T) {
	dsn := &DSN{
		Scheme:    "https",
		PublicKey: "abc123",
		Host:      "sentry.io",
		Port:      443,
		ProjectID: "42",
	}

	result := dsn.String()
	if result != "https://abc123@sentry.io/42" {
		t.Errorf("DSN.String() = %q, want %q", result, "https://abc123@sentry.io/42")
	}
}

func TestDSNGetAPIURL(t *testing.T) {
	dsn := &DSN{
		Scheme:    "https",
		PublicKey: "key",
		Host:      "example.com",
		Port:      443,
		ProjectID: "1",
	}

	result := dsn.GetAPIURL()
	want := "https://example.com/api/1/envelope/"
	if result != want {
		t.Errorf("DSN.GetAPIURL() = %q, want %q", result, want)
	}
}

func TestGenerateKey(t *testing.T) {
	key1 := GenerateKey()
	key2 := GenerateKey()

	// Keys should be 32 characters (16 bytes hex encoded)
	if len(key1) != 32 {
		t.Errorf("GenerateKey() length = %d, want 32", len(key1))
	}

	// Keys should be unique
	if key1 == key2 {
		t.Errorf("GenerateKey() should generate unique keys")
	}
}

func TestGenerateDSN(t *testing.T) {
	result := GenerateDSN("testkey", "localhost:8025", 1, false)
	want := "http://testkey@localhost:8025/1"
	if result != want {
		t.Errorf("GenerateDSN() = %q, want %q", result, want)
	}

	result = GenerateDSN("testkey", "sentry.io", 42, true)
	want = "https://testkey@sentry.io/42"
	if result != want {
		t.Errorf("GenerateDSN() = %q, want %q", result, want)
	}
}
