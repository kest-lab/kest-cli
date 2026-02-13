package support

// Keys returns all keys from a map
func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Values returns all values from a map
func Values[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// Has checks if a key exists in the map
func Has[K comparable, V any](m map[K]V, key K) bool {
	_, ok := m[key]
	return ok
}

// Get gets a value from a map with a default value
func Get[K comparable, V any](m map[K]V, key K, defaultValue V) V {
	if v, ok := m[key]; ok {
		return v
	}
	return defaultValue
}

// Set sets a value in the map and returns the map for chaining
func Set[K comparable, V any](m map[K]V, key K, value V) map[K]V {
	m[key] = value
	return m
}

// Forget removes a key from the map and returns the map
func Forget[K comparable, V any](m map[K]V, keys ...K) map[K]V {
	for _, k := range keys {
		delete(m, k)
	}
	return m
}

// Only returns a map with only the specified keys
func Only[K comparable, V any](m map[K]V, keys ...K) map[K]V {
	result := make(map[K]V)
	for _, k := range keys {
		if v, ok := m[k]; ok {
			result[k] = v
		}
	}
	return result
}

// Except returns a map without the specified keys
func Except[K comparable, V any](m map[K]V, keys ...K) map[K]V {
	keySet := make(map[K]bool)
	for _, k := range keys {
		keySet[k] = true
	}

	result := make(map[K]V)
	for k, v := range m {
		if !keySet[k] {
			result[k] = v
		}
	}
	return result
}

// Merge merges multiple maps into one (later maps override earlier ones)
func Merge[K comparable, V any](maps ...map[K]V) map[K]V {
	result := make(map[K]V)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// MapMap transforms map values using a mapper function
func MapMap[K comparable, V any, R any](m map[K]V, mapper func(K, V) R) map[K]R {
	result := make(map[K]R)
	for k, v := range m {
		result[k] = mapper(k, v)
	}
	return result
}

// FilterMap filters a map based on a predicate
func FilterMap[K comparable, V any](m map[K]V, predicate func(K, V) bool) map[K]V {
	result := make(map[K]V)
	for k, v := range m {
		if predicate(k, v) {
			result[k] = v
		}
	}
	return result
}

// Flip swaps keys and values (values must be comparable)
func Flip[K comparable, V comparable](m map[K]V) map[V]K {
	result := make(map[V]K)
	for k, v := range m {
		result[v] = k
	}
	return result
}

// Pull gets a value and removes it from the map
func Pull[K comparable, V any](m map[K]V, key K, defaultValue V) V {
	if v, ok := m[key]; ok {
		delete(m, key)
		return v
	}
	return defaultValue
}

// IsEmpty checks if a map is empty
func IsEmptyMap[K comparable, V any](m map[K]V) bool {
	return len(m) == 0
}

// IsNotEmptyMap checks if a map is not empty
func IsNotEmptyMap[K comparable, V any](m map[K]V) bool {
	return len(m) > 0
}

// CountMap returns the number of items in a map
func CountMap[K comparable, V any](m map[K]V) int {
	return len(m)
}

// Dot flattens a nested map into dot notation
func Dot(m map[string]interface{}, prefix string) map[string]interface{} {
	result := make(map[string]interface{})
	dotHelper(m, prefix, result)
	return result
}

func dotHelper(m map[string]interface{}, prefix string, result map[string]interface{}) {
	for k, v := range m {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}

		if nested, ok := v.(map[string]interface{}); ok {
			dotHelper(nested, key, result)
		} else {
			result[key] = v
		}
	}
}

// Undot expands a dot-notated map into a nested map
func Undot(m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for k, v := range m {
		parts := splitDotKey(k)
		setNested(result, parts, v)
	}

	return result
}

func splitDotKey(key string) []string {
	var parts []string
	current := ""
	for _, c := range key {
		if c == '.' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(c)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}

func setNested(m map[string]interface{}, keys []string, value interface{}) {
	for i := 0; i < len(keys)-1; i++ {
		key := keys[i]
		if _, ok := m[key]; !ok {
			m[key] = make(map[string]interface{})
		}
		m = m[key].(map[string]interface{})
	}
	m[keys[len(keys)-1]] = value
}
