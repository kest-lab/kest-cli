package console

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
)

// Output provides styled console output methods
type Output struct{}

// NewOutput creates a new Output instance
func NewOutput() *Output {
	return &Output{}
}

// Info prints an info message
func (o *Output) Info(message string, args ...interface{}) {
	color.Cyan("  ℹ "+message, args...)
}

// Success prints a success message
func (o *Output) Success(message string, args ...interface{}) {
	color.Green("  ✓ "+message, args...)
}

// Error prints an error message
func (o *Output) Error(message string, args ...interface{}) {
	color.Red("  ✗ "+message, args...)
}

// Warning prints a warning message
func (o *Output) Warning(message string, args ...interface{}) {
	color.Yellow("  ⚠ "+message, args...)
}

// Line prints a plain line
func (o *Output) Line(message string, args ...interface{}) {
	fmt.Printf("  "+message+"\n", args...)
}

// NewLine prints an empty line
func (o *Output) NewLine() {
	fmt.Println()
}

// Title prints a styled title
func (o *Output) Title(title string) {
	fmt.Println()
	color.New(color.FgWhite, color.Bold).Printf("  %s\n", title)
	fmt.Println("  " + strings.Repeat("─", len(title)+2))
}

// Section prints a section header
func (o *Output) Section(title string) {
	fmt.Println()
	color.Yellow("  %s", title)
}

// Table prints a table with headers and rows
func (o *Output) Table(headers []string, rows [][]string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	// Print headers
	headerLine := "  "
	for _, h := range headers {
		headerLine += color.New(color.FgWhite, color.Bold).Sprint(h) + "\t"
	}
	fmt.Fprintln(w, headerLine)

	// Print separator
	sepLine := "  "
	for range headers {
		sepLine += "───────\t"
	}
	fmt.Fprintln(w, sepLine)

	// Print rows
	for _, row := range rows {
		rowLine := "  "
		for _, cell := range row {
			rowLine += cell + "\t"
		}
		fmt.Fprintln(w, rowLine)
	}

	w.Flush()
}

// TwoColumn prints a formatted two-column detail line.
func (o *Output) TwoColumn(left, right string) {
	dots := strings.Repeat(".", 60-len(left)-len(right))
	fmt.Printf("  %s %s %s\n", left, color.New(color.FgHiBlack).Sprint(dots), right)
}

// Progress prints a progress indicator
func (o *Output) Progress(current, total int, message string) {
	percent := float64(current) / float64(total) * 100
	bar := strings.Repeat("█", int(percent/5)) + strings.Repeat("░", 20-int(percent/5))
	fmt.Printf("\r  %s [%s] %.0f%% ", message, color.GreenString(bar), percent)
	if current == total {
		fmt.Println()
	}
}

// Confirm prompts for yes/no confirmation
func (o *Output) Confirm(question string, defaultYes bool) bool {
	defaultHint := "[y/N]"
	if defaultYes {
		defaultHint = "[Y/n]"
	}

	fmt.Printf("  %s %s: ", question, color.New(color.FgHiBlack).Sprint(defaultHint))

	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))

	if answer == "" {
		return defaultYes
	}
	return answer == "y" || answer == "yes"
}

// Ask prompts for user input
func (o *Output) Ask(question string, defaultValue string) string {
	if defaultValue != "" {
		fmt.Printf("  %s [%s]: ", question, color.New(color.FgHiBlack).Sprint(defaultValue))
	} else {
		fmt.Printf("  %s: ", question)
	}

	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(answer)

	if answer == "" {
		return defaultValue
	}
	return answer
}

// Choice presents a list of options and returns the selected one
func (o *Output) Choice(question string, options []string, defaultIndex int) string {
	fmt.Printf("  %s:\n", question)

	for i, opt := range options {
		if i == defaultIndex {
			fmt.Printf("    %s %s %s\n",
				color.GreenString("[%d]", i),
				opt,
				color.New(color.FgHiBlack).Sprint("(default)"))
		} else {
			fmt.Printf("    [%d] %s\n", i, opt)
		}
	}

	fmt.Print("  > ")
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(answer)

	if answer == "" {
		return options[defaultIndex]
	}

	var idx int
	fmt.Sscanf(answer, "%d", &idx)
	if idx >= 0 && idx < len(options) {
		return options[idx]
	}
	return options[defaultIndex]
}
