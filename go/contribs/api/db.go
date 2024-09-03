package api

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DB_APIs     = "apis"
	DB_CONTRIBS = "contribs"
)

var MongoClient *mongo.Client

func init() {
	log.SetFlags(0)
	log.Default().SetOutput(os.Stderr)
}

func init() {
	log.Println("connecting to MongoDB...")

	uri := os.Getenv("MONGO_DB_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		panic(err.Error())
	}

	MongoClient = client
}
