package main

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion scripts",
	Long: `Generate shell completion scripts for Kest CLI.

To load completions:

Bash:
  $ source <(kest completion bash)
  # To load completions for each session, execute once:
  # Linux:
  $ kest completion bash > /etc/bash_completion.d/kest
  # macOS:
  $ kest completion bash > $(brew --prefix)/etc/bash_completion.d/kest

Zsh:
  $ source <(kest completion zsh)
  # To load completions for each session, execute once:
  $ kest completion zsh > "${fpath[1]}/_kest"

Fish:
  $ kest completion fish | source
  # To load completions for each session, execute once:
  $ kest completion fish > ~/.config/fish/completions/kest.fish

PowerShell:
  PS> kest completion powershell | Out-String | Invoke-Expression
`,
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			return rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		default:
			return cmd.Help()
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
