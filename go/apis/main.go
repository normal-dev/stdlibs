package main

import (
	"context"
	"log"
	"runtime"

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

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}
