package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage environments",
}

var envListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all environments",
	Run: func(cmd *cobra.Command, args []string) {
		// This is a bit tricky because viper doesn't easily give us the raw map
		// for now we'll just show what's in the config
		fmt.Println("Available environments:")
		envs := viper.GetStringMap("environments")
		active := viper.GetString("active_env")
		for name := range envs {
			if name == active {
				fmt.Printf("* %s\n", name)
			} else {
				fmt.Printf("  %s\n", name)
			}
		}
	},
}

var envUseCmd = &cobra.Command{
	Use:   "use [env]",
	Short: "Switch to a different environment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		newEnv := args[0]
		viper.Set("active_env", newEnv)

		// If project config exists, update it
		if _, err := os.Stat(".kest/config.yaml"); err == nil {
			// We need to write it back. Viper's WriteConfig might overwrite comments but it's okay for MVP.
			if err := viper.WriteConfigAs(".kest/config.yaml"); err != nil {
				return err
			}
			fmt.Printf("✓ Switched to environment: %s\n", newEnv)
		} else {
			fmt.Println("No project config found in .kest/config.yaml")
		}
		return nil
	},
}

func init() {
	envCmd.AddCommand(envListCmd)
	envCmd.AddCommand(envUseCmd)
	rootCmd.AddCommand(envCmd)
}
