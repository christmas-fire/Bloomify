package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

// Конфигурация Redis
type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// Загрузка конфигурации Redis
func LoadConfigRedis(path string) (*RedisConfig, error) {
	viper.SetConfigName("config_redis")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config RedisConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error decode config file: %w", err)
	}

	return &config, nil
}
