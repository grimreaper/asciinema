package util

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"code.google.com/p/gcfg"
)

const (
	DEFAULT_API_URL = "https://asciinema.org"
)

type Config struct {
	Api struct {
		Token string
		Url   string
	}
	Record struct {
		Command string
	}
}

type ConfigLoader interface {
	LoadConfig() (*Config, error)
}

type FileConfigLoader struct{}

func (l *FileConfigLoader) LoadConfig() (*Config, error) {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		return nil, errors.New("Need $HOME")
	}

	cfgPath := filepath.Join(homeDir, ".asciinema", "config")

	cfg, err := loadConfigFile(cfgPath)
	if err != nil {
		return nil, err
	}

	if cfg.Api.Url == "" {
		cfg.Api.Url = DEFAULT_API_URL
	}

	if envApiUrl := os.Getenv("ASCIINEMA_API_URL"); envApiUrl != "" {
		cfg.Api.Url = envApiUrl
	}

	return cfg, nil
}

func loadConfigFile(path string) (*Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err = createConfigFile(path); err != nil {
			return nil, err
		}
	}

	var cfg Config
	if err := gcfg.ReadFileInto(&cfg, path); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func createConfigFile(path string) error {
	apiToken := NewUUID().String()
	contents := fmt.Sprintf("[api]\ntoken = %v\n", apiToken)
	return ioutil.WriteFile(path, []byte(contents), 0644)
}