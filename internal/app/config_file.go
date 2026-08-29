package app

import (
	"encoding/json"
	"os"
)

const appConfigSchemaVersion = 1

func appConfigFile() string {
	return dataPath("config.json")
}

func (a *App) readAppConfig() (AppConfig, error) {
	data, err := os.ReadFile(appConfigFile())
	if err != nil {
		if os.IsNotExist(err) {
			return AppConfig{}, nil
		}
		return AppConfig{}, err
	}
	var config AppConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return AppConfig{}, err
	}
	if config.SchemaVersion == 0 {
		config.SchemaVersion = appConfigSchemaVersion
	}
	return config, nil
}

func (a *App) saveAppConfig(config AppConfig) error {
	config.SchemaVersion = appConfigSchemaVersion
	return a.writeJSONFile(appConfigFile(), config)
}
