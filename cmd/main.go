package main

import (
	"context"

	"github.com/gin-gonic/gin"
	handler "github.com/nurmeden/students-service/internal/app/handlers"
	"github.com/nurmeden/students-service/internal/app/repository"
	"github.com/nurmeden/students-service/internal/app/usecase"
	"github.com/nurmeden/students-service/internal/database"
)

func main() {
	// Инициализация логгера
	// logger := log.New(os.Stdout, "", log.LstdFlags)

	client := database.SetupDatabase()
	defer client.Disconnect(context.Background())

	studentRepo, _ := repository.NewStudentRepository(client, "taskdb", "students")

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
	router.Run(":8000")
}
