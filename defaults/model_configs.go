package defaults

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed static/model_configs.json
var modelConfigs []byte

type ModelConfig struct {
	DefaultModels map[string]string `json:"default_models"`
}

var modelConfig ModelConfig

func init() {
	err := json.Unmarshal(modelConfigs, &modelConfig)
	if err != nil {
		fmt.Println("Error unmarshalling model configs:", err)
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
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("error getting executable path: %w", err)
	}
	// there's gotta be a better way to do this
	configFile := filepath.Join(filepath.Dir(exePath), "../defaults/static/model_configs.json")

	err = os.WriteFile(configFile, updatedConfig, 0644)
	if err != nil {
		return fmt.Errorf("error writing updated model config to disk: %w", err)
	}
	return nil
}
