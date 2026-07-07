package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Workspaces []Workspace `yaml:"workspaces"`
}

type Workspace struct {
	Path         string       `yaml:"path"`
	Repositories []Repository `yaml:"repositories"`
}

type Repository struct {
	URL string `yaml:"url"`
	Dir string `yaml:"dir"`
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
