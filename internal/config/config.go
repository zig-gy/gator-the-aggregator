package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFilename = ".gatorconfig.json"

type Config struct {
	DbUrl string			`json:"db_url"`
	CurrentUsername string	`json:"current_user_name"`
}

func Read() (cfg Config, err error) {
	filepath, err := getConfigFilePath()
	if err != nil {
		err = fmt.Errorf("error getting config file path: %v", err)
		return
	}

	file, err := os.ReadFile(filepath)
	if err != nil {
		err = fmt.Errorf("error reading config file: %v", err)
		return
	}

	err = json.Unmarshal(file, &cfg)
	if err != nil {
		err = fmt.Errorf("error decoding json from file: %v", err)
	}
	return
}

func (cfg *Config) SetUser(user string) error {
	filepath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("error getting config file path: %v", err)
	}

	cfg.CurrentUsername = user
	json, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("error encoding json: %v", err)
	}

	if err := os.WriteFile(filepath, json, 0777); err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}
	return nil
}


func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error finding home directory: %v", err)
	}

	return homeDir + "/" + configFilename, nil
}