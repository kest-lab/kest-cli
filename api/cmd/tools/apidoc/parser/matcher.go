package parser

import (
	"strings"
)

// MatchRoutesWithDTOs matches routes with their corresponding DTOs
func MatchRoutesWithDTOs(routes []Route, dtos map[string]*DTO) []Endpoint {
	var endpoints []Endpoint

	for _, route := range routes {
		endpoint := Endpoint{
			Route:   route,
			Tags:    route.Module,
			Summary: generateSummary(route),
		}

		// Try to find matching request DTO
		endpoint.Request = findRequestDTO(route, dtos)

		// Try to find matching response DTO
		endpoint.Response = findResponseDTO(route, dtos)

		endpoints = append(endpoints, endpoint)
	}

	return endpoints
}

// generateSummary generates a summary from route information
func generateSummary(route Route) string {
	// Extract action from handler name
	parts := strings.Split(route.Handler, ".")
	if len(parts) < 2 {
		return route.Path
	}

	action := parts[1]
	module := route.Module

	// Generate human-readable summary
	switch action {
	case "Register":
		return "Register new " + module
	case "Login":
		return module + " login"
	case "Create", "Store":
		return "Create " + module
	case "List", "Index":
		return "List " + module + "s"
	case "Get", "Show":
		return "Get " + module + " details"
	case "Update":
		return "Update " + module
	case "Delete", "Destroy":
		return "Delete " + module
	case "GetProfile":
		return "Get current user profile"
	case "UpdateProfile":
		return "Update current user profile"
	case "ChangePassword":
		return "Change password"
	case "ResetPassword":
		return "Reset password"
	case "DeleteAccount":
		return "Delete account"
	default:
		// Convert camelCase to words
		return camelToWords(action) + " " + module
	}
}

// camelToWords converts camelCase to words
func camelToWords(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune(' ')
		}
		result.WriteRune(r)
	}
	return result.String()
}

// findRequestDTO finds the request DTO for a route
func findRequestDTO(route Route, dtos map[string]*DTO) *DTO {
	// Only POST, PUT, PATCH typically have request bodies
	if route.Method != "POST" && route.Method != "PUT" && route.Method != "PATCH" {
		return nil
	}

	// Extract action from handler
	parts := strings.Split(route.Handler, ".")
	if len(parts) < 2 {
		return nil
	}
	action := parts[1]
	module := route.Module

	// Try common naming patterns
	patterns := []string{
		module + "." + capitalize(module) + action + "Request",
		module + "." + action + "Request",
		module + "." + capitalize(module) + capitalize(action) + "Request",
		module + ".CreateRequest",
		module + ".UpdateRequest",
	}

	// Special cases
	switch action {
	case "Register":
		patterns = append([]string{module + ".UserRegisterRequest"}, patterns...)
	case "Login":
		patterns = append([]string{module + ".UserLoginRequest"}, patterns...)
	case "ChangePassword":
		patterns = append([]string{module + ".UserChangePasswordRequest"}, patterns...)
	case "ResetPassword":
		patterns = append([]string{module + ".UserPasswordResetRequest"}, patterns...)
	case "UpdateProfile":
		patterns = append([]string{module + ".UserUpdateRequest"}, patterns...)
	}

	for _, pattern := range patterns {
		if dto, ok := dtos[pattern]; ok {
			return dto
		}
	}

	return nil
}

// findResponseDTO finds the response DTO for a route
func findResponseDTO(route Route, dtos map[string]*DTO) *DTO {
	parts := strings.Split(route.Handler, ".")
	if len(parts) < 2 {
		return nil
	}
	action := parts[1]
	module := route.Module

	// Try common naming patterns
	patterns := []string{
		module + "." + capitalize(module) + action + "Response",
		module + "." + action + "Response",
		module + ".Response",
		module + "." + capitalize(module),
	}

	// Special cases
	switch action {
	case "Login":
		patterns = append([]string{module + ".UserLoginResponse"}, patterns...)
	case "List":
		patterns = append([]string{module + ".ListResponse"}, patterns...)
	}

	for _, pattern := range patterns {
		if dto, ok := dtos[pattern]; ok {
			return dto
		}
	}

	return nil
}

// capitalize capitalizes the first letter
func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
