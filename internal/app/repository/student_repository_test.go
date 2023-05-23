package repository

import (
	"context"
	"reflect"
	"testing"

	"github.com/go-redis/redis"
	"github.com/nurmeden/students-service/internal/app/model"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestStudentRepository_GetStudentsByCourseID(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		collection *mongo.Collection
		cache      *redis.Client
		logger     *logrus.Logger
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*model.Student
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &StudentRepository{
				client:     tt.fields.client,
				collection: tt.fields.collection,
				cache:      tt.fields.cache,
				logger:     tt.fields.logger,
			}
			got, err := r.GetStudentsByCourseID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("StudentRepository.GetStudentsByCourseID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StudentRepository.GetStudentsByCourseID() = %v, want %v", got, tt.want)
			}
		})
	}
}
