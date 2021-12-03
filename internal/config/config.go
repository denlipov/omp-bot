package config

import (
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

// Build information -ldflags .
const (
	version    string = "dev"
	commitHash string = "-"
)

var (
	cfg *Config
	mx  *sync.Mutex = new(sync.Mutex)
)

// GetConfigInstance returns service config
func GetConfig() Config {
	mx.Lock()
	defer mx.Unlock()
	if cfg != nil {
		return *cfg
	}

	return Config{}
}

// Grpc - contains parameter address grpc.
type Grpc struct {
	ComRequestApiURI   string `yaml:"com_request_api_uri"`
	ServiceConnTimeout int64  `yaml:"service_conn_timeout"`
}

// Project - contains all parameters project information.
type Project struct {
	Debug       bool   `yaml:"debug"`
	Name        string `yaml:"name"`
	Environment string `yaml:"environment"`
	Version     string
	CommitHash  string
}

type Bot struct {
	Token         string `yaml:"token"`
	UpdateTimeout int    `yaml:"update_timeout"`
}

// Config - contains all configuration parameters in config package.
type Config struct {
	Project Project `yaml:"project"`
	Bot     Bot     `yaml:"bot"`
	Grpc    Grpc    `yaml:"grpc"`
}

// ReadConfigYML - read configurations from file and init instance Config.
func ReadConfigYML(filePath string) error {
	if cfg != nil {
		return nil
	}

	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return err
	}

	cfg.Project.Version = version
	cfg.Project.CommitHash = commitHash

	return nil
}
