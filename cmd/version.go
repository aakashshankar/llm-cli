package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

const Version = "0.0.1"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "LLM CLI version",
	Long:  `Prints the version number of the LLM CLI`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
