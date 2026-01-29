package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Version      int                    `yaml:"version"`
	Defaults     Defaults               `yaml:"defaults"`
	Environments map[string]Environment `yaml:"environments"`
	ActiveEnv    string                 `yaml:"active_env"`
	ProjectID    string                 // Relative or absolute path used as ID
	ProjectPath  string                 // Absolute path to the project root
}

type Defaults struct {
	Timeout int               `yaml:"timeout"`
	Headers map[string]string `yaml:"headers"`
}

type Environment struct {
	BaseURL   string            `yaml:"base_url"`
	Variables map[string]string `yaml:"variables"`
}

func LoadConfig() (*Config, error) {
	viper.Reset()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Global config
	home, _ := os.UserHomeDir()
	viper.AddConfigPath(filepath.Join(home, ".kest"))

	// Project detection
	projectRoot, _ := findProjectRoot()
	if projectRoot != "" {
		viper.AddConfigPath(filepath.Join(projectRoot, ".kest"))
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var conf Config
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, err
	}

	conf.ProjectPath = projectRoot
	if projectRoot != "" {
		conf.ProjectID = filepath.Base(projectRoot)
	}

	return &conf, nil
}

func findProjectRoot() (string, error) {
	curr, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(curr, ".kest")); err == nil {
			return curr, nil
		}

		parent := filepath.Dir(curr)
		if parent == curr {
			break
		}
		curr = parent
	}
	return "", nil
}

func (c *Config) GetActiveEnv() Environment {
	if env, ok := c.Environments[c.ActiveEnv]; ok {
		return env
	}
	return Environment{}
}
