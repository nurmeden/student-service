package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
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

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(logfile)
	logrus.SetLevel(logrus.DebugLevel)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	defer redisClient.Close()

	// Проверка подключения к Redis
	_, err = redisClient.Ping().Result()
	if err != nil {
		log.Fatal("Ошибка подключения к Redis:", err)
	}

	client := database.SetupDatabase()
	logrus.Infoln("Creating Database")

	defer client.Disconnect(context.Background())

	counter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "requests_total",
		Help: "Total number of requests processed by the server, partitioned by status code and HTTP method.",
	}, []string{"code", "method"})

	prometheus.MustRegister(counter)

	studentRepo, _ := repository.NewStudentRepository(client, "studentsdb", "students", redisClient)

	studentUsecase := usecase.NewStudentUsecase(*studentRepo)

	studentHandler := handler.NewStudentHandler(studentUsecase)

	router := gin.Default()

	// Регистрация маршрутов
	api := router.Group("/api/")
	{
		api.POST("/students", studentHandler.CreateStudent)
		api.GET("/students/:id", studentHandler.GetStudentByID)
		api.GET("/students/:id/courses", studentHandler.GetStudentCourses)
		api.GET("/students/:id/students", studentHandler.GetStudentByCoursesID)
		auth := api.Group("/auth/")
		{
			auth.POST("/sign-up", studentHandler.CreateStudent)
			auth.POST("/sign-in", studentHandler.SignIn)
		}
	}

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.Run(":8000")
}
