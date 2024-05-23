package cmd

import (
	"bufio"
	"fmt"
	"github.com/aakashshankar/claude-cli/api"
	"github.com/aakashshankar/claude-cli/session"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "claude",
	Short: "A CLI for Claude's completions",
	Long:  `A CLI for Claude's completions. Provide your API keys as environment variables`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			chat()
		}
	},
}

func chat() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Entering chat mode. Type 'exit' to exit")
	session.ClearSession()
	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		if text == "exit\n" {
			break
		}
		api.PromptClaude(text, true, 1024, "claude-3-sonnet-20240229", "You are a helpful assistant.", false)
		fmt.Println()
	}
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
