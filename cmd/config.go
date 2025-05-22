package cmd

import (
	"github.com/ahhcash/llm-cli/ui" // Import the ui package
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Display LLM CLI configuration settings",
	Long: `Displays the current configuration settings for the LLM CLI.
This includes the status of API keys, the default completer,
and default models for each variant. It also provides instructions
on how to set these configurations, primarily via environment variables.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ui.StartConfigTUI()
		return nil
	},
}

func init() {
	// This is where we would add the configCmd to the rootCmd or another parent command.
	// Since rootCmd is in the same package, we can do it directly here or in root.go.
	// For clarity and to follow the pattern of other commands,
	// it's often added in the init() of the file defining rootCmd or in a central place.
	// Let's assume for now it will be added in root.go's init() or similar.
	// If this file's init() is guaranteed to run after root.go's init where rootCmd is defined,
	// we could do rootCmd.AddCommand(configCmd).
	// However, a safer approach is to have a function that root.go calls, or add it in root.go itself.
	// For this task, I will add it in root.go to ensure rootCmd is initialized.
}
