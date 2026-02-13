package support

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"strings"
	"unicode"
)

// Str provides common string manipulation utilities.
type Str struct{}

// After returns everything after the given value in a string
func (Str) After(subject, search string) string {
	if idx := strings.Index(subject, search); idx >= 0 {
		return subject[idx+len(search):]
	}
	return subject
}

// AfterLast returns everything after the last occurrence of the given value
func (Str) AfterLast(subject, search string) string {
	if idx := strings.LastIndex(subject, search); idx >= 0 {
		return subject[idx+len(search):]
	}
	return subject
}

// Before returns everything before the given value in a string
func (Str) Before(subject, search string) string {
	if idx := strings.Index(subject, search); idx >= 0 {
		return subject[:idx]
	}
	return subject
}

// BeforeLast returns everything before the last occurrence of the given value
func (Str) BeforeLast(subject, search string) string {
	if idx := strings.LastIndex(subject, search); idx >= 0 {
		return subject[:idx]
	}
	return subject
}

// Between returns the portion of a string between two values
func (Str) Between(subject, from, to string) string {
	if from == "" || to == "" {
		return subject
	}
	s := Str{}.After(subject, from)
	return Str{}.Before(s, to)
}

// Camel converts a string to camelCase
func (Str) Camel(value string) string {
	words := splitWords(value)
	for i := range words {
		if i == 0 {
			words[i] = strings.ToLower(words[i])
		} else {
			words[i] = strings.Title(strings.ToLower(words[i]))
		}
	}
	return strings.Join(words, "")
}

// Contains checks if a string contains a given substring
func (Str) Contains(haystack, needle string) bool {
	return strings.Contains(haystack, needle)
}

// ContainsAll checks if a string contains all given substrings
func (Str) ContainsAll(haystack string, needles []string) bool {
	for _, needle := range needles {
		if !strings.Contains(haystack, needle) {
			return false
		}
	}
	return true
}

// ContainsAny checks if a string contains any of the given substrings
func (Str) ContainsAny(haystack string, needles []string) bool {
	for _, needle := range needles {
		if strings.Contains(haystack, needle) {
			return true
		}
	}
	return false
}

// EndsWith checks if a string ends with a given suffix
func (Str) EndsWith(haystack, needle string) bool {
	return strings.HasSuffix(haystack, needle)
}

// Finish caps a string with a single instance of a given value
func (Str) Finish(value, cap string) string {
	if !strings.HasSuffix(value, cap) {
		return value + cap
	}
	return value
}

// Headline converts a string to a headline format
func (Str) Headline(value string) string {
	words := splitWords(value)
	for i := range words {
		words[i] = strings.Title(strings.ToLower(words[i]))
	}
	return strings.Join(words, " ")
}

// IsEmpty checks if a string is empty or contains only whitespace
func (Str) IsEmpty(value string) bool {
	return strings.TrimSpace(value) == ""
}

// IsNotEmpty checks if a string is not empty
func (Str) IsNotEmpty(value string) bool {
	return strings.TrimSpace(value) != ""
}

// Kebab converts a string to kebab-case
func (Str) Kebab(value string) string {
	words := splitWords(value)
	for i := range words {
		words[i] = strings.ToLower(words[i])
	}
	return strings.Join(words, "-")
}

// Length returns the length of a string
func (Str) Length(value string) int {
	return len([]rune(value))
}

// Limit truncates a string to the given length
func (Str) Limit(value string, limit int, end string) string {
	runes := []rune(value)
	if len(runes) <= limit {
		return value
	}
	if end == "" {
		end = "..."
	}
	return string(runes[:limit]) + end
}

// Lower converts a string to lowercase
func (Str) Lower(value string) string {
	return strings.ToLower(value)
}

// Mask masks a portion of a string with a repeated character
func (Str) Mask(value string, character rune, index, length int) string {
	runes := []rune(value)
	if index < 0 || index >= len(runes) {
		return value
	}
	end := index + length
	if end > len(runes) {
		end = len(runes)
	}
	for i := index; i < end; i++ {
		runes[i] = character
	}
	return string(runes)
}

// PadLeft pads a string on the left to the given length
func (Str) PadLeft(value string, length int, pad string) string {
	if pad == "" {
		pad = " "
	}
	for len([]rune(value)) < length {
		value = pad + value
	}
	return value
}

// PadRight pads a string on the right to the given length
func (Str) PadRight(value string, length int, pad string) string {
	if pad == "" {
		pad = " "
	}
	for len([]rune(value)) < length {
		value = value + pad
	}
	return value
}

