package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	DatabaseConfig DatabaseConfig `yaml:"database"`
	ServerConfig   ServerConfig   `yaml:"server"`
	LoggerConfig   LoggerConfig   `yaml:"logger"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"-"`
	Password string `yaml:"-"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type LoggerConfig struct {
	Level string `yaml:"level"`
}

func MustLoadConfig() *Config {
	config, err := LoadConfig()
	if err != nil {
		panic(err)
	}

	return config
}

func LoadConfig() (*Config, error) {
	configPath := "configs/config.yaml"
	if os.Getenv("CONFIG_PATH") != "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	envPath := "configs/.env"
	if os.Getenv("ENV_PATH") != "" {
		envPath = os.Getenv("ENV_PATH")
	}

	if err := godotenv.Load(envPath); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(file, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	config.DatabaseConfig.Username = os.Getenv("DATABASE_USERNAME")
	config.DatabaseConfig.Password = os.Getenv("DATABASE_PASSWORD")

	return &config, nil
}

func (c *Config) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Group("database",
			slog.String("host", c.DatabaseConfig.Host),
			slog.String("port", c.DatabaseConfig.Port),
			slog.String("dbname", c.DatabaseConfig.DBName),
			slog.String("sslmode", c.DatabaseConfig.SSLMode),
			slog.Bool("has_username", c.DatabaseConfig.Username != ""),
			slog.Bool("has_password", c.DatabaseConfig.Password != ""),
		),
		slog.Group("server",
			slog.String("host", c.ServerConfig.Port),
		),
		slog.Group("logger",
			slog.String("level", c.LoggerConfig.Level),
		),
	)
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
		c.SSLMode,
	)
}

func (c *Config) Validate() error {
	var errors []string

	if err := c.DatabaseConfig.Validate(); err != nil {
		errors = append(errors, fmt.Sprintf("error validating database config: %s", err))
	}
	if err := c.ServerConfig.Validate(); err != nil {
		errors = append(errors, fmt.Sprintf("error validating server config: %s", err))
	}
	if err := c.LoggerConfig.Validate(); err != nil {
		errors = append(errors, fmt.Sprintf("error validating logger config: %s", err))
	}

	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, "; "))
	}

	return nil
}

func (c *DatabaseConfig) Validate() error {
	var errors []string

	if c.Host == "" {
		errors = append(errors, "host is required")
	}
	if c.Port == "" {
		errors = append(errors, "port is required")
	}
	if c.Username == "" {
		errors = append(errors, "username is required")
	}
	if c.Password == "" {
		errors = append(errors, "password is required")
	}
	if c.DBName == "" {
		errors = append(errors, "dbname is required")
	}
	if c.SSLMode == "" {
		errors = append(errors, "sslmode is required")
	}

	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, ", "))
	}

	return nil
}

func (c *ServerConfig) Validate() error {
	if c.Port == "" {
		return fmt.Errorf("port is required")
	}

	return nil
}

func (c *LoggerConfig) Validate() error {
	if c.Level == "" {
		return fmt.Errorf("level is required")
	}

	return nil
}
