package parser

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// ChangeType represents the type of change
type ChangeType string

const (
	ChangeAdded    ChangeType = "➕ Added"
	ChangeRemoved  ChangeType = "➖ Removed"
	ChangeModified ChangeType = "📝 Modified"
)

// Change represents a detected change
type Change struct {
	Type     ChangeType
	Endpoint string
	Details  string
}

func (c Change) String() string {
	if c.Details != "" {
		return fmt.Sprintf("%s %s (%s)", c.Type, c.Endpoint, c.Details)
	}
	return fmt.Sprintf("%s %s", c.Type, c.Endpoint)
}

// DetectChanges compares current endpoints with previously generated docs
func DetectChanges(outputDir string, endpoints []Endpoint) []Change {
	var changes []Change

	// Load previous state
	stateFile := filepath.Join(outputDir, ".api-state.json")
	previousState := loadState(stateFile)

	// Build current state
	currentState := buildState(endpoints)

	// Compare states
	for key, current := range currentState {
		if previous, exists := previousState[key]; exists {
			if current != previous {
				changes = append(changes, Change{
					Type:     ChangeModified,
					Endpoint: key,
					Details:  "request/response changed",
				})
			}
		} else {
			changes = append(changes, Change{
				Type:     ChangeAdded,
				Endpoint: key,
			})
		}
	}

	for key := range previousState {
		if _, exists := currentState[key]; !exists {
			changes = append(changes, Change{
				Type:     ChangeRemoved,
				Endpoint: key,
			})
		}
	}

	// Save current state for next comparison
	saveState(stateFile, currentState)

	// Sort changes
	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Endpoint < changes[j].Endpoint
	})

	return changes
}

// buildState creates a state map from endpoints
func buildState(endpoints []Endpoint) map[string]string {
	state := make(map[string]string)

	for _, ep := range endpoints {
		key := fmt.Sprintf("%s %s", ep.Route.Method, ep.Route.Path)

		// Create hash of endpoint content
		content := struct {
			Method   string
			Path     string
			Request  *DTO
			Response *DTO
			Auth     bool
		}{
			Method:   ep.Route.Method,
			Path:     ep.Route.Path,
			Request:  ep.Request,
			Response: ep.Response,
			Auth:     !ep.Route.IsPublic,
		}

		data, _ := json.Marshal(content)
		hash := fmt.Sprintf("%x", md5.Sum(data))
		state[key] = hash
	}

	return state
}

// loadState loads previous state from file
func loadState(filename string) map[string]string {
	data, err := os.ReadFile(filename)
	if err != nil {
		return make(map[string]string)
	}

	var state map[string]string
	if err := json.Unmarshal(data, &state); err != nil {
		return make(map[string]string)
	}

	return state
}

// saveState saves current state to file
func saveState(filename string, state map[string]string) {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return
	}
	os.WriteFile(filename, data, 0644)
}

// GenerateSingleModuleDoc generates documentation for a single module
func GenerateSingleModuleDoc(endpoints []Endpoint, outputFile string, moduleName string) error {
	// Filter endpoints by module
	var moduleEndpoints []Endpoint
	for _, ep := range endpoints {
		if ep.Route.Module == moduleName {
			moduleEndpoints = append(moduleEndpoints, ep)
		}
	}

	if len(moduleEndpoints) == 0 {
		return fmt.Errorf("no endpoints found for module: %s", moduleName)
	}

	// Generate using existing function
	modules := map[string][]Endpoint{moduleName: moduleEndpoints}
	return generateModuleFile(modules, outputFile, moduleName)
}

// generateModuleFile generates a single module file
func generateModuleFile(modules map[string][]Endpoint, outputFile string, moduleName string) error {
	endpoints := modules[moduleName]

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# %s API\n\n", capitalize(moduleName)))
	sb.WriteString(fmt.Sprintf("> Generated at: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))

	// Summary table
	sb.WriteString("## Endpoints\n\n")
	sb.WriteString("| Method | Path | Summary | Auth |\n")
	sb.WriteString("|--------|------|---------|------|\n")
	for _, ep := range endpoints {
		auth := "🔓 Public"
		if !ep.Route.IsPublic {
			auth = "🔒 Required"
		}
		sb.WriteString(fmt.Sprintf("| `%s` | `%s` | %s | %s |\n",
			ep.Route.Method, ep.Route.Path, ep.Summary, auth))
	}
	sb.WriteString("\n---\n\n")

	// Detailed docs
	sb.WriteString("## Details\n\n")
	for _, ep := range endpoints {
		writeEndpoint(&sb, ep)
	}

	return os.WriteFile(outputFile, []byte(sb.String()), 0644)
}
