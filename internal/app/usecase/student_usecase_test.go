package usecase

import (
	"reflect"
	"testing"

	"github.com/nurmeden/students-service/internal/app/model"
	"github.com/nurmeden/students-service/internal/app/repository"
	"github.com/sirupsen/logrus"
)

func Test_studentUsecase_SignIn(t *testing.T) {
	type fields struct {
		studentRepo repository.StudentRepository
		logger      *logrus.Logger
	}
	type args struct {
		signInData *model.SignInData
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.AuthToken
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &studentUsecase{
				studentRepo: tt.fields.studentRepo,
				logger:      tt.fields.logger,
			}
			got, err := u.SignIn(tt.args.signInData)
			if (err != nil) != tt.wantErr {
				t.Errorf("studentUsecase.SignIn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("studentUsecase.SignIn() = %v, want %v", got, tt.want)
			}
		})
	}
}
