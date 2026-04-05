package config

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"atomicptr.dev/bits"
	"atomicptr.dev/peon/pkg/fs"
	"github.com/BurntSushi/toml"
	"github.com/adrg/xdg"
)

type contextKey string

const ContextKey contextKey = "atomicptr.dev/peon/pkg/config"

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

	if bits.PathExists(configFile) {
		return FromPath(configFile)
	}

	return Default(), nil
}

func FromPath(path string) (*Config, error) {
	if !bits.PathExists(path) {
		return nil, fmt.Errorf("unknown config file path specified: %s", path)
	}

	var config Config
	_, err := toml.DecodeFile(path, &config)
	if err != nil {
		return nil, err
	}

	slog.Debug("config file found", "path", strings.ReplaceAll(filepath.Clean(path), "\n", ""))

	return &config, nil
}

func FromContext(ctx context.Context) *Config {
	cfg, ok := ctx.Value(ContextKey).(*Config)
	if !ok {
		return Default()
	}

	return cfg
}
