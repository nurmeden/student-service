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
	"go.mongodb.org/mongo-driver/mongo/options"
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
	result, err := r.collection.InsertOne(ctx, student)
	if err != nil {
		r.logger.Errorf("Failed to create student: %v", err)
		return nil, fmt.Errorf("failed to create student: %v", err)
	}
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("failed to convert InsertedID to ObjectID")
	}
	insertedIDStr := insertedID.Hex()

	studentJSON, err := json.Marshal(student)
	if err != nil {
		r.logger.Errorf("Failed to marshal student data for caching: %v", err)
		return nil, err
	}

	err = r.cache.Set(insertedIDStr, studentJSON, 0).Err()
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

func (r *StudentRepository) GetStudentsByCourseID(ctx context.Context, id string) ([]*model.Student, error) {
	cachedResult, err := r.cache.Get(id).Result()
	fmt.Printf("id in the course id: %v\n", id)
	if err == nil {
		students := []*model.Student{}
		fmt.Printf("students in the err: %v\n", students)
		err = json.Unmarshal([]byte(cachedResult), &students)
		if err != nil {
			r.logger.Errorf("Error unmarshalling cached result for students with course ID %s: %s", id, err)
			return nil, err
		}
		return students, nil
	}

	var students []*model.Student

	filter := bson.M{"courses": bson.M{"$in": []string{id}}}
	cursor, err := r.collection.Find(ctx, filter)
	fmt.Printf("cursor: %v\n", cursor)
	if err != nil {
		return nil, fmt.Errorf("failed to query students: %v", err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var student model.Student
		if err := cursor.Decode(&student); err != nil {
			r.logger.Errorf("Error decoding student: %s", err)
			continue
		}
		fmt.Printf("student following: %v\n", student)
		students = append(students, &student)
	}
	fmt.Printf("students: %v\n", students)

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	if len(students) == 0 {
		r.logger.Infof("No students found for course ID %s in database", id)
		return nil, nil
	}

	studentsJSON, err := json.Marshal(students)
	if err != nil {
		r.logger.Errorf("Failed to marshal students data for caching: %v", err)
		return nil, err
	}

	err = r.cache.Set(id, studentsJSON, 0).Err()
	if err != nil {
		r.logger.Errorf("Failed to cache students data: %v", err)
		return nil, err
	}

	return students, nil
}

func (r *StudentRepository) UpdateStudents(ctx context.Context, student *model.Student, studentID primitive.ObjectID) (*model.Student, error) {
	filter := bson.M{"_id": studentID}
	update := bson.M{"$set": bson.M{
		"firstName": student.FirstName,
		"lastName":  student.LastName,
		"age":       student.Age,
		"courses":   student.Courses,
	}}
	fmt.Printf("student.Courses: %v\n", student.Courses)
	fmt.Printf("student.ID: %v\n", studentID)

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	result := r.collection.FindOneAndUpdate(ctx, filter, update, opts)
	if err := result.Err(); err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("student not found")
		}
		return nil, err
	}

	var updatedStudent *model.Student
	if err := result.Decode(&updatedStudent); err != nil {
		return nil, err
	}

	studentJSON, err := json.Marshal(updatedStudent)
	if err != nil {
		return nil, err
	}

	err = r.cache.Set(studentID.Hex(), studentJSON, 0).Err()
	if err != nil {
		fmt.Printf("Failed to update student data in Redis: %v\n", err)
	}

	return updatedStudent, nil
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

func (r *StudentRepository) CheckEmailExistence(ctx context.Context, email string) (bool, error) {
	filter := bson.M{"email": email}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
