package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var chainCmd = &cobra.Command{
	Use:   "chain [file]",
	Short: "Visualize the variable flow chain in a flow file",
	Long: `Parse a .flow.md file and display how variables flow between steps —
which step captures a variable and which steps consume it.`,
	Example: `  # Visualize variable chain
  kest chain login.flow.md`,
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		doc, _ := ParseFlowDocument(string(content))
		if len(doc.Steps) == 0 {
			return fmt.Errorf("no steps found in %s", filePath)
		}

		printChain(doc)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(chainCmd)
}

func printChain(doc FlowDoc) {
	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7D56F4"))
	stepStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00BFFF"))
	varStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700"))
	arrowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))

	// Track: variable → producing step, consuming steps
	type varInfo struct {
		Producer  string
		Consumers []string
	}
	varMap := make(map[string]*varInfo)
	varRefRegex := regexp.MustCompile(`\{\{(\w+)\}\}`)

	// Pass 1: find producers (captures from HTTP and exec steps)
	for _, step := range doc.Steps {
		allCaptures := append(step.Request.Captures, step.Exec.Captures...)
		for _, cap := range allCaptures {
			parts := strings.SplitN(cap, "=", 2)
			if len(parts) == 2 {
				varName := strings.TrimSpace(parts[0])
				varMap[varName] = &varInfo{Producer: step.ID}
			}
		}
	}

	// Pass 2: find consumers ({{var}} references)
	for _, step := range doc.Steps {
		raw := step.Raw
		matches := varRefRegex.FindAllStringSubmatch(raw, -1)
		for _, m := range matches {
			varName := m[1]
			if info, ok := varMap[varName]; ok {
				info.Consumers = append(info.Consumers, step.ID)
			}
		}
	}

	// Print flow graph
	fmt.Printf("\n%s\n\n", title.Render(fmt.Sprintf("══ Variable Chain: %s ══", doc.Meta.Name)))

	// Print step sequence
	for i, step := range doc.Steps {
		name := step.Name
		if name == "" {
			name = step.ID
		}
		fmt.Printf("  %s", stepStyle.Render(name))
		if i < len(doc.Steps)-1 {
			fmt.Printf(" %s ", arrowStyle.Render("→"))
		}
	}
	fmt.Println()

	// Print variable flows
	if len(varMap) == 0 {
		fmt.Println("  No variable captures found.")
		return
	}

	fmt.Printf("  %s\n", title.Render("Variables:"))
	for varName, info := range varMap {
		producer := info.Producer
		if producer == "" {
			producer = "?"
		}
		consumers := "(unused)"
		if len(info.Consumers) > 0 {
			consumers = strings.Join(info.Consumers, ", ")
		}
		fmt.Printf("  %s  %s %s %s %s\n",
			varStyle.Render("{{"+varName+"}}"),
			stepStyle.Render(producer),
			arrowStyle.Render("→"),
			stepStyle.Render(consumers),
			"",
		)
	}

	// Print Mermaid diagram
	fmt.Printf("\n  %s\n", title.Render("Mermaid:"))
	fmt.Println("  " + strings.ReplaceAll(FlowToMermaid(doc), "\n", "\n  "))
}
