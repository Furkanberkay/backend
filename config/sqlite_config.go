package config

import (
	"log"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type EnvSQLiteConfig struct {
	ServerPort string `env:"SERVER_PORT,required"`

	DBDriver string `env:"DB_DRIVER" envDefault:"sqlite"`

	SQLitePath string `env:"SQLITE_PATH" envDefault:"./data/app.db"`

	// Postgres (ileride kullanacaksın)
	DBHost     string `env:"DB_HOST"`
	DBName     string `env:"DB_NAME"`
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`
	DBSslMode  string `env:"DB_SSLMODE" envDefault:"disable"`
	DBPort     string `env:"DB_PORT" envDefault:"5432"`
}

func NewEnvSQLiteConfig() *EnvSQLiteConfig {
	if err := godotenv.Load(); err != nil {
		log.Fatal("unable to load the .env: ", err)
	}

	var cfg EnvSQLiteConfig
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("unable to load variables from the env: %v", err)
	}
	return &cfg
}
