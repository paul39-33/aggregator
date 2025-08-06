package config

import (
	"fmt"
	"os"
	"path/filepath"
	"encoding/json"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	Db_url	 			string `json:"db_url"`
	CurrentUserName 	string `json:"current_user_name"`
}

func getConfigFilePath()(string, error){
	dirname, err := os.UserHomeDir()
	if err != nil {
		return dirname, fmt.Errorf("Error getting user home directory!")
	}
	filepath := filepath.Join(dirname, configFileName)
	return filepath, nil
}

func Read() (Config, error) {
	newConfig := Config{}

	filepath, err := getConfigFilePath()
	if err != nil {
		return newConfig, err
	}

	file, err := os.Open(filepath)
	if err != nil {
		return newConfig, fmt.Errorf("\nError opening the source file!")
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&newConfig); err != nil {
		return newConfig, fmt.Errorf("\nError decoding")
	}
	return newConfig, nil
}

func write(cfg Config) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("Error when trying to Marshal into cfg: %v", err)
	}

	filepath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("Error getting filepath: %v", err)
	}

	return os.WriteFile(filepath, data, 0644)
}

func (cfg *Config) SetUser(username string) error {
	cfg.CurrentUserName = username
	return write(*cfg)
}



