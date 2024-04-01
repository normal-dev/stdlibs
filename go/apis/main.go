package main

import (
	"context"
	"log"
	"os"
	"runtime"
	"strings"

	goapis "apis-go/api"

	"apis-go/model"

	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	log.SetFlags(0)
	log.Default().SetOutput(os.Stderr)
}

var mongoColl = goapis.MongoClient.Database("apis").Collection("go")

func main() {
	log.Printf("using Go version %s", runtime.Version()[2:])

	log.Println("cleaning...")
	checkErr(clean(context.TODO()))

	apis := goapis.Get()
	log.Printf("saving %d apis...", len(apis))
	for _, api := range apis {
		// TODO: Save all at once
		checkErr(saveAPI(context.TODO(), api))
	}

	ns := make(map[string]struct{})
	for _, api := range apis {
		ns[api.Ns] = struct{}{}
	}
	log.Printf("saving catalogue...")
	checkErr(saveCat(context.TODO(), ns, len(apis)))
}

func saveAPI(ctx context.Context, api goapis.API) error {
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
	_, err := mongoColl.InsertOne(ctx, doc)
	return err
}

func saveCat(ctx context.Context, nss map[string]struct{}, napis int) error {
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
		Vids: map[string]string{
			"archive/tar": "FBoHtOuFnHY",
			"errors":      "SoGSHWe28D0",
			"fmt":         "uuDo2S8qbcc",
		},
	})
	return err
}

func clean(ctx context.Context) error {
	_, err := mongoColl.DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}
	return err
}

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func logErr(err error) {
	if err != nil {
		log.Println(err.Error())
	}
}
