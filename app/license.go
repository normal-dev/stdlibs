package main

type license struct {
	ID    any `json:"_id" bson:"_id"`
	Repos []struct {
		Author string    `json:"author" bson:"author"`
		Repo   [2]string `json:"repo" bson:"repo"`
		Type   string    `json:"type" bson:"type"`
	} `json:"repos" bson:"repos"`
}
