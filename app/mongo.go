package main

import (
	"context"
	"errors"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	db_apis     = "apis"
	db_contribs = "contribs"

	catalogue_id = "_cat"
	licenses_id  = "_licenses"
)

var mongoClient *mongo.Client

func init() {
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

	mongoClient = client
}

func init() {
	// TODO: Create indices for "go.contribs.locus" and "go.apis.ns/apis._id"
}

func mongoCollFromTech(tech, db string) (*mongo.Collection, error) {
	var mongoColl *mongo.Collection
	switch tech {
	case tech_go:
		mongoColl = mongoClient.Database(db).Collection("go")

	case tech_node:
		mongoColl = mongoClient.Database(db).Collection("node")

	case tech_python:
		mongoColl = mongoClient.Database(db).Collection("python")

	default:
		return nil, errors.New("can't find tech")
	}

	return mongoColl, nil
}
