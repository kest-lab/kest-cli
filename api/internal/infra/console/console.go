package console

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/fatih/color"
)

// Command represents a CLI command
type Command interface {
	Name() string
	Description() string
	Usage() string
	Run(args []string) error
}

// Application is the CLI application container
type Application struct {
	name     string
	version  string
	commands map[string]Command
	groups   map[string][]string
}

// New creates a new CLI application
func New(name, version string) *Application {
	return &Application{
		name:     name,
		version:  version,
		commands: make(map[string]Command),
		groups:   make(map[string][]string),
	}
}

// RegisterAs adds a command to the application with a specific name (alias)
func (app *Application) RegisterAs(name string, cmd Command) {
	app.commands[name] = cmd

	// Extract group from command name (e.g., "make:model" -> "make")
	if idx := strings.Index(name, ":"); idx > 0 {
		group := name[:idx]
		app.groups[group] = append(app.groups[group], name)
	} else {
		app.groups[""] = append(app.groups[""], name)
	}
}

// Register adds a command to the application
func (app *Application) Register(cmd Command) {
	app.RegisterAs(cmd.Name(), cmd)
}

// Run executes the CLI application
func (app *Application) Run(args []string) error {
	if len(args) < 2 {
		app.showHelp()
		return nil
	}

	cmdName := args[1]

	switch cmdName {
	case "help", "-h", "--help":
		app.showHelp()
		return nil
	case "version", "-v", "--version":
		app.showVersion()
		return nil
	case "list":
		app.listCommands()
		return nil
	}

	cmd, ok := app.commands[cmdName]
	if !ok {
		return fmt.Errorf("command '%s' not found. Run '%s list' for available commands", cmdName, app.name)
	}

	startTime := time.Now()
	err := cmd.Run(args[2:])
	elapsed := time.Since(startTime)

	if err != nil {
		color.Red("  ✗ Error: %v", err)
		return err
	}

	color.Green("  ✓ Done in %v", elapsed.Round(time.Millisecond))
	return nil
}

func (app *Application) showVersion() {
	fmt.Printf("%s version %s\n", app.name, app.version)
}

func (app *Application) showHelp() {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	fmt.Println()
	fmt.Printf("%s %s\n\n", green(app.name), app.version)
	fmt.Printf("%s\n", yellow("Usage:"))
	fmt.Printf("  %s <command> [arguments]\n\n", app.name)
	fmt.Printf("%s\n", yellow("Available commands:"))

	app.listCommands()
}

func (app *Application) listCommands() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	// Get sorted group names
	groups := make([]string, 0, len(app.groups))
	for g := range app.groups {
		groups = append(groups, g)
	}
	sort.Strings(groups)

	for _, group := range groups {
		cmds := app.groups[group]
		sort.Strings(cmds)

		if group != "" {
			fmt.Fprintf(w, " %s\n", yellow(group))
		}

		for _, cmdName := range cmds {
			cmd := app.commands[cmdName]
			fmt.Fprintf(w, "  %s\t%s\n", green(cmdName), cmd.Description())
		}
	}

	w.Flush()
	fmt.Println()
}
