package mocks

import (
	"context"

	"github.com/nurmeden/students-service/internal/app/model"
	"github.com/stretchr/testify/mock"
)

type MockStudentUsecase struct {
	mock.Mock
}

func (m *MockStudentUsecase) CreateStudent(ctx context.Context, student *model.Student) (*model.Student, error) {
	args := m.Called(ctx, student)
	return args.Get(0).(*model.Student), args.Error(1)
}

func (m *MockStudentUsecase) GetStudentByID(ctx context.Context, id string) (*model.Student, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Student), args.Error(1)
}

func (m *MockStudentUsecase) GetStudentByCoursesID(ctx context.Context, id string) (*model.Student, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Student), args.Error(1)
}

func (m *MockStudentUsecase) UpdateStudent(ctx context.Context, student_id string, student *model.Student) (*model.Student, error) {
	args := m.Called(ctx, student_id, student)
	return args.Get(0).(*model.Student), args.Error(1)
}

func (m *MockStudentUsecase) DeleteStudent(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockStudentUsecase) SignIn(signInData *model.SignInData) (*model.AuthToken, error) {
	args := m.Called(signInData)
	return args.Get(0).(*model.AuthToken), args.Error(1)
}

func (m *MockStudentUsecase) SaveRefreshToken(userID string, refreshToken string) error {
	args := m.Called(userID, refreshToken)
	return args.Error(0)
}

func (m *MockStudentUsecase) GenerateToken(studentID string) (string, error) {
	args := m.Called(studentID)
	return args.String(0), args.Error(1)
}

func (m *MockStudentUsecase) ValidateRefreshToken(refreshToken string) (string, error) {
	args := m.Called(refreshToken)
	return args.String(0), args.Error(1)
}

func (m *MockStudentUsecase) DeleteRefreshToken(userID string) error {
	args := m.Called(userID)
	return args.Error(0)
}
