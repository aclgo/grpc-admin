package config

import (
	"log"
	"os"
)

type Config struct {
	ApiPort            string
	DbDriver           string `mapstructure:"DB_DRIVER"`
	DbURI              string `mapstructure:"DB_URI"`
	DevMode            bool
	Encoding           string
	LogLevel           string
	Observability      Observability
	AdminClientService ClientServiceConfig `mapstructure:"ADMIN_SERVICE"`
	UserClientService  ClientServiceConfig `mapstructure:"USER_SERVICE"`
	MailClientService  ClientServiceConfig `mapstructure:"MAIL_SERVICE"`
}

type ClientServiceConfig struct {
	Name string `mapstructure:"SERVICE_NAME"`
	Addr string `mapstructure:"SERVICE_ADDR"`
}

type Observability struct {
	CollectorURL string
	ZipkinURL    string
}

func NewConfig(path string) *Config {
	var cfg Config

	_, ok := os.LookupEnv("PROD")
	if !ok {

		log.Printf("env variable PROD is false, reading from file")

		err := cfg.LoadFromFile(path)
		if err != nil {
			log.Fatalf("LoadFromFile: %v", err)
		}

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

func (c Config) LoadFromFile(path string) error {
	return nil
}
