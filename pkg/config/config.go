package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

var Cfg *Config

type (
	Config struct {
		App    `yaml:"app"`
		HTTP   `yaml:"http"`
		Logger `yaml:"logger"`
		DB     `yaml:"db"`
		Redis  `yaml:"redis"`
		JWT    `yaml:"jwt"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	Logger struct {
		Level string `env-required:"true" yaml:"level" env:"LOG_LEVEL"`
	}

	DB struct {
		URL string `env-required:"true" yaml:"url" env:"DB_URL"`
	}

	Redis struct {
		URL string `env-required:"true" yaml:"url" env:"REDIS_URL"`
	}

	JWT struct {
		Secret   string `env-required:"true" yaml:"secret" env:"JWT_SECRET"`
		Duration int64  `env-required:"true" yaml:"duration" env:"JWT_DURATION"`
	}
)

func newConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func Init() error {
	cfg, err := newConfig()
	if err != nil {
		return err
	}

	Cfg = cfg

	return nil
}
