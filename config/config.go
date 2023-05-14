package config

import (
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	JWTKey string `env:"JWT_KEY" envDefault:"supersecret"`
}

type HTTP struct {
	AppPort string `env:"APP_PORT" envDefault:"8000"`
}

type Database struct {
	DBHost         string `env:"MONGO_HOST" envDefault:"studentsdb"`
	DBPort         string `env:"MONGO_PORT" envDefault:"27017"`
	DBName         string `env:"DATABASE_NAME" envDefault:"studentsdb"`
	CollectionName string `env:"DATABASE_NAME" envDefault:"students"`
}

func PrepareEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln(err)
	}
}
