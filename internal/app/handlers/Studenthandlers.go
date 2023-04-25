package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nurmeden/students-service/internal/app/model"
	"github.com/nurmeden/students-service/internal/app/usecase"
)

type StudentHandler struct {
	studentUsecase usecase.StudentUsecase
	//logger         logger.Logger
}

// func NewStudentHandler(studentUsecase usecase.StudentUsecase, logger logger.Logger) *StudentHandler {
func NewStudentHandler(studentUsecase usecase.StudentUsecase) *StudentHandler {
	return &StudentHandler{
		studentUsecase: studentUsecase,
		// logger:         logger,
	}
}

func (h *StudentHandler) CreateStudent(c *gin.Context) {
	var student model.Student

	err := c.ShouldBindJSON(&student)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode request body"})
		return
	}

	createdStudent, err := h.studentUsecase.CreateStudent(context.Background(), &student)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create student"})
		return
	}

	c.JSON(http.StatusCreated, createdStudent)
}

func (h *StudentHandler) GetStudentByID(c *gin.Context) {
	studentID := c.Param("id")
	fmt.Printf("studentID: %v\n", studentID)
	student, err := h.studentUsecase.GetStudentByID(context.Background(), studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("student get by id: %v\n", student)

	c.JSON(http.StatusOK, student)
}

func (h *StudentHandler) GetStudentByCoursesID(c *gin.Context) {
	courseID := c.Param("id")
	fmt.Printf("courseID: %v\n", courseID)
	student, err := h.studentUsecase.GetStudentByCoursesID(context.Background(), courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("student get by id: %v\n", student)

	c.JSON(http.StatusOK, student)
}

func (h *StudentHandler) SignIn(c *gin.Context) {
	var signInData model.SignInData
	// Получаем данные аутентификации из тела запроса
	err := c.ShouldBindJSON(&signInData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode request body"})
		return
	}

	// Вызываем метод SignIn в StudentUsecase для проведения аутентификации
	authResult, err := h.studentUsecase.SignIn(&signInData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to authenticate"})
		return
	}

	// Возвращаем токен аутентификации в ответе
	c.JSON(http.StatusOK, gin.H{"token": authResult.Token})
}

func (sc *StudentHandler) GetStudentCourses(c *gin.Context) {
	// Получаем идентификатор студента из URL-параметров
	studentID := c.Param("id")

	// Отправляем запрос к второму микросервису, отвечающему за курсы, используя HTTP-запрос
	// Например, можно использовать стандартный пакет net/http для выполнения GET-запроса
	resp, err := http.Get("http://localhost:8080/api/courses/" + studentID + "/courses")
	fmt.Printf("resp: %v\n", resp)
	if err != nil {
		// Обработка ошибки
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get student courses"})
		return
	}
	defer resp.Body.Close()

	fmt.Printf("resp: %v\n", resp)

	// Чтение ответа и обработка данных
	// Например, можно использовать пакет encoding/json для декодирования JSON-ответа
	var course *model.CourseResponse
	err = json.NewDecoder(resp.Body).Decode(&course)
	if err != nil {
		// Обработка ошибки
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode student courses"})
		return
	}

	// Отправляем данные о курсах в качестве ответа
	c.JSON(http.StatusOK, gin.H{"courses": course})
}
