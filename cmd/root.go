package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "llm",
	Short: "A CLI tool to interact with LLMs",
	Long:  `A CLI tool to interact with LLMs. Provide your API keys as environment variables`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
