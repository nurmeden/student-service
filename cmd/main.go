package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis"
	"github.com/nurmeden/students-service/internal/database"
	"github.com/sirupsen/logrus"
)

func main() {
	logfile, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		log.Fatal(err)
	}

	defer logfile.Close()

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(logfile)
	logger.SetLevel(logrus.DebugLevel)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	defer redisClient.Close()

	_, err = redisClient.Ping().Result()
	if err != nil {
		fmt.Println(err.Error())
		logger.Fatal("Ошибка подключения к Redis:", err)
	}

	client, err := database.SetupDatabase()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer client.Disconnect(context.Background())
}
