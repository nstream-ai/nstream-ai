package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type UserConfig struct {
	Email     string `json:"email"`
	OrgName   string `json:"org_name"`
	Role      string `json:"role"`
	AuthToken string `json:"auth_token"`
}

type ClusterConfig struct {
	Name          string `json:"name"`
	CloudProvider string `json:"cloud_provider"`
	Region        string `json:"region"`
	Bucket        string `json:"bucket"`
	Role          string `json:"role"`
	ClusterToken  string `json:"cluster_token"`
}

type Config struct {
	User    UserConfig    `json:"user"`
	Cluster ClusterConfig `json:"cluster"`
}

func GetConfigPath() string {
	return filepath.Join(os.Getenv("HOME"), ".nstreamconfig")
}

func LoadConfig() (*Config, error) {
	configPath := GetConfigPath()
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func SaveConfig(config *Config) error {
	configPath := GetConfigPath()
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0600)
}

func ConfigExists() bool {
	configPath := GetConfigPath()
	_, err := os.Stat(configPath)
	return !os.IsNotExist(err)
}
