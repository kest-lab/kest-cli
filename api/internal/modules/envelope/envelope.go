package envelope

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"
)

// Common errors
var (
	ErrEmptyEnvelope     = errors.New("empty envelope")
	ErrInvalidHeader     = errors.New("invalid envelope header")
	ErrInvalidItemHeader = errors.New("invalid item header")
	ErrPayloadTooLarge   = errors.New("payload too large")
)

// MaxPayloadSize is the maximum size of an envelope payload (40MB)
const MaxPayloadSize = 40 * 1024 * 1024

// Envelope represents a Sentry-compatible envelope
type Envelope struct {
	Header *Header `json:"-"`
	Items  []*Item `json:"-"`
}

// Header represents the envelope header
type Header struct {
	EventID string            `json:"event_id,omitempty"`
	SentAt  time.Time         `json:"sent_at,omitempty"`
	Dsn     string            `json:"dsn,omitempty"`
	Sdk     *SdkInfo          `json:"sdk,omitempty"`
	Trace   map[string]string `json:"trace,omitempty"`
}

// SdkInfo contains SDK metadata
type SdkInfo struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

// ItemType represents the type of envelope item
type ItemType string

const (
	ItemTypeEvent       ItemType = "event"
	ItemTypeTransaction ItemType = "transaction"
	ItemTypeAttachment  ItemType = "attachment"
	ItemTypeSession     ItemType = "session"
	ItemTypeCheckIn     ItemType = "check_in"
	ItemTypeLog         ItemType = "log"
)

// ItemHeader represents the header of an envelope item
type ItemHeader struct {
	Type        ItemType `json:"type"`
	Length      *int     `json:"length,omitempty"`
	Filename    string   `json:"filename,omitempty"`
	ContentType string   `json:"content_type,omitempty"`
	ItemCount   *int     `json:"item_count,omitempty"`
}

// Item represents a single item within an envelope
type Item struct {
	Header  *ItemHeader `json:"-"`
	Payload []byte      `json:"-"`
}

// Parse parses a Sentry envelope from an io.Reader
// Envelope format:
//   - First line: JSON envelope header
//   - Subsequent pairs: JSON item header + payload
//   - Separated by newlines (\n)
func Parse(r io.Reader) (*Envelope, error) {
	reader := bufio.NewReader(r)

	// Read envelope header (first line)
	headerLine, err := reader.ReadBytes('\n')
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("failed to read envelope header: %w", err)
	}
	if len(headerLine) == 0 {
		return nil, ErrEmptyEnvelope
	}

	// Trim newline
	headerLine = bytes.TrimSuffix(headerLine, []byte("\n"))

	var header Header
	if err := json.Unmarshal(headerLine, &header); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidHeader, err)
	}

	envelope := &Envelope{
		Header: &header,
		Items:  make([]*Item, 0),
	}

	// Read items
	for {
		// Read item header
		itemHeaderLine, err := reader.ReadBytes('\n')
		if err == io.EOF {
			if len(itemHeaderLine) == 0 {
				break // End of envelope
			}
		} else if err != nil {
			return nil, fmt.Errorf("failed to read item header: %w", err)
		}

		// Skip empty lines
		itemHeaderLine = bytes.TrimSuffix(itemHeaderLine, []byte("\n"))
		if len(itemHeaderLine) == 0 {
			continue
		}

		var itemHeader ItemHeader
		if err := json.Unmarshal(itemHeaderLine, &itemHeader); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidItemHeader, err)
		}

		// Read payload
		var payload []byte
		if itemHeader.Length != nil && *itemHeader.Length > 0 {
			// Length is specified - read exact bytes
			length := *itemHeader.Length
			if length > MaxPayloadSize {
				return nil, ErrPayloadTooLarge
			}
			payload = make([]byte, length)
			if _, err := io.ReadFull(reader, payload); err != nil {
				return nil, fmt.Errorf("failed to read payload: %w", err)
			}
			// Consume trailing newline if present
			if b, err := reader.ReadByte(); err == nil && b != '\n' {
				reader.UnreadByte()
			}
		} else {
			// No length - read until newline
			payload, err = reader.ReadBytes('\n')
			if err != nil && err != io.EOF {
				return nil, fmt.Errorf("failed to read payload: %w", err)
			}
			payload = bytes.TrimSuffix(payload, []byte("\n"))
		}

		item := &Item{
			Header:  &itemHeader,
			Payload: payload,
		}
		envelope.Items = append(envelope.Items, item)
	}

	return envelope, nil
}

// ParseBytes parses a Sentry envelope from bytes
func ParseBytes(data []byte) (*Envelope, error) {
	return Parse(bytes.NewReader(data))
}

// GetEventItems returns all items of type "event"
func (e *Envelope) GetEventItems() []*Item {
	var events []*Item
	for _, item := range e.Items {
		if item.Header.Type == ItemTypeEvent {
			events = append(events, item)
		}
	}
	return events
}

