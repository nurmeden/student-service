package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetupDatabase(ctx context.Context) (*mongo.Client, error) {
	co := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, co)
	if err != nil {
		fmt.Printf("Failed to connect to MongoDB: %v", err)
		return nil, err
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
		return nil, err
	}
	return client, nil
}
