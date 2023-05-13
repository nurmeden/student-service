package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetupDatabase() (*mongo.Client, error) {
	co := options.Client().ApplyURI("mongodb://studentsdb:27017")
	client, err := mongo.NewClient(co)
	if err != nil {
		fmt.Printf("Failed to connect to MongoDB: %v", err)
		return nil, err
	}
	if err != nil {
		fmt.Printf("Failed to connect to MongoDB: %v", err)
		return nil, err
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("FFFFFFFFFFFFFFFFFFFFFF", err)
		return nil, err
	}
	fmt.Printf("client: %v\n", client)
	return client, nil
}
