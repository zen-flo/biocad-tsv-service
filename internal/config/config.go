package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type ServerConfig struct {
	Port string `yaml:"port"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type DirConfig struct {
	Input  string `yaml:"input"`
	Output string `yaml:"output"`
}

type Config struct {
	Server ServerConfig `yaml:"server"`
	DB     DBConfig     `yaml:"db"`
	Dirs   DirConfig    `yaml:"dirs"`
}

// LoadConfig reads the YAML file and returns Config
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}

// String brings the config to a string for easy logging
func (c *Config) String() string {
	return fmt.Sprintf("Server{port=%s}, DB{host=%s, port=%d, user=%s}, Dirs{input=%s, output=%s}",
		c.Server.Port, c.DB.Host, c.DB.Port, c.DB.User, c.Dirs.Input, c.Dirs.Output)
}

// Validate checks if the config fields are valid
func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("server.port is required")
	}
	if c.DB.Host == "" {
		return fmt.Errorf("db host is required")
	}
	if c.DB.Port <= 0 || c.DB.Port > 65535 {
		return fmt.Errorf("db port must be between 1 and 65535")
	}
	if c.DB.User == "" {
		return fmt.Errorf("db user is required")
	}
	if c.DB.Password == "" {
		return fmt.Errorf("db password is required")
	}
	if c.DB.Name == "" {
		return fmt.Errorf("db name is required")
	}
	if c.Dirs.Input == "" {
		return fmt.Errorf("dirs input is required")
	}
	if c.Dirs.Output == "" {
		return fmt.Errorf("dirs output is required")
	}
	return nil
}
