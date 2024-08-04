package cmd

import (
	"github.com/aakashshankar/llm-cli/session"
	"github.com/spf13/cobra"
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
	return &cobra.Command{
		Use:   "switch <uuid>",
		Short: "Switch to a session",
		Long:  "Switch to a session",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			session.SwitchSession(args[0])
		},
	}
}

func listSessionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all sessions",
		Long:  "List all sessions",
		Run: func(cmd *cobra.Command, args []string) {
			session.ListSessions()
		},
	}
}

func sessionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "session",
		Short: "Manage existing sessions",
		Long:  "Switch between existing sessions",
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			if err != nil {
				return
			}
		},
	}
}
