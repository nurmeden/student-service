package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	// _ "github.com/nurmeden/students-service/cmd/docs"
	"github.com/nurmeden/students-service/internal/app/model"
	"github.com/nurmeden/students-service/internal/app/usecase"
	"github.com/sirupsen/logrus"
)

const endpoint = "http://localhost:8080/api/courses/"

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

// CreateStudent godoc
// @Summary Create a new student
// @Description Create a new student with the input payload
// @Tags Students
// @Accept  json
// @Produce  json
// @Param student body model.Student true "Student data"
// @Success 201 {object} model.Student
// @Router /api/students/ [post]
// @Router /api/auth/sign-up/ [post]
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

// GetStudentByID godoc
// @Summary Get student by ID
// @Description Get student by ID
// @Tags students
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Success 200 {object} model.Student
// @Router /students/{id} [get]
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

// UpdateStudents godoc
// @Summary Update a student by ID
// @Description Update a student by ID
// @Tags students
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Param student body model.Student true "Student object"
// @Success 200 {object} model.Student

// @Router /students/{id} [put]
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

// DeleteStudent godoc
// @Summary Delete a student by ID
// @Description Delete a student by its ID
// @Tags students
// @Param id path int true "Student ID"
// @Success 200 {string} string "message: Студент успешно удален"

// @Router /students/{id} [delete]
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

// SignIn godoc
// @Summary Sign in a student
// @Description Authenticates a student using their email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param signInData body model.SignInData true "Sign in data"
// @Success 200 {object} TokenResponse

// @Router /api/sign-in [post]
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

	// refreshToken := uuid.New().String()
	// err = h.studentUsecase.SaveRefreshToken(authResult.UserID, refreshToken)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save refresh token"})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{"token": authResult.Token})
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
