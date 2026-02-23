package config

import (
	"encoding/json"
	"os"
)

func ReadConfig() (Config, error) {

	configFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	NewConfig := Config{}
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return Config{}, err
	}

	err = json.Unmarshal(data, &NewConfig)
	if err != nil {
		return Config{}, err
	}

	return NewConfig, nil

}

func (cfg *Config) SetUser(user string) error {
	cfg.CurrentUserName = user

	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	err = os.WriteFile(configFilePath, jsonData, 0666)
	if err != nil {
		return err
	}

	return nil

}

func getConfigFilePath() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	filepath := homedir + "/" + filename

	return filepath, nil
}
