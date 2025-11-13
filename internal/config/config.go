package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Server     ServerConfig   `envPrefix:"SERVER_"`
	Db         DatabaseConfig `envPrefix:"DB_"`
	AdminToken string         `env:"ADMINTOKEN,required"`
}

type ServerConfig struct {
	Host string `env:"HOST" envDefault:"localhost"`
	Port int    `env:"PORT" envDefault:"3530"`
}
type DatabaseConfig struct {
	Host     string `env:"HOST,required"`
	Port     int    `env:"PORT,required"`
	User     string `env:"USER,required"`
	Password string `env:"PASSWORD,required"`
	Dbname   string `env:"NAME,required"`
}

func LoadConfig() (*Config, error) {
	config, err := env.ParseAsWithOptions[Config](env.Options{Prefix: "OSA_"})
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return &config, nil
}
