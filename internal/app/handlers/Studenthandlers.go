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
	"github.com/sirupsen/logrus"
)

type StudentHandler struct {
	studentUsecase usecase.StudentUsecase
}

// func NewStudentHandler(studentUsecase usecase.StudentUsecase, logger logger.Logger) *StudentHandler {
func NewStudentHandler(studentUsecase usecase.StudentUsecase) *StudentHandler {
	return &StudentHandler{
		studentUsecase: studentUsecase,
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

	c.JSON(http.StatusOK, gin.H{"token": authResult.Token})
}

func (sc *StudentHandler) GetStudentCourses(c *gin.Context) {
	studentID := c.Param("id")

	resp, err := http.Get("http://localhost:8080/api/courses/" + studentID + "/courses")
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
