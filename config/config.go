package config

import (
	"os"

	"github.com/joho/godotenv"
)

type database struct {
	Url string
}

type Config struct {
	Database database
}

func New() *Config {
	godotenv.Load()
	return &Config{
		Database: database{
			Url: os.Getenv("DATABASE_URL") ,
		},
	}
}