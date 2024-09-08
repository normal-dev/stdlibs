package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	inaugural = time.Date(
		2023,
		7,
		25,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	today = time.Now().UTC()
)

type (
	repository struct {
		Tech      string    `json:"tech" bson:"tech"`
		RepoName  string    `json:"repo_name" bson:"repo_name"`
		RepoOwner string    `json:"repo_owner" bson:"repo_owner"`
		Locusn    int       `json:"locusn" bson:"locusn"`
		Updated   time.Time `json:"updated" bson:"updated"`
	}

	contribution struct {
		Tech string `json:"tech" bson:"tech"`
		Repo struct {
			RepoOwner string `json:"repo_owner" bson:"repo_owner"`
			RepoName  string `json:"repo_name" bson:"repo_name"`
		}
		Locusn int `json:"locusn" bson:"locusn"`
	}
)

func (c repository) EqRepo(repoOwner, repoName string) bool {
	return c.RepoOwner == repoOwner && c.RepoName == repoName
}

func main() {
	ctx := context.Background()
	coll := mongoClient.Database("app").Collection("repos")

	cur, err := coll.Find(ctx, bson.D{})
	if err != nil {
		panic(err.Error())
	}
	defer cur.Close(ctx)
	var repos []repository
	if err := cur.All(ctx, &repos); err != nil {
		panic(err)
	}

	contribs, err := fetchContribs(ctx)
	if err != nil {
		panic(err.Error())
	}

OuterLoop:
	for _, contrib := range contribs {
		var (
			repoOwner = contrib.Repo.RepoOwner
			repoName  = contrib.Repo.RepoName
		)

		// Update
		for _, repo := range repos {
			if !repo.EqRepo(repoOwner, repoName) {
				continue
			}

			if repo.Locusn != contrib.Locusn {
				repo.Locusn = contrib.Locusn
				repo.Updated = today.UTC()
			}
			_, err := coll.UpdateOne(ctx,
				bson.D{
					{Key: "repo_owner", Value: repo.RepoOwner},
					{Key: "repo_name", Value: repo.RepoName},
				},
				bson.D{
					{Key: "$set", Value: repo},
				},
			)
			if err != nil {
				panic(err.Error())
			}

			continue OuterLoop
		}

		// Insert
		_, err := coll.InsertOne(ctx, repository{
			RepoName:  repoName,
			RepoOwner: repoOwner,
			Locusn:    contrib.Locusn,
			Tech:      contrib.Tech,
			Updated:   today.UTC(),
		})
		if err != nil {
			panic(err.Error())
		}
	}
}

func fetchContribs(ctx context.Context) ([]contribution, error) {
	var contribs []contribution

	f := func(ctx context.Context, tech string) {
		mongoColl := mongoClient.Database("contribs").Collection(tech)

		pipeline := mongo.Pipeline{
			bson.D{
				{
					Key: "$match", Value: bson.D{
						{
							Key: "_id", Value: bson.D{
								{
									Key: "$nin", Value: bson.A{
										"_licenses",
										"_cat",
									},
								},
							},
						},
					},
				},
			},
			bson.D{
				{
					Key: "$sort", Value: bson.D{
						{Key: "_id", Value: 1},
					},
				},
			},
			bson.D{
				{
					Key: "$group", Value: bson.D{
						{
							Key: "_id", Value: bson.D{
								{Key: "repo_name", Value: "$repo_name"},
								{Key: "repo_owner", Value: "$repo_owner"},
							},
						},
						{
							Key: "locusn", Value: bson.D{
								{
									Key: "$sum", Value: bson.D{
										{Key: "$size", Value: "$locus"},
									},
								},
							},
						},
					},
				},
			},
			bson.D{
				{
					Key: "$project", Value: bson.D{
						{Key: "_id", Value: 0},
						{Key: "repo", Value: "$_id"},
						{Key: "locusn", Value: 1},
					},
				},
			},
		}
		cur, err := mongoColl.Aggregate(ctx, pipeline)
		if err != nil {
			panic(err.Error())
		}
		defer cur.Close(ctx)

		var cs []contribution
		if err := cur.All(ctx, &cs); err != nil {
			panic(err.Error())
		}

		for idx := range cs {
			cs[idx].Tech = tech
		}

		contribs = append(contribs, cs...)
	}

	f(ctx, "go")
	f(ctx, "node")

	return contribs, nil
}
