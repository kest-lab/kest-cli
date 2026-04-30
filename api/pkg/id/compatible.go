package id

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

var ErrInvalidIDFormat = errors.New("invalid ID format")

// Compatible accepts UUID strings and legacy numeric IDs while normalizing them
// to the string form used across the API.
type Compatible string

func (id Compatible) String() string {
	return string(id)
}

func Parse(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", ErrInvalidIDFormat
	}

	if _, err := uuid.Parse(raw); err == nil {
		return raw, nil
	}

	if _, err := strconv.ParseUint(raw, 10, 64); err == nil {
		return raw, nil
	}

	return "", ErrInvalidIDFormat
}

func Normalize(value any) (string, error) {
	switch typed := value.(type) {
	case nil:
		return "", ErrInvalidIDFormat
	case Compatible:
		return Parse(string(typed))
	case *Compatible:
		if typed == nil {
			return "", ErrInvalidIDFormat
		}
		return Parse(string(*typed))
	case string:
		return Parse(typed)
	case []byte:
		return Parse(string(typed))
	case json.Number:
		return Parse(typed.String())
	case int:
		return strconv.Itoa(typed), nil
	case int8:
		return strconv.FormatInt(int64(typed), 10), nil
	case int16:
		return strconv.FormatInt(int64(typed), 10), nil
	case int32:
		return strconv.FormatInt(int64(typed), 10), nil
	case int64:
		return strconv.FormatInt(typed, 10), nil
	case uint:
		return strconv.FormatUint(uint64(typed), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(typed), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(typed), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(typed), 10), nil
	case uint64:
		return strconv.FormatUint(typed, 10), nil
	case float64:
		if typed < 0 || math.Trunc(typed) != typed {
			return "", ErrInvalidIDFormat
		}
		return strconv.FormatUint(uint64(typed), 10), nil
	case float32:
		if typed < 0 || math.Trunc(float64(typed)) != float64(typed) {
			return "", ErrInvalidIDFormat
		}
		return strconv.FormatUint(uint64(typed), 10), nil
	default:
		return "", fmt.Errorf("%w: %T", ErrInvalidIDFormat, value)
	}
}

func (id *Compatible) UnmarshalJSON(data []byte) error {
	trimmed := bytes.TrimSpace(data)
	if bytes.Equal(trimmed, []byte("null")) {
		*id = ""
		return nil
	}

	var stringValue string
	if err := json.Unmarshal(trimmed, &stringValue); err == nil {
		normalized, parseErr := Parse(stringValue)
		if parseErr != nil {
			return parseErr
		}
		*id = Compatible(normalized)
		return nil
	}

	decoder := json.NewDecoder(bytes.NewReader(trimmed))
	decoder.UseNumber()

	var numberValue json.Number
	if err := decoder.Decode(&numberValue); err == nil {
		normalized, parseErr := Parse(numberValue.String())
		if parseErr != nil {
			return parseErr
		}
		*id = Compatible(normalized)
		return nil
	}

	return ErrInvalidIDFormat
}

func (id *Compatible) UnmarshalText(text []byte) error {
	normalized, err := Parse(string(text))
	if err != nil {
		return err
	}

	*id = Compatible(normalized)
	return nil
}
