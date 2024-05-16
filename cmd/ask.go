package cmd

import (
	"github.com/aakashshankar/claude-cli/api"
	"github.com/spf13/cobra"
)

var stream bool
var askCmd = &cobra.Command{
	Use:   "ask <prompt>",
	Short: "Prompt claude",
	Long:  `This command prompts the Claude API with the given prompt and returns the response`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		api.PromptClaude(args[0], stream)
	},
}

func init() {
	askCmd.PersistentFlags().BoolVarP(&stream, "stream", "s", false, "Stream the response")
	rootCmd.AddCommand(askCmd)
}
