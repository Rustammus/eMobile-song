package config

import (
	"eMobile/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"sync"
)

type Config struct {
	Server struct {
		Port           string `yaml:"port" env:"APP_PORT" env-required:""`
		LogLevel       string `yaml:"log_level" env:"APP_LOG_LEVEL" env-required:""`
		EnableSwag     bool   `yaml:"enable_swag" env:"APP_ENABLE_SWAG" env-default:"true"`
		Host           string `yaml:"host" env:"APP_HOST" env-default:"localhost"`
		ExternalHost   string `yaml:"external_host" env:"APP_EXTERNAL_HOST" env-required:""`
		ExternalPort   string `yaml:"external_port" env:"APP_EXTERNAL_PORT" env-required:""`
		InfoServiceUrl string `yaml:"info_service_url" env:"APP_INFO_SERVICE_URL" env-required:""`
		PagLimit       int    `yaml:"pag_limit" env:"APP_PAG_LIMIT" env-required:""`
	} `yaml:"server"`
	Storage Storage `yaml:"storage"`
}

type Storage struct {
	Username  string `yaml:"username" env:"POSTGRES_USER" env-required:""`
	Password  string `yaml:"password" env:"POSTGRES_PASSWORD" env-required:""`
	Host      string `yaml:"host" env:"POSTGRES_HOST" env-required:""`
	Port      string `yaml:"port" env:"POSTGRES_PORT" env-required:""`
	Database  string `yaml:"database" env:"POSTGRES_DB" env-required:""`
	Migration string `yaml:"migration" env:"POSTGRES_MIGRATION"`
}

var once sync.Once
var instance *Config

func GetConfig(logger logging.Logger) *Config {
	once.Do(func() {
		confPath := "app_dev.env"
		if value, ok := os.LookupEnv("CONFIG_FILE"); ok {
			confPath = value
		} else {
			logger.Info("env: CONFIG_FILE not set. Default: \"app_dev.env\"")
		}

		logger.Infof("Try to read config file %s", confPath)
		instance = &Config{}
		err := cleanenv.ReadConfig(confPath, instance)
		if err != nil {
			logger.Infof("Failed to read config file. Try read env.\t%s", err.Error())
			instance = &Config{}
			err = cleanenv.ReadEnv(instance)
			if err != nil {
				logger.Fatalf("Failed to read environment variables. Abort start app.\t%s", err.Error())
			}
		}

		logger.Info("Successfully read config.")
	})
	return instance
}
