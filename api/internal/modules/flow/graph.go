package flow

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

// VariableMappingRule represents one explicit variable handoff on an edge.
type VariableMappingRule struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type graphStep struct {
	ID        uint
	ClientKey string
	SortOrder int
}

type graphEdge struct {
	SourceKey string
	TargetKey string
	Mapping   string
}

func parseVariableMappingRules(raw string) ([]VariableMappingRule, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil, nil
	}

	var rules []VariableMappingRule
	if err := json.Unmarshal([]byte(trimmed), &rules); err != nil {
		return nil, newFlowError(http.StatusUnprocessableEntity, "variable_mapping must be valid JSON")
	}

	for _, rule := range rules {
		if strings.TrimSpace(rule.Source) == "" || strings.TrimSpace(rule.Target) == "" {
			return nil, newFlowError(http.StatusUnprocessableEntity, "variable_mapping rules require non-empty source and target")
		}
	}

	return rules, nil
}

func normalizeStepClientKey(stepID uint, clientKey string) string {
	trimmed := strings.TrimSpace(clientKey)
	if trimmed != "" && !(stepID > 0 && trimmed == "step-0") {
		return trimmed
	}

	return fmt.Sprintf("step-%d", stepID)
}

func validateSaveGraph(steps []SaveStepRequest, edges []SaveEdgeRequest) error {
	graphSteps := make([]graphStep, 0, len(steps))
	for _, step := range steps {
		clientKey := strings.TrimSpace(step.ClientKey)
		if clientKey == "" {
			return newFlowError(http.StatusUnprocessableEntity, "step client_key is required")
		}
		if strings.TrimSpace(step.Name) == "" {
			return newFlowError(http.StatusUnprocessableEntity, "step name is required")
		}
		if strings.TrimSpace(step.Method) == "" {
			return newFlowError(http.StatusUnprocessableEntity, "step method is required")
		}
		if strings.TrimSpace(step.URL) == "" {
			return newFlowError(http.StatusUnprocessableEntity, "step url is required")
		}

		graphSteps = append(graphSteps, graphStep{
			ClientKey: clientKey,
			SortOrder: step.SortOrder,
		})
	}

	graphEdges := make([]graphEdge, 0, len(edges))
	for _, edge := range edges {
		sourceKey := strings.TrimSpace(edge.SourceClientKey)
		targetKey := strings.TrimSpace(edge.TargetClientKey)
		if sourceKey == "" || targetKey == "" {
			return newFlowError(http.StatusUnprocessableEntity, "edge source_client_key and target_client_key are required")
		}

		graphEdges = append(graphEdges, graphEdge{
			SourceKey: sourceKey,
			TargetKey: targetKey,
			Mapping:   edge.VariableMapping,
		})
	}

	return validateGraph(graphSteps, graphEdges)
}

func validateStoredGraph(steps []*FlowStepPO, edges []*FlowEdgePO) error {
	graphSteps := make([]graphStep, 0, len(steps))
	for _, step := range steps {
		graphSteps = append(graphSteps, graphStep{
			ID:        step.ID,
			ClientKey: normalizeStepClientKey(step.ID, step.ClientKey),
			SortOrder: step.SortOrder,
		})
	}

	stepKeysByID := make(map[uint]string, len(graphSteps))
	for _, step := range graphSteps {
		stepKeysByID[step.ID] = step.ClientKey
	}

	graphEdges := make([]graphEdge, 0, len(edges))
	for _, edge := range edges {
		graphEdges = append(graphEdges, graphEdge{
			SourceKey: stepKeysByID[edge.SourceStepID],
			TargetKey: stepKeysByID[edge.TargetStepID],
			Mapping:   edge.VariableMapping,
		})
	}

	return validateGraph(graphSteps, graphEdges)
}

func validateGraph(steps []graphStep, edges []graphEdge) error {
	if len(steps) == 0 {
		return nil
	}

	stepIndex := make(map[string]graphStep, len(steps))
	for _, step := range steps {
		if step.ClientKey == "" {
			return newFlowError(http.StatusUnprocessableEntity, "step client_key is required")
		}
		if _, exists := stepIndex[step.ClientKey]; exists {
			return newFlowError(http.StatusUnprocessableEntity, fmt.Sprintf("duplicate step client_key %q", step.ClientKey))
		}
		stepIndex[step.ClientKey] = step
	}

	adjacency := make(map[string][]string, len(steps))
	inDegree := make(map[string]int, len(steps))
	for _, step := range steps {
		inDegree[step.ClientKey] = 0
	}

	for _, edge := range edges {
		if edge.SourceKey == "" || edge.TargetKey == "" {
			return newFlowError(http.StatusUnprocessableEntity, "edge references must resolve to existing steps")
		}
		if edge.SourceKey == edge.TargetKey {
			return newFlowError(http.StatusUnprocessableEntity, "flow edges cannot form self loops")
		}
		if _, ok := stepIndex[edge.SourceKey]; !ok {
			return newFlowError(http.StatusUnprocessableEntity, fmt.Sprintf("edge source step %q does not exist", edge.SourceKey))
		}
		if _, ok := stepIndex[edge.TargetKey]; !ok {
			return newFlowError(http.StatusUnprocessableEntity, fmt.Sprintf("edge target step %q does not exist", edge.TargetKey))
		}
		if _, err := parseVariableMappingRules(edge.Mapping); err != nil {
			return err
		}

		adjacency[edge.SourceKey] = append(adjacency[edge.SourceKey], edge.TargetKey)
		inDegree[edge.TargetKey]++
	}

	queue := make([]string, 0)
	for _, step := range steps {
		if inDegree[step.ClientKey] == 0 {
			queue = append(queue, step.ClientKey)
		}
	}
	sort.Slice(queue, func(i, j int) bool {
		return stepIndex[queue[i]].SortOrder < stepIndex[queue[j]].SortOrder
	})

	visited := 0
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		visited++

		nextSteps := adjacency[current]
		sort.Slice(nextSteps, func(i, j int) bool {
			return stepIndex[nextSteps[i]].SortOrder < stepIndex[nextSteps[j]].SortOrder
		})
		for _, next := range nextSteps {
			inDegree[next]--
			if inDegree[next] == 0 {
				queue = append(queue, next)
				sort.Slice(queue, func(i, j int) bool {
					return stepIndex[queue[i]].SortOrder < stepIndex[queue[j]].SortOrder
				})
			}
		}
	}

	if visited != len(steps) {
		return newFlowError(http.StatusUnprocessableEntity, "flow edges cannot form cycles")
	}

	return nil
}
