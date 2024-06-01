package cmd

import (
	"fmt"
	"github.com/aakashshankar/llm-cli/assist"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "llm",
	Short: "A CLI tool to interact with LLMs",
	Long:  `A CLI tool to interact with LLMs. Provide your API keys as environment variables`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			err := cmd.Help()
			if err != nil {
				return
			}
			os.Exit(1)
		}
		_, err := assist.Assist(args[0])
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
