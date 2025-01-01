package main

import "mongo"

const (
	catalogue_id = "_cat"
	licenses_id  = "_licenses"
)

var mongoColl = mongo.Client.Database("contribs").Collection("go")
