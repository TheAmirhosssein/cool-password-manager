package config

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/TheAmirhosssein/cool-password-manage/internal/utils/oprfutils"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

var (
	config *Config
	once   sync.Once
)

type (
	Config struct {
		APP    `yaml:"app"`
		HTTP   `yaml:"http"`
		DB     `yaml:"db"`
		Redis  `yaml:"redis"`
		Opaque `yaml:"opaque"`
	}

	APP struct {
		Name              string `env-required:"true" yaml:"name"`
		Version           string `env-required:"true" yaml:"version"`
		AESKey            string `env-required:"true" yaml:"aes_key" env:"AES_KEY"`
		RootPath          string
		TemplatePath      string `env-required:"true" yaml:"template_path" env:"TEMPLATE_PATH"`
		StaticPath        string `env-required:"true" yaml:"static_path" env:"STATIC_PATH"`
		TwoFactorDuration int    `env-required:"true" yaml:"two_factor_duration" env:"TWO_FACTOR_DURATION"`
		SecretKey         string `env-required:"true" yaml:"secret_key" env:"SECRET_KEY"`
		DefaultPage       int    `env-required:"true" yaml:"default_page" env:"DEFAULT_PAGE"`
		DefaultPageSize   int    `env-required:"true" yaml:"default_page_size" env:"DEFAULT_PAGE_SIZE"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port"`
		Host string `env-required:"true" yaml:"host"`
	}

	DB struct {
		URL string `env-required:"true" env:"POSTGRES_URL"`
	}

	Redis struct {
		URL string `env-required:"true" env:"REDIS_URL"`
	}

	Opaque struct {
		ServerID             string `env-required:"true" yaml:"server_id" env:"SERVER_ID"`
		PublicKeyPath        string `env-required:"true" yaml:"public_key_path" env:"PUBLIC_KEY_PATH"`
		PrivateKeyPath       string `env-required:"true" yaml:"private_key_path" env:"PRIVATE_KEY_PATH"`
		OprfKeyPath          string `env-required:"true" yaml:"oprf_key_path" env:"ORFP_KEY_PATH"`
		RegistrationDuration int    `env-required:"true" yaml:"registration_duration" env:"RegistrationDuration"`
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
	return &Config{Opaque: createTestCodes()}
}

func (c *Config) GetAESSecretKey() ([]byte, error) {
	if InTestMode() {
		return base64.StdEncoding.DecodeString("syaZbz9ca3SZ51GUdyx3F//e89Hgfr2XuHHn4VdnMQU=")
	}

	keyBytes, err := base64.StdEncoding.DecodeString(c.APP.AESKey)
	return keyBytes, err
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

	index := strings.Index(path, `/internal`)
	if index != -1 {
		path = path[:index]
	}

	return path, err
}

func createTestCodes() Opaque {
	rootPath, err := getRootPath()
	if err != nil {
		panic(err)
	}

	keysPath := fmt.Sprint(rootPath, "/internal/infrastructure/opaque/keys/test/")
	err = oprfutils.GenerateAndSaveKeys(keysPath)
	if err != nil {
		panic(err)
	}

	return Opaque{
		ServerID:       "something",
		PublicKeyPath:  fmt.Sprint(keysPath, "public.bin"),
		PrivateKeyPath: fmt.Sprint(keysPath, "private.bin"),
		OprfKeyPath:    fmt.Sprint(keysPath, "oprf.bin"),
	}
}
