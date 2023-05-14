package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nurmeden/students-service/internal/app/handlers/mocks"
	"github.com/nurmeden/students-service/internal/app/model"
	"github.com/nurmeden/students-service/internal/app/usecase"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestStudentHandler_SignIn(t *testing.T) {
	mockStudentUsecase := &mocks.MockStudentUsecase{}
	mockLogger := logrus.New()

	type fields struct {
		studentUsecase usecase.StudentUsecase
		logger         *logrus.Logger
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		expectedCode int
		expectedBody string
		mockFn       func()
		jsonBody     string
	}{
		{
			name: "Successful sign in",
			fields: fields{
				studentUsecase: mockStudentUsecase,
				logger:         mockLogger,
			},
			args: args{
				c: &gin.Context{}, // Replace with actual context
			},
			expectedCode: http.StatusOK,
			expectedBody: "{\"token\":\"auth_token\"}", // \"refresh_token\":\"refresh_token\"
			mockFn: func() {
				mockStudentUsecase.On("SignIn", mock.Anything, mock.Anything).Return(&model.AuthToken{Token: "auth_token"}, nil)
			},
			jsonBody: "{\"email\":\"test@test.com\",\"password\":\"password\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				mockStudentUsecase.AssertExpectations(t)
			})
			tt.mockFn()
			h := &StudentHandler{
				studentUsecase: tt.fields.studentUsecase,
				logger:         tt.fields.logger,
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/signin", strings.NewReader(tt.jsonBody))

			h.SignIn(c)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())
		})
	}

}

func TestStudentHandler_CreateStudent(t *testing.T) {
	mockStudentUsecase := &mocks.MockStudentUsecase{}
	// mockLogger := logrus.New()
	type fields struct {
		studentUsecase usecase.StudentUsecase
		logger         *logrus.Logger
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		expectedCode int
		expectedBody string
		mockFn       func()
		jsonBody     string
	}{
		{
			name: "Success",
			fields: fields{
				studentUsecase: mockStudentUsecase,
				logger:         logrus.New(),
			},
			args: args{
				c: &gin.Context{},
			},
			expectedCode: http.StatusCreated,
			expectedBody: "{\"_id\":\"gdfhdhd\", \"firstName\":\"Dulat\", \"lastName\":\"Nurmeden\", \"password\":\"qwerty\", \"email\":\"nurmeden@gmail.com\", \"age\":\"eht\", \"courses\":null}",
			mockFn: func() {
				mockStudentUsecase.On("CreateStudent", mock.Anything, mock.Anything).Return(nil)

				// mockStudentUsecase.On("Create", mock.Anything, mock.Anything).Return(&model.Student{}, nil)
			},
			jsonBody: "{\"email\":\"test@test.com\",\"password\":\"password\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				mockStudentUsecase.AssertExpectations(t)
			})
			tt.mockFn()
			h := &StudentHandler{
				studentUsecase: tt.fields.studentUsecase,
				logger:         tt.fields.logger,
			}
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/sign-up", strings.NewReader(tt.jsonBody))

			h.CreateStudent(c)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())
		})
	}
}
