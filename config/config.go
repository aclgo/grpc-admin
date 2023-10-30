package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	ApiPort  string `mapstructure:"API_PORT"`
	DbDriver string `mapstructure:"DB_DRIVER"`
	DbURI    string `mapstructure:"DB_URI"`
	DevMode  bool   `mapstructure:"DEV_MODE"`
	Encoding string `mapstructure:"ENCODING"`
	LogLevel string `mapstructure:"LOG_LEVEL"`
	Tracer
	Metric
	AdminClientServiceConfig
	UserClientServiceConfig
	MailClientServiceConfig
}

type AdminClientServiceConfig struct {
	Name string `mapstructure:"ADMIN_SERVICE_NAME"`
	Addr string `mapstructure:"ADMIN_SERVICE_ADDR"`
}

type UserClientServiceConfig struct {
	Name string `mapstructure:"USER_SERVICE_NAME"`
	Addr string `mapstructure:"USER_SERVICE_ADDR"`
}

type MailClientServiceConfig struct {
	Name string `mapstructure:"MAIL_SERVICE_NAME"`
	Addr string `mapstructure:"MAIL_SERVICE_ADDR"`
}

type Metric struct {
	Name        string `mapstructure:"METRIC_NAME"`
	ExporterURL string `mapstructure:"METRIC_EXPORTER_NAME"`
}

type Tracer struct {
	Name        string `mapstructure:"TRACER_NAME"`
	ExporterURL string `mapstructure:"TRACER_EXPORTER_NAME"`
}

func NewConfig(path string) *Config {
	var cfg Config

	_, ok := os.LookupEnv("PROD")
	if !ok {

		log.Printf("env variable PROD is false, reading from file")

		cfg.LoadFromFile(path)
		return &cfg
	}

	err := cfg.LoadFromEnv()
	if err != nil {
		log.Fatalf("LoadFromEnv: %v", err)
	}

	return &cfg

}

func (c Config) LoadFromEnv() error {
	return nil
}

func (c Config) LoadFromFile(path string) {

	viper.SetConfigName("app")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("LoadFile.ReadInConfig: %v", err)
	}

	if err := viper.Unmarshal(&c); err != nil {
		log.Fatalf("LoadFile.Unmarshal: %v", err)
	}
}
