package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var runCmd = &cobra.Command{
	Use:   "run [file]",
	Short: "Run a Kest scenario file (.kest)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runScenario(args[0])
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func runScenario(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fmt.Printf("\n--- Step %d: %s ---\n", lineNum, line)
		err := parseAndExecuteLine(line)
		if err != nil {
			fmt.Printf("Error at line %d: %v\n", lineNum, err)
			return err
		}
	}
	return scanner.Err()
}

func parseAndExecuteLine(line string) error {
	// Simple shell-like splitter (imperfect but works for basic cases)
	// For better robustness, use a real shell word splitter library
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return fmt.Errorf("invalid command format")
	}

	method := strings.ToLower(parts[0])
	url := parts[1]

	// Use a temporary flag set to parse the rest of the line
	fs := pflag.NewFlagSet("scenario", pflag.ContinueOnError)
	var data string
	var headers, queries, captures, asserts []string
	var noRec bool

	fs.StringVarP(&data, "data", "d", "", "")
	fs.StringSliceVarP(&headers, "header", "H", []string{}, "")
	fs.StringSliceVarP(&queries, "query", "q", []string{}, "")
	fs.StringSliceVarP(&captures, "capture", "c", []string{}, "")
	fs.StringSliceVarP(&asserts, "assert", "a", []string{}, "")
	fs.BoolVar(&noRec, "no-record", false, "")

	err := fs.Parse(parts[2:])
	if err != nil {
		return err
	}

	return ExecuteRequest(RequestOptions{
		Method:   method,
		URL:      url,
		Data:     data,
		Headers:  headers,
		Queries:  queries,
		Captures: captures,
		Asserts:  asserts,
		NoRecord: noRec,
	})
}
