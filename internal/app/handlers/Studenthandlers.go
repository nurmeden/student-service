package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nurmeden/students-service/internal/app/model"
	"github.com/nurmeden/students-service/internal/app/usecase"
	"github.com/sirupsen/logrus"
)

const endpoint = "http://localhost:8080/api/courses/"

var blackListedTokens = make(map[string]bool)

type StudentHandler struct {
	studentUsecase usecase.StudentUsecase
	logger         *logrus.Logger
}

func NewStudentHandler(studentUsecase usecase.StudentUsecase, logger *logrus.Logger) *StudentHandler {
	return &StudentHandler{
		studentUsecase: studentUsecase,
		logger:         logger,
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
		logrus.Errorf("Failed to create student: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create student"})
		return
	}

	c.JSON(http.StatusCreated, createdStudent)
}

func (h *StudentHandler) GetStudentByID(c *gin.Context) {
	fmt.Printf("c.Request.URL: %v\n", c.Request.URL)
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

func (h *StudentHandler) UpdateStudents(c *gin.Context) {
	studentID := c.Param("id")

	var studentUpdateInput model.Student

	if err := c.ShouldBindJSON(&studentUpdateInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	student, err := h.studentUsecase.UpdateStudent(context.Background(), studentID, &studentUpdateInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": student})
}

func (h *StudentHandler) DeleteStudent(c *gin.Context) {
	studentID := c.Param("id")

	err := h.studentUsecase.DeleteStudent(context.Background(), studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Студент успешно удален"})
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
	err := c.ShouldBindJSON(&signInData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode request body"})
		return
	}

	authResult, err := h.studentUsecase.SignIn(&signInData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to authenticate"})
		return
	}

	refreshToken := uuid.New().String()
	err = h.studentUsecase.SaveRefreshToken(authResult.UserID, refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": authResult.Token, "refresh_token": refreshToken})
}

func (sc *StudentHandler) GetStudentCourses(c *gin.Context) {
	studentID := c.Param("id")

	resp, err := http.Get(endpoint + studentID + "/courses")
	fmt.Printf("resp: %v\n", resp)
	if err != nil {
		// Обработка ошибки
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get student courses"})
		return
	}
	defer resp.Body.Close()

	fmt.Printf("resp: %v\n", resp)

	var course *model.CourseResponse
	err = json.NewDecoder(resp.Body).Decode(&course)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode student courses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"courses": course})
}

func (h *StudentHandler) Logout(c *gin.Context) {
	userID := c.GetString("user_id")                   // получаем ID пользователя из контекста Gin
	err := h.studentUsecase.DeleteRefreshToken(userID) // удаляем refresh token из базы данных
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete refresh token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func (h *StudentHandler) RefreshToken(c *gin.Context) {
	// Получаем refresh token из запроса
	refreshToken := c.PostForm("refresh_token")
	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token is missing"})
		return
	}

	// Проверяем, что refresh token является действительным
	userID, err := h.studentUsecase.ValidateRefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh_token"})
		return
	}

	// Генерируем новый access token
	token, err := h.studentUsecase.GenerateToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	// Отправляем новый access token в ответе
	c.JSON(http.StatusOK, gin.H{"token": token})
}
