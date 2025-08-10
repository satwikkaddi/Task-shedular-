package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	RedisAddr string
	HTTPPort  string
	Env       string
	Timeout   time.Duration
}

func LoadConfig() *Config {
	viper.SetDefault("REDIS_ADDR", "127.0.0.1:6379")
	viper.SetDefault("HTTP_PORT", "8080")
	iper := viper.New()
	iper.AutomaticEnv()

	return &Config{
		RedisAddr: iper.GetString("REDIS_ADDR"),
		HTTPPort:  iper.GetString("HTTP_PORT"),
		Env:       iper.GetString("ENV"),
		Timeout:   10 * time.Second,
	}
}
