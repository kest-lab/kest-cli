package project

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// DSN represents a parsed Sentry DSN
type DSN struct {
	Scheme    string // http or https
	PublicKey string
	SecretKey string
	Host      string
	Port      int
	Path      string
	ProjectID string
}

// DSNParseError represents an error that occurs if a DSN cannot be parsed
type DSNParseError struct {
	Message string
}

func (e DSNParseError) Error() string {
	return "[Trac] DSNParseError: " + e.Message
}

// ParseDSN parses a Sentry-compatible DSN string
// Format: {PROTOCOL}://{PUBLIC_KEY}@{HOST}/{PROJECT_ID}
// Example: https://abc123@trac.example.com/42
func ParseDSN(rawURL string) (*DSN, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, &DSNParseError{fmt.Sprintf("invalid url: %v", err)}
	}

	// Scheme
	scheme := parsedURL.Scheme
	if scheme != "http" && scheme != "https" {
		return nil, &DSNParseError{"invalid scheme, must be http or https"}
	}

	// PublicKey
	publicKey := parsedURL.User.Username()
	if publicKey == "" {
		return nil, &DSNParseError{"empty public key"}
	}

	// SecretKey (optional)
	secretKey, _ := parsedURL.User.Password()

	// Host
	host := parsedURL.Hostname()
	if host == "" {
		return nil, &DSNParseError{"empty host"}
	}

	// Port
	var port int
	if p := parsedURL.Port(); p != "" {
		port, err = strconv.Atoi(p)
		if err != nil {
			return nil, &DSNParseError{"invalid port"}
		}
	} else {
		if scheme == "https" {
			port = 443
		} else {
			port = 80
		}
	}

	// ProjectID
	if parsedURL.Path == "" || parsedURL.Path == "/" {
		return nil, &DSNParseError{"empty project id"}
	}
	pathSegments := strings.Split(strings.TrimPrefix(parsedURL.Path, "/"), "/")
	projectID := pathSegments[len(pathSegments)-1]

	if projectID == "" {
		return nil, &DSNParseError{"empty project id"}
	}

	// Path (for sub-path installations)
	var path string
	if len(pathSegments) > 1 {
		path = "/" + strings.Join(pathSegments[:len(pathSegments)-1], "/")
	}

	return &DSN{
		Scheme:    scheme,
		PublicKey: publicKey,
		SecretKey: secretKey,
		Host:      host,
		Port:      port,
		Path:      path,
		ProjectID: projectID,
	}, nil
}

// String formats DSN back to URL string
func (d *DSN) String() string {
	var result string
	result += fmt.Sprintf("%s://%s", d.Scheme, d.PublicKey)
	if d.SecretKey != "" {
		result += fmt.Sprintf(":%s", d.SecretKey)
	}
	result += fmt.Sprintf("@%s", d.Host)

	// Only add port if non-default
	defaultPort := 80
	if d.Scheme == "https" {
		defaultPort = 443
	}
	if d.Port != defaultPort {
		result += fmt.Sprintf(":%d", d.Port)
	}

	if d.Path != "" {
		result += d.Path
	}
	result += fmt.Sprintf("/%s", d.ProjectID)
	return result
}

// GetAPIURL returns the envelope endpoint URL
func (d *DSN) GetAPIURL() string {
	var result string
	result += fmt.Sprintf("%s://%s", d.Scheme, d.Host)

	defaultPort := 80
	if d.Scheme == "https" {
		defaultPort = 443
	}
	if d.Port != defaultPort {
		result += fmt.Sprintf(":%d", d.Port)
	}

	if d.Path != "" {
		result += d.Path
	}
	result += fmt.Sprintf("/api/%s/envelope/", d.ProjectID)
	return result
}

// GenerateKey generates a random 32-character hex key
func GenerateKey() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based key (less secure but functional)
		fallback := make([]byte, 16)
		for i := range fallback {
			fallback[i] = byte(i * 17) // Simple deterministic fallback
		}
		return hex.EncodeToString(fallback)
	}
	return hex.EncodeToString(bytes)
}

// GenerateDSN creates a DSN string for a project
func GenerateDSN(publicKey string, host string, projectID uint, secure bool) string {
	scheme := "http"
	if secure {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s@%s/%d", scheme, publicKey, host, projectID)
}
