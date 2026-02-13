package support

import (
	"reflect"
	"strconv"
	"strings"
)

// DataGet retrieves a value from a nested data structure using "dot" notation
// Example: DataGet(data, "user.profile.name")
func DataGet(target any, key string, defaultVal ...any) any {
	if target == nil {
		return getDefault(defaultVal)
	}

	if key == "" {
		return target
	}

	keys := strings.Split(key, ".")
	current := target

	for _, k := range keys {
		current = getSegment(current, k)
		if current == nil {
			return getDefault(defaultVal)
		}
	}

	return current
}

// DataSet sets a value in a nested data structure using "dot" notation
// Note: target must be a pointer to a map
func DataSet(target any, key string, value any) bool {
	if target == nil || key == "" {
		return false
	}

	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr {
		return false
	}

	v = v.Elem()
	if v.Kind() != reflect.Map {
		return false
	}

	keys := strings.Split(key, ".")
	return setNestedReflect(v, keys, value)
}

// DataHas checks if a key exists in a nested data structure using "dot" notation
func DataHas(target any, key string) bool {
	if target == nil || key == "" {
		return false
	}

	keys := strings.Split(key, ".")
	current := target

	for _, k := range keys {
		current = getSegment(current, k)
		if current == nil {
			return false
		}
	}

	return true
}

// DataFill fills a value only if it doesn't exist
func DataFill(target any, key string, value any) bool {
	if DataHas(target, key) {
		return false
	}
	return DataSet(target, key, value)
}

// DataForget removes a key from a nested data structure
func DataForget(target any, key string) bool {
	if target == nil || key == "" {
		return false
	}

	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr {
		return false
	}

	v = v.Elem()
	if v.Kind() != reflect.Map {
		return false
	}

	keys := strings.Split(key, ".")
	if len(keys) == 1 {
		v.SetMapIndex(reflect.ValueOf(key), reflect.Value{})
		return true
	}

	// Navigate to parent
	parentKey := strings.Join(keys[:len(keys)-1], ".")
	parent := DataGet(target, parentKey)
	if parent == nil {
		return false
	}

	parentV := reflect.ValueOf(parent)
	if parentV.Kind() == reflect.Map {
		lastKey := keys[len(keys)-1]
		parentV.SetMapIndex(reflect.ValueOf(lastKey), reflect.Value{})
		return true
	}

	return false
}

// getSegment retrieves a single segment from a value
func getSegment(target any, key string) any {
	if target == nil {
		return nil
	}

	v := reflect.ValueOf(target)

	// Handle pointers
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Map:
		mapVal := v.MapIndex(reflect.ValueOf(key))
		if !mapVal.IsValid() {
			return nil
		}
		return mapVal.Interface()

	case reflect.Struct:
		field := v.FieldByName(key)
		if !field.IsValid() {
			// Try case-insensitive match
			t := v.Type()
			for i := 0; i < t.NumField(); i++ {
				if strings.EqualFold(t.Field(i).Name, key) {
					field = v.Field(i)
					break
				}
			}
		}
		if field.IsValid() && field.CanInterface() {
			return field.Interface()
		}
		return nil

	case reflect.Slice, reflect.Array:
		index, err := strconv.Atoi(key)
		if err != nil {
			return nil
		}
		if index < 0 || index >= v.Len() {
			return nil
		}
		return v.Index(index).Interface()

	default:
		return nil
	}
}

// setNestedReflect sets a value in a nested map structure using reflection
func setNestedReflect(v reflect.Value, keys []string, value any) bool {
	if len(keys) == 0 {
		return false
	}

	if len(keys) == 1 {
		v.SetMapIndex(reflect.ValueOf(keys[0]), reflect.ValueOf(value))
		return true
	}

	key := keys[0]
	mapVal := v.MapIndex(reflect.ValueOf(key))

	if !mapVal.IsValid() {
		// Create nested map
		newMap := reflect.MakeMap(reflect.TypeOf(map[string]any{}))
		v.SetMapIndex(reflect.ValueOf(key), newMap)
		return setNestedReflect(newMap, keys[1:], value)
	}

	// Navigate into existing value
	if mapVal.Kind() == reflect.Interface {
		mapVal = mapVal.Elem()
	}

	if mapVal.Kind() == reflect.Map {
		return setNestedReflect(mapVal, keys[1:], value)
	}

	return false
}

// getDefault returns the default value from a variadic slice
func getDefault(defaultVal []any) any {
	if len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return nil
}
