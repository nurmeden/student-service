package config

import (
	"errors"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	JWTKey string `env:"JWT_KEY" envDefault:"supersecret"`
	Logger Logger
}

type Logger struct {
	Development       bool
	DisableCaller     bool
	DisableStacktrace bool
	Encoding          string
	Level             string
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

func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}

		return nil, err
	}

	return v, nil
}

func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, err
}
