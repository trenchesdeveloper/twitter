package config

import (
	"log"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

type database struct {
	Url string
}

type jwt struct {
	Secret string
	Issuer string
}

type Config struct {
	Database database
	Jwt      jwt
}

func LoadEnv(fileName string) {
	re := regexp.MustCompile(`^(.*` + "twitter" + `)`)

	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	err := godotenv.Load(string(rootPath) + "/" + fileName)

	if err != nil {
		log.Println("error loading .env file222")
		err := godotenv.Load()
		if err != nil {
			return
		}
	}
}

func New() *Config {
	return &Config{
		Database: database{
			Url: os.Getenv("DATABASE_URL"),
		},
		Jwt: jwt{
			Secret: os.Getenv("JWT_SECRET"),
			Issuer: os.Getenv("DOMAIN"),
		},
	}
}
