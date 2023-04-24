package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"student-service/internal/app/model"
	"student-service/internal/app/usecase"

	"github.com/gin-gonic/gin"
)

type StudentHandler struct {
	studentUsecase usecase.StudentUsecase
	// logger         logger.Logger
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

// func (h *StudentHandler) CreateStudent(w http.ResponseWriter, r *http.Request) {
// 	var student model.Student
// 	err := json.NewDecoder(r.Body).Decode(&student)
// 	if err != nil {
// 		// h.logger.Errorf("Failed to decode request body: %v", err)
// 		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
// 		return
// 	}

// 	createdStudent, err := h.studentUsecase.CreateStudent(&student)
// 	if err != nil {
// 		// h.logger.Errorf("Failed to create student: %v", err)
// 		http.Error(w, "Failed to create student", http.StatusInternalServerError)
// 		return
// 	}

// 	response, err := json.Marshal(createdStudent)
// 	if err != nil {
// 		// h.logger.Errorf("Failed to marshal response: %v", err)
// 		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	w.Write(response)
// }
