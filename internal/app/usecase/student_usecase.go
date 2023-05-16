package usecase

import (
	"context"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"github.com/nurmeden/students-service/internal/app/model"
	"github.com/nurmeden/students-service/internal/app/repository"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type StudentUsecase interface {
	CreateStudent(ctx context.Context, student *model.Student) (*model.Student, error)
	GetStudentByID(ctx context.Context, id string) (*model.Student, error)
	GetStudentByCoursesID(ctx context.Context, id string) (*model.Student, error)
	UpdateStudent(ctx context.Context, student_id string, student *model.Student) (*model.Student, error)
	DeleteStudent(ctx context.Context, id string) error
	SignIn(ctx context.Context, signInData *model.SignInData) (*model.AuthToken, error)
	SaveRefreshToken(userID string, refreshToken string) error
	GenerateToken(studentID string) (string, error)
	ValidateRefreshToken(refreshToken string) (string, error)
	DeleteRefreshToken(userID string) error
	GetByEmail(ctx context.Context, email string) (*model.Student, error)
	CheckEmailExistence(ctx context.Context, email string) (bool, error)
}

const jwtSecret = "dfhdfjhgdjkff"

type studentUsecase struct {
	studentRepo  repository.StudentRepository
	logger       *logrus.Logger
	jwtSecret    []byte
	jwtGenerator *jwt.Token
	cache        *redis.Client
}

func NewStudentUsecase(studentRepo repository.StudentRepository, logger *logrus.Logger, cache *redis.Client) StudentUsecase {
	jwtGenerator := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{})
	jwtGenerator.Claims.(jwt.MapClaims)["exp"] = time.Now().Add(time.Hour * 24).Unix()

	return &studentUsecase{
		studentRepo:  studentRepo,
		logger:       logger,
		jwtSecret:    []byte(jwtSecret),
		jwtGenerator: jwtGenerator,
		cache:        cache,
	}
}

func (u *studentUsecase) CreateStudent(ctx context.Context, student *model.Student) (*model.Student, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(student.Password), bcrypt.DefaultCost)
	if err != nil {
		u.logger.Errorf("Error generating password hash: %v", err)
		return nil, err
	}
	student.Password = string(hashedPassword)

	return u.studentRepo.CreateStudent(ctx, student)
}

func (u *studentUsecase) GetStudentByID(ctx context.Context, id string) (*model.Student, error) {
	return u.studentRepo.GetStudentByID(ctx, id)
}

func (u *studentUsecase) GetStudentByCoursesID(ctx context.Context, id string) (*model.Student, error) {
	return u.studentRepo.GetStudentByCoursesID(ctx, id)
}

func (u *studentUsecase) UpdateStudent(ctx context.Context, student_id string, student *model.Student) (*model.Student, error) {
	student, err := u.studentRepo.GetStudentByID(ctx, student_id)
	if err != nil {
		return nil, err
	}

	return u.studentRepo.UpdateStudents(ctx, student)
}

func (u *studentUsecase) DeleteStudent(ctx context.Context, id string) error {
	return u.studentRepo.Delete(ctx, id)
}

func (u *studentUsecase) GetByEmail(ctx context.Context, email string) (*model.Student, error) {
	return u.studentRepo.GetByEmail(ctx, email)
}

func (u *studentUsecase) CheckEmailExistence(ctx context.Context, email string) (bool, error) {
	return u.studentRepo.CheckEmailExistence(ctx, email)
}

func (u *studentUsecase) SignIn(ctx context.Context, signInData *model.SignInData) (*model.AuthToken, error) {
	student, err := u.studentRepo.GetByEmail(ctx, signInData.Email)
	if err != nil {
		u.logger.Errorf("Error retrieving student with email %s: %v", signInData.Email, nil)
		return nil, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(student.Password), []byte(signInData.Password))
	if err != nil {
		u.logger.Errorf("Incorrect password for student with email %s", signInData.Email)
		return nil, err
	}
	token, err := u.GenerateToken(student.ID)
	if err != nil {
		return nil, err
	}

	authToken := &model.AuthToken{
		UserID:    student.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour * 1),
	}
	return authToken, nil
}

func (u *studentUsecase) SaveRefreshToken(userID string, refreshToken string) error {
	return u.cache.Set(userID, refreshToken, time.Hour*24*30).Err()
}

func (uc *studentUsecase) GenerateToken(studentID string) (string, error) {
	// Generate a new JWT token for the given student
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["userID"] = studentID
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (u *studentUsecase) ValidateRefreshToken(refreshToken string) (string, error) {
	// Проверяем, что refresh token является действительным
	_, err := u.cache.Get(refreshToken).Result()
	if err != nil {
		if err == redis.Nil {
			return "", errors.New("refresh token not found")
		}
		return "", err
	}

	// Разбираем токен и проверяем, что он является действительным
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		// Проверяем, что токен был подписан нашим ключом
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})

	if err != nil {
		return "", err
	}

	// Извлекаем идентификатор пользователя из токена
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid refresh token")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("invalid refresh token")
	}

	return userID, nil
}

func (u *studentUsecase) DeleteRefreshToken(userID string) error {
	return u.cache.Del(userID).Err() // удаляем ключ из Redis
}
