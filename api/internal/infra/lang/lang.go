package lang

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Translator handles i18n translations
// Translator with JSON file support for internationalization.
type Translator struct {
	mu           sync.RWMutex
	locale       string
	fallback     string
	translations map[string]map[string]string // locale -> key -> value
	paths        []string
	loaded       map[string]bool
}

var (
	translator *Translator
	once       sync.Once
)

// Global returns the global translator instance
func Global() *Translator {
	once.Do(func() {
		translator = New("en", "en")
	})
	return translator
}

// New creates a new translator
func New(locale, fallback string) *Translator {
	return &Translator{
		locale:       locale,
		fallback:     fallback,
		translations: make(map[string]map[string]string),
		paths:        make([]string, 0),
		loaded:       make(map[string]bool),
	}
}

// SetLocale sets the current locale
func (t *Translator) SetLocale(locale string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.locale = locale
}

// GetLocale returns the current locale
func (t *Translator) GetLocale() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.locale
}

// SetFallback sets the fallback locale
func (t *Translator) SetFallback(fallback string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.fallback = fallback
}

// AddPath adds a translation files path
func (t *Translator) AddPath(path string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.paths = append(t.paths, path)
}

// Load loads translations for a locale from all registered paths
func (t *Translator) Load(locale string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.loaded[locale] {
		return nil
	}

	if t.translations[locale] == nil {
		t.translations[locale] = make(map[string]string)
	}

	for _, basePath := range t.paths {
		// Load JSON files: lang/en.json or lang/en/*.json
		jsonFile := filepath.Join(basePath, locale+".json")
		if err := t.loadJSONFile(locale, jsonFile); err != nil && !os.IsNotExist(err) {
			return err
		}

		// Load from directory
		dirPath := filepath.Join(basePath, locale)
		if info, err := os.Stat(dirPath); err == nil && info.IsDir() {
			files, _ := filepath.Glob(filepath.Join(dirPath, "*.json"))
			for _, file := range files {
				if err := t.loadJSONFile(locale, file); err != nil {
					return err
				}
			}
		}
	}

	t.loaded[locale] = true
	return nil
}

// loadJSONFile loads translations from a JSON file
func (t *Translator) loadJSONFile(locale, filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var translations map[string]interface{}
	if err := json.Unmarshal(data, &translations); err != nil {
		return fmt.Errorf("failed to parse %s: %w", filename, err)
	}

	// Flatten nested translations
	t.flattenTranslations(locale, "", translations)
	return nil
}

// flattenTranslations flattens nested translation maps
func (t *Translator) flattenTranslations(locale, prefix string, data map[string]interface{}) {
	for key, value := range data {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		switch v := value.(type) {
		case string:
			t.translations[locale][fullKey] = v
		case map[string]interface{}:
			t.flattenTranslations(locale, fullKey, v)
		}
	}
}

// Get retrieves a translation with optional replacements
func (t *Translator) Get(key string, replacements ...map[string]string) string {
	t.mu.RLock()
	locale := t.locale
	fallback := t.fallback
	t.mu.RUnlock()

	// Try current locale
	if value := t.getForLocale(locale, key); value != "" {
		return t.replace(value, replacements...)
	}

	// Try fallback
	if locale != fallback {
		if value := t.getForLocale(fallback, key); value != "" {
			return t.replace(value, replacements...)
		}
	}

	// Return key as fallback
	return key
}

// getForLocale retrieves a translation for a specific locale
func (t *Translator) getForLocale(locale, key string) string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if translations, ok := t.translations[locale]; ok {
		if value, ok := translations[key]; ok {
			return value
		}
	}
	return ""
}

// replace applies replacements to a translation string
func (t *Translator) replace(value string, replacements ...map[string]string) string {
	if len(replacements) == 0 {
		return value
	}

	result := value
	for _, repl := range replacements {
		for k, v := range repl {
			result = strings.ReplaceAll(result, ":"+k, v)
		}
	}
	return result
}

// Has checks if a translation exists
func (t *Translator) Has(key string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if translations, ok := t.translations[t.locale]; ok {
		if _, ok := translations[key]; ok {
			return true
		}
	}
	return false
}

// Choice handles pluralization
func (t *Translator) Choice(key string, count int, replacements ...map[string]string) string {
	value := t.Get(key, replacements...)

	// Simple pluralization: "one|many" format
	parts := strings.Split(value, "|")
	if len(parts) == 2 {
		if count == 1 {
			return t.replace(parts[0], replacements...)
		}
		return t.replace(parts[1], replacements...)
	}

	// More complex: "{0} none|{1} one|[2,*] many"
	// TODO: Implement more complex pluralization rules

	return value
}

// Add adds a translation programmatically
func (t *Translator) Add(locale, key, value string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.translations[locale] == nil {
		t.translations[locale] = make(map[string]string)
	}
	t.translations[locale][key] = value
}

// AddMany adds multiple translations
func (t *Translator) AddMany(locale string, translations map[string]string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.translations[locale] == nil {
		t.translations[locale] = make(map[string]string)
	}
	for k, v := range translations {
		t.translations[locale][k] = v
	}
}

// --- Convenience functions ---

// SetLocale sets the locale on the global translator
func SetLocale(locale string) {
	Global().SetLocale(locale)
}

// GetLocale gets the locale from the global translator
func GetLocale() string {
	return Global().GetLocale()
}

// AddPath adds a path to the global translator
func AddPath(path string) {
	Global().AddPath(path)
}

// Load loads translations for a locale
func Load(locale string) error {
	return Global().Load(locale)
}

// Get retrieves a translation
func Get(key string, replacements ...map[string]string) string {
	return Global().Get(key, replacements...)
}

// Trans is an alias for Get
func Trans(key string, replacements ...map[string]string) string {
	return Get(key, replacements...)
}

// __ is a shorthand for Get translation.
func __(key string, replacements ...map[string]string) string {
	return Get(key, replacements...)
}

// Has checks if a translation exists
func Has(key string) bool {
	return Global().Has(key)
}

// Choice handles pluralization
func Choice(key string, count int, replacements ...map[string]string) string {
	return Global().Choice(key, count, replacements...)
}
