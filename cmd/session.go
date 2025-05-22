package cmd

import (
	"fmt"
	"github.com/ahhcash/llm-cli/session"
	"github.com/ahhcash/llm-cli/ui" // Import the ui package
	"github.com/spf13/cobra"
	"os"
)

func init() {
	sessionCmd := sessionCommand()
	listCmd := listSessionCommand()
	inspectCmd := inspectSessionCommand()
	switchCmd := switchSessionCommand()

	sessionCmd.AddCommand(listCmd)
	sessionCmd.AddCommand(inspectCmd)
	sessionCmd.AddCommand(switchCmd)

	rootCmd.AddCommand(sessionCmd)
}

func inspectSessionCommand() *cobra.Command {
	var yaml bool

	inspectCmd := &cobra.Command{
		Use:   "inspect <uuid>",
		Short: "Inspect a session",
		Long:  "Inspect a session",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			session.InspectSession(args[0], yaml)
		},
	}

	inspectCmd.Flags().BoolVarP(&yaml, "yaml", "y", false, "Output in YAML format")

	return inspectCmd
}

func switchSessionCommand() *cobra.Command {
	switchCmd := &cobra.Command{
		Use:   "switch [uuid]", // Changed to indicate optional UUID
		Short: "Switch to a session, or list sessions if no UUID provided",
		Long:  "Switch to a specific session by UUID. If no UUID is provided, a TUI will be shown to select a session.",
		Args:  cobra.MaximumNArgs(1), // Allow 0 or 1 argument
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				// No UUID provided, launch TUI
				ui.StartSessionListTUI()
				return nil
			}
			// UUID provided, switch directly
			session.SwitchSession(args[0])
			// Add a confirmation message, as the TUI does
			fmt.Printf("Switched to session: %s\n", args[0]) 
			return nil
		},
	}
	return switchCmd
}

func listSessionCommand() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all sessions using an interactive TUI",
		Long:  "Displays an interactive TUI to list all available sessions.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ui.StartSessionListTUI()
			return nil
		},
	}
	return listCmd
}

func sessionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "session",
		Short: "Manage existing sessions",
		Long:  "Switch between existing sessions",
		RunE: func(cmd *cobra.Command, args []string) error {
			// If 'llm session' is called without subcommands, show help.
			if len(args) == 0 { // Should not happen if subcommands are defined, but good practice.
				return cmd.Help()
			}
			// If a subcommand was intended but mistyped, Cobra handles it.
			// Default to help if we reach here somehow.
			return cmd.Help()
		},
	}
	// Ensure that if 'llm session' is run with no further args, it shows help.
	sessionCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if cmd.Use == "session" && len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
	}
	return sessionCmd
}
