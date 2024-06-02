package defaults

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
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
