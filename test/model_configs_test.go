package test

//import (
//	"encoding/json"
//	"github.com/aakashshankar/llm-cli/defaults"
//	"github.com/stretchr/testify/assert"
//	"os"
//	"path/filepath"
//	"testing"
//)
//
//func setupTestConfig(t *testing.T) string {
//	tmpDir := t.TempDir()
//
//	testConfig := defaults.ModelConfig{
//		DefaultModels: map[string]string{
//			"claude":  "claude-3-sonnet-20240229",
//			"mistral": "mistral-large-latest",
//		},
//	}
//	configData, err := json.MarshalIndent(testConfig, "", "  ")
//	if err != nil {
//		t.Fatalf("Failed to marshal test config: %v", err)
//	}
//
//	err = os.MkdirAll(filepath.Join(tmpDir, "..defaults", "static"), 0755)
//	if err != nil {
//		t.Fatalf("Failed to create test config directory: %v", err)
//	}
//	configFilePath := filepath.Join(tmpDir, "..defaults", "static", "model_configs.json")
//	err = os.WriteFile(configFilePath, configData, 0644)
//	if err != nil {
//		t.Fatalf("Failed to write test config file: %v", err)
//	}
//
//	return configFilePath
//}
//
//
//func TestSetDefaultModel(t *testing.T) {
//	ass := assert.New(t)
//
//	configFilePath := setupTestConfig(t)
//	defaults.LoadConfig(configFilePath)
//
//	err := defaults.SetDefaultModel("claude", "claude-3-haiku-42069")
//	ass.NoError(err)
//}
//
//func TestDefaultModels(t *testing.T) {
//	ass := assert.New(t)
//
//	//configFilePath := setupTestConfig(t)
//	//defaults.LoadConfig(configFilePath)
//
//	claude, _ := defaults.GetDefaultModel("claude")
//	mistral, _ := defaults.GetDefaultModel("mistral")
//	ass.Equal(claude, "claude-3-sonnet-20240229")
//	ass.Equal(mistral, "mistral-large-latest")
//
//	_, err := defaults.GetDefaultModel("unknown")
//	ass.Error(err)
//}
