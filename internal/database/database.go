package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetupDatabase(ctx context.Context, mongoURI string) (*mongo.Client, error) {
	co := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(ctx, co)
	if err != nil {
		fmt.Printf("Failed to connect to MongoDB: %v", err)
		return nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
		return nil, err
	}
	return client, nil
}
