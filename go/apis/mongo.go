package main

import (
	"context"
	"mongo"
	"runtime"
	"strings"

	goapis "apis-go/api"
	"apis-go/model"

	"go.mongodb.org/mongo-driver/bson"
)

var mongoColl = mongo.Client.Database("apis").Collection("go")

func newDoc(api goapis.API) bson.D {
	doc := bson.D{
		bson.E{Key: "_id", Value: api.ID()},
		bson.E{Key: "doc", Value: api.Doc},
		bson.E{Key: "name", Value: api.Name},
		bson.E{Key: "type", Value: api.Type},
		bson.E{Key: "ns", Value: api.Ns},
	}
	if api.Value != nil {
		doc = append(doc, bson.E{Key: "value", Value: *api.Value})
	}
	return doc
}

func insertAPIs(ctx context.Context, docs []any) error {
	_, err := mongoColl.InsertMany(ctx, docs)
	return err
}

func insertCat(ctx context.Context, nss map[string]struct{}, napis int) error {
	var ns []string
	for pkg := range nss {
		ns = append(ns, pkg)
	}
	_, err := mongoColl.InsertOne(ctx, model.Cat{
		ID:      model.CAT_ID,
		NAPIs:   napis,
		NNs:     len(nss),
		Ns:      ns,
		Version: strings.TrimPrefix(runtime.Version(), "go"),
		Vids:    map[string]string{},
	})
	return err
}
