package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Caddy    CaddyConfig    `yaml:"caddy"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type DatabaseConfig struct {
	Path string `yaml:"path"`
}

type CaddyConfig struct {
	BinaryPath        string `yaml:"binary_path"`
	CaddyfilePath     string `yaml:"caddyfile_path"`
	CaddyfileSitesDir string `yaml:"caddyfile_sites_dir"`
	ServerDomain      string `yaml:"server_domain"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Server: ServerConfig{Port: 8080},
		Database: DatabaseConfig{
			Path: "./data/caddy.db",
		},
		Caddy: CaddyConfig{
			BinaryPath:        "/usr/sbin/caddy",
			CaddyfilePath:     "/app/caddy_config/Caddyfile",
			CaddyfileSitesDir: "/app/caddy_config/sites",
			ServerDomain:      "",
		},
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
