package database

import (
	"fmt"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetupDatabase() (*mongo.Client, error) {
	co := options.Client().ApplyURI(viper.GetString("MONGODB_URI"))
	client, err := mongo.NewClient(co)
	if err != nil {
		fmt.Printf("Failed to connect to MongoDB: %v", err)
		return nil, err
	}
}
