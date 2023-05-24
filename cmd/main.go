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
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	docs "github.com/nurmeden/students-service/docs"
)

// @title Students API
// @version 1
// @description API for managing students
// @host		localhost:8000
// @BasePath /api/
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
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	defer redisClient.Close()
	fmt.Printf("redisClient: %v\n", redisClient)
	ll, err := redisClient.Ping().Result()
	fmt.Println(ll)
	if err != nil {
		fmt.Println(err.Error())
		logger.Fatal("Ошибка подключения к Redis:", err)
	}

	fmt.Println(ll)

	// // dbName := os.Getenv("DATABASE_NAME")
	// mongoURI := os.Getenv("MONGODB_URI")
	// // collectionName := os.Getenv("COLLECTION_NAME")

	client, err := database.SetupDatabase(context.Background())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer client.Disconnect(context.Background())

	counter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "requests_total",
		Help: "Total number of requests processed by the server, partitioned by status code and HTTP method.",
	}, []string{"code", "method"})

	prometheus.MustRegister(counter)

	studentRepo, _ := repository.NewStudentRepository(client, "studentsdb", "students", redisClient, logger)

	studentUsecase := usecase.NewStudentUsecase(*studentRepo, logger, redisClient)

	studentHandler := handler.NewStudentHandler(studentUsecase, logger)

	router := gin.Default()
	docs.SwaggerInfo.BasePath = "/api"
	api := router.Group("/api/")
	{
		api.POST("/students/", studentHandler.CreateStudent)
		api.GET("/students/:id", studentHandler.GetStudentByID)
		api.GET("/students/:id/students", studentHandler.GetStudentsByCourseID)
		studentsGroup := api.Group("/students")
		studentsGroup.Use(handler.AuthMiddleware())
		{
			studentsGroup.PUT("/:id", studentHandler.UpdateStudents)
			studentsGroup.DELETE("/:id", studentHandler.DeleteStudent)
			studentsGroup.GET("/:id/courses", studentHandler.GetStudentCourses)
		}
		auth := api.Group("/auth/")
		{
			auth.POST("/sign-up", studentHandler.CreateStudent)
			auth.POST("/sign-in", studentHandler.SignIn)
			auth.POST("/logout", studentHandler.Logout)
			auth.POST("/refresh-token", studentHandler.RefreshToken)
		}

		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.Run(":8001")
}
