package handler

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nurmeden/students-service/internal/app/usecase"
)

func TestStudentHandler_SignIn(t *testing.T) {
	type fields struct {
		studentUsecase usecase.StudentUsecase
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Valid credentials",
			fields: fields{
				studentUsecase: &mockStudentUsecase{},
			},
			args: args{
				c: &gin.Context{},
			},
		},
		{
			name: "Invalid email",
			fields: fields{
				studentUsecase: &mockStudentUsecase{},
			},
			args: args{
				c: &gin.Context{
					Request: &http.Request{
						Form: url.Values{
							"email":    {"invalid_email"},
							"password": {"password123"},
						},
					},
				},
			},
		},
		{
			name: "Incorrect password",
			fields: fields{
				studentUsecase: &mockStudentUsecase{},
			},
			args: args{
				c: &gin.Context{
					Request: &http.Request{
						Form: url.Values{
							"email":    {"valid_email@example.com"},
							"password": {"wrong_password"},
						},
					},
				},
			},
		},
		{
			name: "Usecase error",
			fields: fields{
				studentUsecase: &mockStudentUsecase{
					signInFunc: func(ctx context.Context, email, password string) (*domain.Student, error) {
						return nil, errors.New("usecase error")
					},
				},
			},
			args: args{
				c: &gin.Context{
					Request: &http.Request{
						Form: url.Values{
							"email":    {"valid_email@example.com"},
							"password": {"password123"},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &StudentHandler{
				studentUsecase: tt.fields.studentUsecase,
			}
			h.SignIn(tt.args.c)
		})
	}
}
