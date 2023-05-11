package db

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Define a struct for the file metadata
type FileMetadata struct {
	ID        string    `bson:"_id"`
	Container string    `bson:"container_name"`
	Folder    string    `bson:"folder"`
	FileName  string    `bson:"file_name"`
	Created   time.Time `bson:"created"`
}

func getMongoUrl() string {
	url := os.Getenv("MONGODB_URL")
	if url == "" {
		url = "mongodb://localhost:27017"
	}
	return url
}

func getMongoDatabaseName() string {
	databaseName := os.Getenv("MONGODB_DATABASE")
	if databaseName == "" {
		databaseName = "file-storage-service"
	}
	return databaseName
}

func getMongoCollectionName() string {
	collectionName := os.Getenv("MONGODB_COLLECTION")
	if collectionName == "" {
		collectionName = "file-metadata"
	}
	return collectionName
}

func SaveFileMetadata(metadata FileMetadata) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(getMongoUrl())
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	collection := client.Database(getMongoDatabaseName()).Collection(getMongoCollectionName())
	return collection.InsertOne(ctx, metadata)
}
