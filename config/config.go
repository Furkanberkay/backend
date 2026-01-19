package config

import (
	"log"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type EnvConfig struct {
	ServerPort string `env:"SERVER_PORT,required"`
	JWTSecret  string `env:"SECRET_KEY,required"`
	DBDriver   string `env:"DB_DRIVER,required"`

	DBHost     string `env:"DB_HOST"`
	DBName     string `env:"DB_NAME"`
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`
	DBPort     string `env:"DB_PORT"`
	DBSslMode  string `env:"DB_SSLMODE"`

	SQLitePath string `env:"SQLITE_PATH"`
}

func NewEnvConfig() *EnvConfig {
	_ = godotenv.Load()

	config := EnvConfig{}

	if err := env.Parse(&config); err != nil {
		log.Fatalf("Config yüklenemedi: %v", err)
	}

	return &config
}
