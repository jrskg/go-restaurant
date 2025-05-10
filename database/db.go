package database

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/jrskg/go-restaurant/constants"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConnectDB() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoUri := os.Getenv("MONGO_URI")
	if mongoUri == "" {
		log.Fatal("Please provide mongo uri")
	}

	client, err := mongo.Connect(options.Client().ApplyURI(mongoUri))
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Connected to database")
	return client
}

var DBClient = ConnectDB()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database(constants.DB_NAME).Collection(collectionName)
	return collection
}

func DisconnectDB(client *mongo.Client) {
	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
}
