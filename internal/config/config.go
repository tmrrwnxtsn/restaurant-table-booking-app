package config

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/qiangxue/go-env"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

const envVarsPrefix = "API_"

// Config содержит настройки сервиса.
type Config struct {
	// BindAddr представляет адрес хоста.
	BindAddr string `yaml:"bind_addr" env:"BIND_ADDR"`
	// DSN (data source name) является строкой подключения к базе данных.
	DSN string `yaml:"dsn" env:"DSN,secret"`
	// LogLevel представляет уровень логгирования.
	LogLevel string `yaml:"log_level" env:"LOG_LEVEL"`
}

// Validate проверяет, достаточно ли настроек для запуска сервиса.
func (c Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.BindAddr, validation.Required),
		validation.Field(&c.DSN, validation.Required),
		validation.Field(&c.LogLevel, validation.Required),
	)
}

// Load загружает настройки сервиса из переменных среды и, если их не окажется, из yml-файла.
func Load(ymlConfigPath string) (*Config, error) {
	cfg := Config{}

	// загрузка конфигурационных значений из yml-файла
	bytes, err := ioutil.ReadFile(ymlConfigPath)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(bytes, &cfg); err != nil {
		return nil, err
	}

	// загрузка конфигурационных значений из переменных среды, имеющих префикс envVarsPrefix
	if err = env.New(envVarsPrefix, log.Printf).Load(&cfg); err != nil {
		panic(err)
	}

	if err = cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}
