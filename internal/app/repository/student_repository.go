package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/nurmeden/students-service/internal/app/model"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type StudentRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
	cache      *redis.Client
	logger     *logrus.Logger
}

func NewStudentRepository(client *mongo.Client, dbName string, collectionName string, cache *redis.Client, logger *logrus.Logger) (*StudentRepository, error) {
	r := &StudentRepository{
		client: client,
		cache:  cache,
		logger: logger,
	}

	collection := client.Database(dbName).Collection(collectionName)
	r.collection = collection

	return r, nil
}

func (r *StudentRepository) CreateStudent(ctx context.Context, student *model.Student) (*model.Student, error) {
	r.logger.Infof("Creating new student: %+v", student)

	err := r.client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}
	_, err = r.collection.InsertOne(ctx, student)
	if err != nil {
		r.logger.Errorf("Failed to create student: %v", err)
		return nil, fmt.Errorf("failed to create student: %v", err)
	}

	studentJSON, err := json.Marshal(student)
	if err != nil {
		r.logger.Errorf("Failed to marshal student data for caching: %v", err)
		return nil, err
	}

	err = r.cache.Set(student.ID, studentJSON, 0).Err()
	if err != nil {
		r.logger.Errorf("Failed to cache student data: %v", err)
		return nil, err
	}

	r.logger.Infof("Student created successfully")

	return student, nil
}

func (r *StudentRepository) GetStudentByID(ctx context.Context, id string) (*model.Student, error) {
	cachedResult, err := r.cache.Get(id).Result()
	if err == nil {
		student := &model.Student{}
		err = json.Unmarshal([]byte(cachedResult), student)
		if err != nil {
			r.logger.Errorf("Error unmarshalling cached result for student with ID %s: %s", id, err)
			return nil, err
		}
		return student, nil
	} else if err != redis.Nil {
		r.logger.Errorf("Error getting cached result for student with ID %s: %s", id, err)
		return nil, err
	}

	var student model.Student
	studentId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Errorf("Invalid ID")
	}
	filter := bson.M{"_id": studentId}
	err = r.collection.FindOne(ctx, filter).Decode(&student)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			r.logger.Infof("Student with id %s not found in database", id)
			return nil, nil // Если студент не найден, возвращаем nil и ошибку nil
		}
		r.logger.Errorf("failed to read student: %v", err)
		return nil, fmt.Errorf("failed to read student: %v", err)
	}

	studentJSON, err := json.Marshal(student)
	if err != nil {
		r.logger.Errorf("Failed to marshal student data for caching: %v", err)
		return nil, err
	}

	err = r.cache.Set(id, studentJSON, 0).Err()
	if err != nil {
		r.logger.Errorf("Failed to cache student data: %v", err)
		return nil, err
	}

	return &student, nil
}

func (r *StudentRepository) GetStudentByCoursesID(ctx context.Context, id string) (*model.Student, error) {
	cachedResult, err := r.cache.Get(id).Result()
	if err == nil {
		student := &model.Student{}
		err = json.Unmarshal([]byte(cachedResult), student)
		if err != nil {
			r.logger.Errorf("Error unmarshalling cached result for student with course ID %s: %s", id, err)
			return nil, err
		}
		return student, nil
	}

	var student model.Student

	filter := bson.M{"courses": id}
	err = r.collection.FindOne(ctx, filter).Decode(&student)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			r.logger.Infof("Student with coursesId %s not found in database", id)
			return nil, nil // Если студент не найден, возвращаем nil и ошибку nil
		}
		return nil, fmt.Errorf("failed to read student: %v", err)
	}

	studentJSON, err := json.Marshal(student)
	if err != nil {
		r.logger.Errorf("Failed to marshal student data for caching: %v", err)
		return nil, err
	}

	err = r.cache.Set(id, studentJSON, 0).Err()
	if err != nil {
		r.logger.Errorf("Failed to cache student data: %v", err)
		return nil, err
	}
	return &student, nil
}

func (r *StudentRepository) UpdateStudents(ctx context.Context, student *model.Student) (*model.Student, error) {
	filter := bson.M{"_id": student.ID}
	update := bson.M{"$set": bson.M{
		"firstName": student.FirstName,
		"lastName":  student.LastName,
		"age":       student.Age,
	}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Errorf("failed to update student: %v", err)
		return nil, fmt.Errorf("failed to update student: %v", err)
	}
	return student, nil
}

func (r *StudentRepository) Delete(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete student: %v", err)
	}
	return nil
}

func (r *StudentRepository) GetByEmail(ctx context.Context, email string) (*model.Student, error) {
	filter := bson.M{"email": email}

	var student model.Student
	err := r.collection.FindOne(ctx, filter).Decode(&student)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}
	fmt.Printf("student: %v\n", student)

	return &student, nil
}

func (r *StudentRepository) GetByEmailForSignUp(ctx context.Context, email string) *model.Student {
	filter := bson.M{"email": email}

	var student model.Student
	err := r.collection.FindOne(ctx, filter).Decode(&student)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		return nil
	}
	fmt.Printf("student: %v\n", student)

	return &student
}
