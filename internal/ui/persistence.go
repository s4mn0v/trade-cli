package ui

import (
	"encoding/json"
	"os"
)

type SavedConfig struct {
	APIKey     string `json:"api_key"`
	SecretKey  string `json:"secret_key"`
	Passphrase string `json:"passphrase"`
}

// GetConfigPath returns the path to ~/.bitget-trade-cli.json
func GetConfigPath() string {
	// home, _ := os.UserHomeDir()
	// return filepath.Join(home, ".bitget-trade-cli.json")
	return ".bitget-trade-cli"
}

// SaveSession writes the credentials to a hidden file
func SaveSession(config SavedConfig) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(GetConfigPath(), data, 0o600) // Read/Write for owner only (Secure)
}

// LoadSession reads the credentials from the hidden file
func LoadSession() (SavedConfig, error) {
	var config SavedConfig
	path := GetConfigPath()

	// Check if file exists to avoid unnecessary errors
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return config, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)
	return config, err
}
