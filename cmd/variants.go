package cmd

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/aakashshankar/llm-cli/api"
	"github.com/aakashshankar/llm-cli/defaults"
	"github.com/spf13/cobra"
	"os"
)

//go:embed static/variants.json
var variants []byte

var stream bool
var tokens int
var system string
var clr bool

func LoadVariants() []string {
	var vars []string
	err := json.Unmarshal(variants, &vars)
	if err != nil {
		fmt.Println("error unmarshaling variants: %w", err)
	}
	return vars
}

func init() {
	llmVariants := LoadVariants()
	for _, v := range llmVariants {
		// absolute nightmare of a pitfall
		variant := v
		llmCmd := &cobra.Command{
			Use:   variant + " <prompt>",
			Short: "Interact with the " + variant + " LLM",
			Long:  "Prompt the " + variant + " LLM, with persistent context",
			Args:  cobra.ExactArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				model, err := cmd.Flags().GetString("model")
				if err != nil {
					fmt.Println("Error getting model flag:", err)
					os.Exit(1)
				}
				api.Prompt(variant, args[0], stream, tokens, model, system, clr)
			},
		}

		defaultModel, err := defaults.GetDefaultModel(variant)
		if err != nil {
			fmt.Println("Error getting default model:", err)
			os.Exit(1)
		}

		llmCmd.Flags().String("model", defaultModel, "Model that responds to the prompt")
		llmCmd.Flags().BoolVarP(&stream, "stream", "s", false, "Stream the response")
		llmCmd.Flags().IntVarP(&tokens, "tokens", "t", 1024, "Maximum number of tokens preferred in the output")
		llmCmd.Flags().StringVarP(&system, "system", "S", "", "Set the system prompt")
		llmCmd.Flags().BoolVarP(&clr, "clear", "c", false, "Clear the current session and start a new one")

		chatCmd := &cobra.Command{
			Use:   "chat",
			Short: "Interact with the " + variant + " LLM in chat mode",
			Long:  "Chat with the " + variant + " LLM",
			Run: func(cmd *cobra.Command, args []string) {
				api.Chat(variant)
			},
		}
		llmCmd.AddCommand(chatCmd)
		rootCmd.AddCommand(llmCmd)
	}
}
