package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/adrg/xdg"
	"github.com/atomicptr/peon/pkg/fs"
)

type contextKey string

const ContextKey contextKey = "github.com/atomicptr/peon/pkg/config"

type Config struct {
	//
}

func Default() *Config {
	return &Config{}
}

func FromEnv() (*Config, error) {
	envPath, ok := fs.ValidatePath(os.Getenv("PEON_CONFIG"))
	if ok {
		return FromPath(envPath)
	}

	configFile := filepath.Join(xdg.ConfigHome, "peon", "config.toml")

	if fs.Exists(configFile) {
		return FromPath(configFile)
	}

	return Default(), nil
}

func FromPath(path string) (*Config, error) {
	if !fs.Exists(path) {
		return nil, fmt.Errorf("unknown config file path specified: %s", path)
	}

	var config Config
	_, err := toml.DecodeFile(path, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func FromContext(ctx context.Context) *Config {
	cfg, ok := ctx.Value(ContextKey).(*Config)
	if !ok {
		return Default()
	}

	return cfg
}
