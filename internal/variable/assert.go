package variable

import (
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
)

// Assert checks if the response body matches the assertion expression (e.g. status=200, body.id=1)
func Assert(status int, body []byte, assertion string) (bool, string) {
	parts := strings.SplitN(assertion, "=", 2)
	if len(parts) != 2 {
		return false, fmt.Sprintf("invalid assertion format: %s", assertion)
	}

	key := strings.TrimSpace(parts[0])
	expected := strings.TrimSpace(parts[1])

	if key == "status" {
		actual := fmt.Sprintf("%d", status)
		if actual == expected {
			return true, ""
		}
		return false, fmt.Sprintf("status mismatch: expected %s, got %s", expected, actual)
	}

	if strings.HasPrefix(key, "body.") {
		query := key[5:]
		result := gjson.Get(string(body), query)
		if !result.Exists() {
			return false, fmt.Sprintf("body path not found: %s", query)
		}
		actual := result.String()
		if actual == expected {
			return true, ""
		}
		return false, fmt.Sprintf("body mismatch at %s: expected %s, got %s", query, expected, actual)
	}

	return false, fmt.Sprintf("unknown assertion key: %s", key)
}
