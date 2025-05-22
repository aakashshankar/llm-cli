package cmd

import (
	"fmt"
	"os"

	"github.com/ahhcash/llm-cli/assist"
	"github.com/ahhcash/llm-cli/ui" // Import the ui package
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "llm",
	Short: "A CLI tool to interact with LLMs",
	Long:  `A CLI tool to interact with LLMs. Provide your API keys as environment variables.`,
	// Args will be handled by subcommands or the RunE logic below.
	// Setting Args: cobra.ArbitraryArgs or similar could also work if we want to allow 'llm some prompt' directly.
	// For now, let's rely on subcommands for specific actions beyond the main menu or direct assist.
	// The current `assist.Assist` call seems to be for `llm <prompt>` if no other subcommand matches.
	// We need to ensure `llm` alone triggers the menu.
	Args: cobra.MaximumNArgs(1), // Keep this as per original for assist.Assist
	RunE: func(cmd *cobra.Command, args []string) error {
		// This RunE will only be executed if no subcommands are matched.
		if len(args) == 0 && !cmd.Flags().HasFlags() { // No arguments and no flags passed to `llm` itself
			selection, err := ui.StartMainMenuListTUI()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error starting main menu: %v\n", err)
				return err
			}

			switch selection {
			case "chat":
				defaultCompleter := os.Getenv("DEFAULT_COMPLETER")
				if defaultCompleter == "" {
					fmt.Println("DEFAULT_COMPLETER environment variable is not set.")
					fmt.Println("Please set it (e.g., export DEFAULT_COMPLETER=openai) or use 'llm <variant> chat'.")
					return nil
				}
				fmt.Printf("Attempting to start chat with default completer: %s\n", defaultCompleter)
				// Assuming system prompt is handled within StartChatTUI or by user flags on direct command
				ui.StartChatTUI("", defaultCompleter) 
			case "sessions":
				ui.StartSessionListTUI()
			case "config":
				ui.StartConfigTUI()
			case "quit":
				fmt.Println("Exiting.")
			default:
				// Should not happen with current menu items
				fmt.Printf("Unknown selection: %s\n", selection)
			}
			return nil
		} else if len(args) == 1 { // If one argument is passed, assume it's for assist.Assist
			// This maintains the original behavior for `llm <prompt>`
			_, err := assist.Assist(args[0])
			if err != nil {
				fmt.Println("Error:", err)
				// os.Exit(1) // In RunE, we should return errors
				return err
			}
			return nil
		}
		// If args > 1 or flags are present but no subcommand matched, Cobra will show help.
		// Or, we can explicitly call cmd.Help()
		return cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// Cobra already prints errors to stderr by default.
		// We might want to os.Exit(1) here if Execute() itself returns an error
		// not already handled by a subcommand's RunE.
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(configCmd) // configCmd is defined in cmd/config.go
	// Other commands like sessionCmd, variant commands are added in their respective files.
}
