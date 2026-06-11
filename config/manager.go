package config

import (
	"os"

	"github.com/vimek-go/pr-analizer-action/models"

	yaml "gopkg.in/yaml.v3"
)

func Load(path string) (*models.Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg models.Config
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	err = cfg.Validate()
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
