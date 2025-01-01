package main

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func init() {
	uri := os.Getenv("MONGO_DB_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	var (
		serverAPI = options.ServerAPI(options.ServerAPIVersion1)
		opts      = options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

		ctx = context.Background()
	)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err.Error())
	}

	mongoClient = client
}
