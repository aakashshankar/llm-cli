package cmd

import (
	"fmt"
	"github.com/aakashshankar/claude-cli/api"
	"github.com/spf13/cobra"
)

var askCmd = &cobra.Command{
	Use:   "ask <prompt>",
	Short: "Prompt claude",
	Long:  `This command prompts the Claude API with the given prompt and returns the response`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(api.PromptClaude(args[0]))
	},
}

func init() {
	rootCmd.AddCommand(askCmd)
}
