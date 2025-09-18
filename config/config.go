package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

var (
	config *Config
	once   sync.Once
)

type (
	Config struct {
		APP  `yaml:"app"`
		HTTP `yaml:"http"`
		DB   `yaml:"db"`
	}

	APP struct {
		Name                   string `env-required:"true" yaml:"name"`
		Version                string `env-required:"true" yaml:"version"`
		SecretKey              string `env-required:"true" yaml:"secret_key" env:"SECRET_KEY"`
		RootPath               string
		AuthenticatorTemplates string `env-required:"true" yaml:"auth_templates" env:"AUTH_TEMPLATE"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port"`
		Host string `env-required:"true" yaml:"host"`
	}

	DB struct {
		URL string `env-required:"true" env:"POSTGRES_URL"`
	}
)

func newConfig() (*Config, error) {
	conf := &Config{}
	if _, err := os.Stat(".env"); !errors.Is(err, os.ErrNotExist) {
		err := godotenv.Load(".env")
		if err != nil {
			return nil, err
		}
	}

	err := cleanenv.ReadConfig("./config/config.yaml", conf)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(conf)
	if err != nil {
		return nil, fmt.Errorf("env error: %w", err)
	}

	conf.APP.RootPath, err = getRootPath()
	if err != nil {
		return nil, fmt.Errorf("getting root path error: %w", err)
	}

	return conf, nil
}

func GetConfig() *Config {
	once.Do(func() {
		var err error
		config, err = newConfig()
		if err != nil {
			panic(err.Error())
		}
	})
	return config
}

func GetTestConfig() *Config {
	return &Config{}
}

func (c *Config) GetSecretKey() string {
	if InTestMode() {
		return "secretestkey"
	}

	return c.APP.SecretKey
}

func InTestMode() bool {
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test.") {
			return true
		}
	}
	return false
}

func getRootPath() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}

	index := strings.Index(path, `\internal`)
	if index != -1 {
		path = path[:index]
	}

	return path, err
}
