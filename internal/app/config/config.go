package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	ServiceHost string
	ServicePort int

	MinioHost      string
	MinioPort      int
	MinioAccessKey string
	MinioSecretKey string
	MinioBucket    string

	SecretKey        string
	ExpiresAtMinutes time.Duration

	RedisHost           string
	RedisPort           int
	RedisUser           string
	RedisPassword       string
	RedisDB             int
	RedisDialTimeoutSec time.Duration
	RedisReadTimeoutSec time.Duration
	RedisAppPrefix      string
}

func NewConfig() (*Config, error) {
	var err error

	configName := "config"
	_ = godotenv.Load()
	if os.Getenv("CONFIG_NAME") != "" {
		configName = os.Getenv("CONFIG_NAME")
	}

	viper.SetConfigName(configName)
	viper.SetConfigType("toml")
	viper.AddConfigPath("config")
	viper.AddConfigPath(".")
	viper.WatchConfig()

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfg := &Config{}           // создаем объект конфига
	err = viper.Unmarshal(cfg) // читаем информацию из файла,
	// конвертируем и затем кладем в нашу переменную cfg
	if err != nil {
		return nil, err
	}

	log.Info("config parsed")

	return cfg, nil
}
