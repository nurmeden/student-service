package repository

import (
	"context"
	"reflect"
	"testing"

	"github.com/go-redis/redis"
	"github.com/nurmeden/students-service/internal/app/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func TestStudentRepository_Create(t *testing.T) {
	type fields struct {
		client     *mongo.Client
		collection *mongo.Collection
		cache      *redis.Client
	}
	type args struct {
		ctx     context.Context
		student *model.Student
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
			got, err := r.Create(tt.args.ctx, tt.args.student)
			if (err != nil) != tt.wantErr {
				t.Errorf("StudentRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StudentRepository.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func createTestDBConnection(t *testing.T) (*mongo.Client, *mongo.Collection) {
	// Connect to MongoDB
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("Error connecting to MongoDB: %v", err)
	}

	// Create a new database and collection for testing
	dbName := "testdb"
	collName := "students"
	err = client.Ping(ctx, nil)
	if err != nil {
		t.Fatalf("Error pinging MongoDB: %v", err)
	}
	db := client.Database(dbName)
	collection := db.Collection(collName)

	// Clean up any existing data in the test collection
	_, err = collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		t.Fatalf("Error cleaning up test data: %v", err)
	}

	return client, collection
}
