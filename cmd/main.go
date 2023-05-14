package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	handler "github.com/nurmeden/students-service/internal/app/handlers"
	"github.com/nurmeden/students-service/internal/app/repository"
	"github.com/nurmeden/students-service/internal/app/usecase"
	"github.com/nurmeden/students-service/internal/database"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	if err := godotenv.Load(".env"); err != nil {
		log.Print("No .env file found")
	}

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

	dbName := os.Getenv("DATABASE_NAME")
	mongoURI := os.Getenv("MONGODB_URI")
	collectionName := os.Getenv("COLLECTION_NAME")
	fmt.Printf("mongoURI: %v\n", mongoURI)

	client, err := database.SetupDatabase(context.Background(), mongoURI)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("client: %v\n", client)

	defer client.Disconnect(context.Background())

	counter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "requests_total",
		Help: "Total number of requests processed by the server, partitioned by status code and HTTP method.",
	}, []string{"code", "method"})

	prometheus.MustRegister(counter)

	studentRepo, _ := repository.NewStudentRepository(client, dbName, collectionName, redisClient, logger)

	studentUsecase := usecase.NewStudentUsecase(*studentRepo, logger, redisClient)

	studentHandler := handler.NewStudentHandler(studentUsecase, logger)

	router := gin.Default()

	api := router.Group("/api/")
	{
		studentsGroup := api.Group("/students/")
		studentsGroup.Use(handler.AuthMiddleware())
		{
			api.POST("/students/", studentHandler.CreateStudent)
			api.GET("/students/:id", studentHandler.GetStudentByID)
			api.PUT("/students/:id", studentHandler.UpdateStudents)
			api.DELETE("/students/:id", studentHandler.DeleteStudent)
			api.GET("/students/:id/courses", studentHandler.GetStudentCourses)
			api.GET("/students/:id/students", studentHandler.GetStudentByCoursesID)

		}
		auth := api.Group("/auth/")
		{
			auth.POST("/sign-up", studentHandler.CreateStudent)
			auth.POST("/sign-in", studentHandler.SignIn)
			auth.POST("/logout", studentHandler.Logout)
			auth.POST("/refresh-token", studentHandler.RefreshToken)
		}
	}
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.Run(":8000")
}
