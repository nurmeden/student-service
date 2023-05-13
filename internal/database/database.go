package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetupDatabase() *mongo.Client {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://studentsdb:27017"))
	if err != nil {
		fmt.Printf("Failed to connect to MongoDB: %v", err)
		return nil
	}
	if err != nil {
		fmt.Printf("Failed to connect to MongoDB: %v", err)
		return nil
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("FFFFFFFFFFFFFFFFFFFFFF", err)
		return nil
	}
	fmt.Printf("client: %v\n", client)
	return client
}
