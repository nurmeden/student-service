package repository

import (
	"context"
	"github.com/go-redis/redis"
	"github.com/nurmeden/students-service/internal/app/model"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
	"testing"
)

func TestStudentRepository_GetStudentByCoursesID(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		collection *mongo.Collection
		cache      *redis.Client
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Student
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
			}
			got, err := r.GetStudentByCoursesID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStudentByCoursesID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetStudentByCoursesID() got = %v, want %v", got, tt.want)
			}
		})
	}
}
