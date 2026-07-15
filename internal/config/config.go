package config

import (
	"backup-maker/internal/validator"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
)

type Config struct {
	SourceDir        string   `json:"source_dir" validate:"required,dir_exists"`
	DestinationDir   string   `json:"destination_dir" validate:"required,dir_exists"`
	BackupOnStartup  bool     `json:"autostart"`
	MaxBackups       int      `json:"max_backups" validate:"required,min=1,max=10"`
	EnabledProviders []string `json:"enabled_providers" validate:"required,min=1"`
}

func NewDefaultConfig() *Config {
	return &Config{
		SourceDir:        "",
		DestinationDir:   "",
		BackupOnStartup:  false,
		MaxBackups:       3,
		EnabledProviders: []string{"local"},
	}
}

func SaveConfig(cfg *Config) error {
	err := validator.Instance.Struct(cfg)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", " ")
	if err != nil {
		return err
	}

	path, err := getConfigPath()
	if err != nil {
		return err
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func LoadConfig() (*Config, error) {
	path, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	cfg := NewDefaultConfig()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}

	err = json.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configFilePath := filepath.Join(home, ".config", "backup-maker", "config.json")

	dir := filepath.Dir(configFilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	return configFilePath, nil
}

func ClearConfig() error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	err = os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func BackupOnStartupActive() bool {
	cfg, err := LoadConfig()
	if err != nil {
		return false
	}

	return cfg.BackupOnStartup
}

func GetEnabledProvidersIds() []string {
	cfg, err := LoadConfig()
	if err != nil {
		return []string{}
	}

	return cfg.EnabledProviders
}

func isEnabledProvider(providerId string) bool {
	enabledProvidersIds := GetEnabledProvidersIds()
	return slices.Contains(enabledProvidersIds, providerId)
}
