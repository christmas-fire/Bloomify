package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

// Конфигурация PostgreSQL
type PostgreSQLConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Sslmode  string `yaml:"sslmode"`
}

// Загрузка конфигурации PostgreSQL
func LoadConfigPostgres(path string) (*PostgreSQLConfig, error) {
	viper.SetConfigName("config_postgres")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config PostgreSQLConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error decode config file: %w", err)
	}

	return &config, nil
}
