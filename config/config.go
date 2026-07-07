package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DefaultStartPath string      `yaml:"default_start_path,omitempty"`
	Workspaces       []Workspace `yaml:"workspaces"`
}

type Workspace struct {
	Path         string       `yaml:"path"`
	Repositories []Repository `yaml:"repositories"`
}

type Repository struct {
	URL string `yaml:"url"`
	Dir string `yaml:"dir"`
}

// GetDefaultConfigPath returns the default configuration file path according to the XDG Base Directory Specification.
// It checks for $XDG_CONFIG_HOME and falls back to $HOME/.config.
func GetDefaultConfigPath() (string, error) {
	var configHome string
	// 1. Check for XDG_CONFIG_HOME environment variable.
	if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
		configHome = xdgConfigHome
	} else {
		// 2. Fallback to the user's home directory.
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configHome = filepath.Join(homeDir, ".config")
	}

	return filepath.Join(configHome, "repolift", "config.yaml"), nil
}


func Load(filepath string) (*Config, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
