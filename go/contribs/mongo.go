package main

import "mongo"

var mongoColl = mongo.MongoClient.Database("contribs").Collection("go")
