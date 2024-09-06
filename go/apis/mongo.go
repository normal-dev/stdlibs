package main

import "mongo"

var mongoColl = mongo.MongoClient.Database("apis").Collection("go")
