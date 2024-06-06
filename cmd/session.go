package cmd

import (
	"github.com/aakashshankar/llm-cli/session"
	"github.com/spf13/cobra"
)

func init() {
	sessionCmd := &cobra.Command{
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

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all sessions",
		Long:  "List all sessions",
		Run: func(cmd *cobra.Command, args []string) {
			session.ListSession()
		},
	}

	sessionCmd.AddCommand(listCmd)
	rootCmd.AddCommand(sessionCmd)
}
