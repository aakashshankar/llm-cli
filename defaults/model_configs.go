package defaults

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type ModelConfig struct {
	DefaultModels map[string]string `json:"default_models"`
}

var configFilePath string
var modelConfig ModelConfig

func init() {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error determining executable path:", err)
		os.Exit(1)
	}
	configFilePath = filepath.Join(filepath.Dir(exePath), "../defaults", "static", "model_configs.json")
	LoadConfig(configFilePath)
}

func LoadConfig(configFilePath string) {
	bytes, err := os.ReadFile(configFilePath)
	if err != nil {
		fmt.Println("Error reading config file:", err)
		os.Exit(1)
	}
	err = json.Unmarshal(bytes, &modelConfig)
	if err != nil {
		fmt.Println("Error unmarshalling config file:", err)
		os.Exit(1)
	}
}

func GetDefaultModel(variant string) (string, error) {
	defaultModels, ok := modelConfig.DefaultModels[variant]
	if !ok {
		return "", fmt.Errorf("unknown LLM type: %s", variant)
	}
	return defaultModels, nil
}

func SetDefaultModel(variant string, model string) error {
	modelConfig.DefaultModels[variant] = model
	// write model config to disk
	updatedConfig, err := json.MarshalIndent(modelConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling model config to JSON: %w", err)
	}

	err = os.WriteFile(configFilePath, updatedConfig, 0644)
	if err != nil {
		return fmt.Errorf("error writing updated model config to disk: %w", err)
	}
	return nil
}
