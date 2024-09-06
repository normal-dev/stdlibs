package main

import (
	"context"
	"log"
	"runtime"
	"strings"

	"apis-go/model"

	goapis "apis-go/api"

	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	log.SetFlags(0)
}

func main() {
	ctx := context.TODO()

	log.Printf("version: %s", runtime.Version()[2:])
	_, err := mongoColl.DeleteMany(ctx, bson.M{})
	checkErr(err)

	apis := goapis.Get()
	docs := make([]any, len(apis))
	for i, api := range apis {
		docs[i] = newDoc(api)
	}
	checkErr(insertAPIs(ctx, docs))

	ns := make(map[string]struct{})
	for _, api := range apis {
		ns[api.Ns] = struct{}{}
	}

	checkErr(insertCat(ctx, ns, len(apis)))
}

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

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}