// Pascal converts a string to PascalCase
func (Str) Pascal(value string) string {
	words := splitWords(value)
	for i := range words {
		words[i] = strings.Title(strings.ToLower(words[i]))
	}
	return strings.Join(words, "")
}

// Random generates a random string of the given length
func (Str) Random(length int) string {
	bytes := make([]byte, (length+1)/2)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)[:length]
}

// Remove removes all occurrences of the given value from the string
func (Str) Remove(search, subject string) string {
	return strings.ReplaceAll(subject, search, "")
}

// Replace replaces a given value in the string
func (Str) Replace(search, replace, subject string) string {
	return strings.ReplaceAll(subject, search, replace)
}

// ReplaceFirst replaces the first occurrence of a given value
func (Str) ReplaceFirst(search, replace, subject string) string {
	return strings.Replace(subject, search, replace, 1)
}

// Reverse reverses a string
func (Str) Reverse(value string) string {
	runes := []rune(value)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Slug generates a URL-friendly slug from the given string
func (Str) Slug(value string, separator string) string {
	if separator == "" {
		separator = "-"
	}
	// Convert to lowercase
	value = strings.ToLower(value)
	// Replace non-alphanumeric characters with separator
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	value = reg.ReplaceAllString(value, separator)
	// Trim separators from ends
	value = strings.Trim(value, separator)
	return value
}

// Snake converts a string to snake_case
func (Str) Snake(value string) string {
	words := splitWords(value)
	for i := range words {
		words[i] = strings.ToLower(words[i])
	}
	return strings.Join(words, "_")
}

// Start prepends a single instance of a given value to a string
func (Str) Start(value, prefix string) string {
	if !strings.HasPrefix(value, prefix) {
		return prefix + value
	}
	return value
}

// StartsWith checks if a string starts with a given prefix
func (Str) StartsWith(haystack, needle string) bool {
	return strings.HasPrefix(haystack, needle)
}

// Studly converts a string to StudlyCase (PascalCase)
func (Str) Studly(value string) string {
	return Str{}.Pascal(value)
}

// Substr returns a substring of the given string
func (Str) Substr(value string, start int, length int) string {
	runes := []rune(value)
	if start < 0 {
		start = len(runes) + start
	}
	if start < 0 {
		start = 0
	}
	if start >= len(runes) {
		return ""
	}
	end := start + length
	if length < 0 || end > len(runes) {
		end = len(runes)
	}
	return string(runes[start:end])
}

// Title converts a string to Title Case
func (Str) Title(value string) string {
	return strings.Title(strings.ToLower(value))
}

// Trim trims whitespace from both ends of a string
func (Str) Trim(value string) string {
	return strings.TrimSpace(value)
}

// TrimLeft trims whitespace from the left of a string
func (Str) TrimLeft(value string) string {
	return strings.TrimLeft(value, " \t\n\r")
}

// TrimRight trims whitespace from the right of a string
func (Str) TrimRight(value string) string {
	return strings.TrimRight(value, " \t\n\r")
}

// Upper converts a string to uppercase
func (Str) Upper(value string) string {
	return strings.ToUpper(value)
}

// Words limits the number of words in a string
func (Str) Words(value string, words int, end string) string {
	parts := strings.Fields(value)
	if len(parts) <= words {
		return value
	}
	if end == "" {
		end = "..."
	}
	return strings.Join(parts[:words], " ") + end
}

// Wrap wraps a string with given strings
func (Str) Wrap(value, before, after string) string {
	return before + value + after
}

// UUID generates a random UUID v4
func (Str) UUID() string {
	uuid := make([]byte, 16)
	rand.Read(uuid)
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10
	return hex.EncodeToString(uuid[:4]) + "-" +
		hex.EncodeToString(uuid[4:6]) + "-" +
		hex.EncodeToString(uuid[6:8]) + "-" +
		hex.EncodeToString(uuid[8:10]) + "-" +
		hex.EncodeToString(uuid[10:])
}

// splitWords splits a string into words based on various delimiters
func splitWords(s string) []string {
	// Replace common delimiters with spaces
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")

	// Handle camelCase and PascalCase
	var result []string
	var current strings.Builder

	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			if current.Len() > 0 {
				result = append(result, current.String())
				current.Reset()
			}
		}
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			current.WriteRune(r)
		} else if current.Len() > 0 {
			result = append(result, current.String())
			current.Reset()
		}
	}

	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}