// GetTransactionItems returns all items of type "transaction"
func (e *Envelope) GetTransactionItems() []*Item {
	var transactions []*Item
	for _, item := range e.Items {
		if item.Header.Type == ItemTypeTransaction {
			transactions = append(transactions, item)
		}
	}
	return transactions
}

// Serialize converts the envelope back to wire format
func (e *Envelope) Serialize() ([]byte, error) {
	var buf bytes.Buffer

	// Write header
	headerBytes, err := json.Marshal(e.Header)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal envelope header: %w", err)
	}
	buf.Write(headerBytes)
	buf.WriteByte('\n')

	// Write items
	for _, item := range e.Items {
		// Update length if needed
		if item.Header.Length == nil || *item.Header.Length != len(item.Payload) {
			length := len(item.Payload)
			item.Header.Length = &length
		}

		itemHeaderBytes, err := json.Marshal(item.Header)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal item header: %w", err)
		}
		buf.Write(itemHeaderBytes)
		buf.WriteByte('\n')
		buf.Write(item.Payload)
		buf.WriteByte('\n')
	}

	return buf.Bytes(), nil
}

// PayloadAsString returns the item payload as a string
func (i *Item) PayloadAsString() string {
	return string(i.Payload)
}

// PayloadAsJSON unmarshals the payload into the target
func (i *Item) PayloadAsJSON(target interface{}) error {
	return json.Unmarshal(i.Payload, target)
}

// String implements Stringer for Header
func (h *Header) String() string {
	return fmt.Sprintf("EventID=%s, SentAt=%s, Dsn=%s", h.EventID, h.SentAt, h.Dsn)
}

// String implements Stringer for ItemHeader
func (h *ItemHeader) String() string {
	length := 0
	if h.Length != nil {
		length = *h.Length
	}
	return fmt.Sprintf("Type=%s, Length=%d", h.Type, length)
}

// intPtr is a helper to create pointer to int
func intPtr(i int) *int {
	return &i
}

// NewItem creates a new envelope item
func NewItem(itemType ItemType, payload []byte) *Item {
	length := len(payload)
	return &Item{
		Header: &ItemHeader{
			Type:   itemType,
			Length: &length,
		},
		Payload: payload,
	}
}

// NewEnvelope creates a new envelope with the given event ID
func NewEnvelope(eventID string) *Envelope {
	return &Envelope{
		Header: &Header{
			EventID: eventID,
			SentAt:  time.Now().UTC(),
		},
		Items: make([]*Item, 0),
	}
}

// AddItem adds an item to the envelope
func (e *Envelope) AddItem(item *Item) {
	e.Items = append(e.Items, item)
}

// ExtractProjectIDFromPath extracts project ID from the API path
// Example: /api/123/envelope/ -> "123"
func ExtractProjectIDFromPath(path string) string {
	// Path format: /api/{project_id}/envelope/
	parts := bytes.Split([]byte(path), []byte("/"))
	for i, part := range parts {
		if string(part) == "api" && i+1 < len(parts) {
			return string(parts[i+1])
		}
	}
	return ""
}

// ExtractPublicKeyFromAuth extracts the public key from X-Sentry-Auth header
// Format: Sentry sentry_version=7, sentry_client=sentry.go/0.x.x, sentry_key=YOUR_KEY
func ExtractPublicKeyFromAuth(authHeader string) string {
	if authHeader == "" {
		return ""
	}

	// Parse the header
	parts := bytes.Split([]byte(authHeader), []byte(","))
	for _, part := range parts {
		part = bytes.TrimSpace(part)
		if bytes.HasPrefix(part, []byte("sentry_key=")) {
			return string(bytes.TrimPrefix(part, []byte("sentry_key=")))
		}
	}

	return ""
}

// ParseSentryVersion extracts sentry_version from auth header
func ParseSentryVersion(authHeader string) string {
	parts := bytes.Split([]byte(authHeader), []byte(","))
	for _, part := range parts {
		part = bytes.TrimSpace(part)
		if bytes.HasPrefix(part, []byte("Sentry sentry_version=")) {
			return string(bytes.TrimPrefix(part, []byte("Sentry sentry_version=")))
		}
		if bytes.HasPrefix(part, []byte("sentry_version=")) {
			val := bytes.TrimPrefix(part, []byte("sentry_version="))
			// Handle case like "Sentry sentry_version=7"
			valParts := bytes.Split(val, []byte(" "))
			if len(valParts) > 0 {
				return string(valParts[len(valParts)-1])
			}
			return string(val)
		}
	}
	return ""
}

// ValidateSentryVersion checks if the SDK version is supported
func ValidateSentryVersion(version string) bool {
	v, err := strconv.Atoi(version)
	if err != nil {
		return false
	}
	// Support version 7 (current) and future versions
	return v >= 7
}
