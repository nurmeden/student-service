package repository

import (
	"context"
	"fmt"
	"log"
	"student-service/internal/app/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type StudentRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewStudentRepository(client *mongo.Client, dbName string, collectionName string) (*StudentRepository, error) {
	r := &StudentRepository{
		client: client,
	}

	// Получение коллекции студентов
	collection := client.Database(dbName).Collection(collectionName)
	r.collection = collection

	return r, nil
}

// Create - создание нового студента
func (r *StudentRepository) Create(ctx context.Context, student *model.Student) (*model.Student, error) {
	err := r.client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}
	_, err = r.collection.InsertOne(ctx, student)
	if err != nil {
		fmt.Println("osinda")
		return nil, fmt.Errorf("failed to create student: %v", err)
	}
	fmt.Println("bari jaksi")
	return student, nil
}

// Read - чтение информации о студенте по ID
func (r *StudentRepository) Read(ctx context.Context, id string) (*model.Student, error) {
	var student model.Student
	studentId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid id")
	}
	filter := bson.M{"_id": studentId}
	err = r.collection.FindOne(ctx, filter).Decode(&student)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Если студент не найден, возвращаем nil и ошибку nil
		}
		return nil, fmt.Errorf("failed to read student: %v", err)
	}
	return &student, nil
}

// Update - обновление информации о студенте
func (r *StudentRepository) Update(ctx context.Context, student *model.Student) (*model.Student, error) {
	filter := bson.M{"_id": student.ID}
	update := bson.M{"$set": bson.M{
		"firstName": student.FirstName,
		"lastName":  student.LastName,
		"age":       student.Age,
	}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update student: %v", err)
	}
	return student, nil
}

// Delete - удаление студента по ID
func (r *StudentRepository) Delete(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete student: %v", err)
	}
	return nil
}
