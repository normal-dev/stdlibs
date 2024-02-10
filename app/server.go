package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/normal-dev/stdlibs/model"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	db_apis     = "apis"
	db_contribs = "contribs"
)

var mongoClient *mongo.Client

func init() {
	log.Println("connecting to MongoDB...")

	uri := os.Getenv("MONGO_DB_URI")
	if uri == "" {
		log.Printf("can't find MongoDB URI, falling back to %s", "mongodb://localhost:27017")
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
	log.SetFlags(0)
	log.Default().SetOutput(os.Stderr)
}

const (
	tech_go     = "go"
	tech_node   = "node"
	tech_python = "python"
)

func main() {
	noClient := flag.Bool("no-client", false, "")
	flag.Parse()

	router := gin.Default()

	// Trust all proxies
	router.SetTrustedProxies(nil)

	// CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"https://stdlibs-app-m3pnp47eca-uc.a.run.app",
			"https://www.stdlibs.com",
			"https://stdlibs.com",
		},
		AllowMethods:     []string{"GET", "OPTIONS", "HEAD"},
		AllowCredentials: false,
		MaxAge:           24 * time.Hour,
	}))

	// Resolve UTF-8 characters
	router.UseRawPath = true
	router.UnescapePathValues = false

	// Website/client
	if wantsClient := !fromPtr(noClient); wantsClient {
		log.Print("loading client assets...")
		router.Static("/assets", "./website/assets")
		router.LoadHTMLGlob("website/index.html")
		router.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", nil)
		})
	}

	// Cache pages for one hour
	store := persistence.NewInMemoryStore(time.Hour * 6)

	// /node
	router.GET("/api/:tech", cache.CachePage(store, time.Hour*12, func(ctx *gin.Context) {
		// Union type of API and contributions catalogue
		type cat struct {
			NAPIs     int      `json:"n_apis" bson:"n_apis"`
			NContribs int      `json:"n_contribs" bson:"n_contribs"`
			NNs       int      `json:"n_ns" bson:"n_ns"`
			NRepos    int      `json:"n_repos" bson:"n_repos"`
			Ns        []string `json:"ns" bson:"ns"`
			Version   string   `json:"version" bson:"version"`
		}

		var c cat

		{
			var (
				err       error
				mongoColl *mongo.Collection
			)
			mongoColl, err = mongoCollFromCtx(ctx, db_apis)
			if err != nil {
				log.Println(err.Error())
				ctx.Status(http.StatusBadRequest)
				return
			}

			err = mongoColl.FindOne(context.TODO(), bson.D{{Key: "_id", Value: model.CAT_ID}}).Decode(&c)
			if err != nil {
				// TODO: Check for not found error
				log.Println(err.Error())
				ctx.Status(http.StatusInternalServerError)
				return
			}
		}
		{
			var (
				err       error
				mongoColl *mongo.Collection
			)
			mongoColl, err = mongoCollFromCtx(ctx, db_contribs)
			if err != nil {
				log.Println(err.Error())
				ctx.Status(http.StatusBadRequest)
				return
			}

			err = mongoColl.FindOne(context.TODO(), bson.D{{Key: "_id", Value: model.CAT_ID}}).Decode(&c)
			if err != nil {
				// TODO: Check for not found error
				log.Println(err.Error())
				ctx.Status(http.StatusInternalServerError)
				return
			}
		}

		ctx.JSON(http.StatusOK, c)
	}))
	// /go/gen
	router.GET("/api/gen", cache.CachePage(store, time.Hour*6, func(ctx *gin.Context) {
		const maxcontribs = 3

		contribs := make([]bson.M, 0)
		// Go
		{
			mongoColl := mongoClient.Database(db_contribs).Collection("go")
			size := rand.Intn(6-3) + 3 // 3-6
			filter := bson.D{
				{Key: "apis", Value: bson.D{
					{Key: "$size", Value: size},
				}},
			}
			ncontribs, err := mongoColl.CountDocuments(context.TODO(), filter)
			if err != nil {
				log.Println(err.Error())
				ctx.Status(http.StatusInternalServerError)
				return
			}
			skip := rand.Int63n(ncontribs)
			cur, err := mongoColl.Find(context.TODO(), filter, &options.FindOptions{
				Limit: toPtr[int64](maxcontribs),
				Skip:  toPtr[int64](skip),
			})
			if err != nil {
				log.Println(err.Error())
				ctx.Status(http.StatusInternalServerError)
				return
			}
			for cur.Next(context.TODO()) {
				var contrib primitive.M
				if err := cur.Decode(&contrib); err != nil {
					log.Println(err.Error())
					ctx.Status(http.StatusInternalServerError)
					return
				}
				contribs = append(contribs, contrib)
			}
		}

		// Node.js
		{
			mongoColl := mongoClient.Database(db_contribs).Collection("node")
			filter := bson.D{
				{Key: "apis", Value: bson.D{
					{Key: "$size", Value: 5},
				}},
			}
			ncontribs, err := mongoColl.CountDocuments(context.TODO(), filter)
			if err != nil {
				log.Println(err.Error())
				ctx.Status(http.StatusInternalServerError)
				return
			}
			skip := rand.Int63n(ncontribs)
			cur, err := mongoColl.Find(context.TODO(), filter, &options.FindOptions{
				Limit: toPtr[int64](maxcontribs),
				Skip:  toPtr[int64](skip),
			})
			if err != nil {
				log.Println(err.Error())
				ctx.Status(http.StatusInternalServerError)
				return
			}
			for cur.Next(context.TODO()) {
				var contrib primitive.M
				if err := cur.Decode(&contrib); err != nil {
					log.Println(err.Error())
					ctx.Status(http.StatusInternalServerError)
					return
				}
				contribs = append(contribs, contrib)
			}
		}

		// Shuffle contributions
		for i := range contribs {
			j := rand.Intn(i + 1)
			contribs[i], contribs[j] = contribs[j], contribs[i]
		}

		ctx.JSON(http.StatusOK, contribs)
	}))
	// /node/licenses
	router.GET("/api/:tech/licenses", cache.CachePage(store, time.Hour*12, func(ctx *gin.Context) {
		var (
			err       error
			mongoColl *mongo.Collection
		)
		mongoColl, err = mongoCollFromCtx(ctx, db_contribs)
		if err != nil {
			log.Println(err.Error())
			ctx.Status(http.StatusBadRequest)
			return
		}

		var licenses model.Licenses
		err = mongoColl.FindOne(context.TODO(), bson.D{{Key: "_id", Value: model.LICENSES_ID}}).Decode(&licenses)
		if err != nil {
			log.Println(err.Error())
			ctx.Status(http.StatusInternalServerError)
			return
		}
		ctx.JSON(http.StatusOK, licenses)
	}))
	// /go/context
	router.GET("/api/:tech/:ns", cache.CachePage(store, time.Hour*12, func(ctx *gin.Context) {
		var (
			err       error
			mongoColl *mongo.Collection
		)
		mongoColl, err = mongoCollFromCtx(ctx, db_apis)
		if err != nil {
			log.Println(err.Error())
			ctx.Status(http.StatusBadRequest)
			return
		}

		ns := ctx.Param("ns")
		ns, err = url.QueryUnescape(ns)
		if err != nil {
			log.Println(err.Error())
			ctx.Status(http.StatusBadRequest)
			return
		}

		filter := bson.D{
			{Key: "ns", Value: ns},
			{Key: "_id", Value: bson.D{primitive.E{Key: "$ne", Value: model.CAT_ID}}},
		}
		cur, err := mongoColl.Find(context.TODO(), filter)
		if err != nil {
			log.Println(err.Error())
			ctx.Status(http.StatusInternalServerError)
			return
		}
		apis := make([]bson.M, 0)
		if err := cur.All(context.TODO(), &apis); err != nil {
			log.Println(err.Error())
			ctx.Status(http.StatusInternalServerError)
			return
		}
		ctx.JSON(http.StatusOK, apis)
	}))
	// /node/crypto/verify
	router.GET("/api/:tech/:ns/:api", func(ctx *gin.Context) {
		var (
			err       error
			mongoColl *mongo.Collection
		)
		mongoColl, err = mongoCollFromCtx(ctx, db_contribs)
		if err != nil {
			log.Println(err.Error())
			ctx.Status(http.StatusBadRequest)
			return
		}

		ns := ctx.Param("ns")
		ns, err = url.QueryUnescape(ns)
		if err != nil {
			log.Println(err.Error())
			ctx.Status(http.StatusBadRequest)
			return
		}
		api := ctx.Param("api")
		api, err = url.QueryUnescape(api)
		if err != nil {
			log.Println(err.Error())
			ctx.Status(http.StatusBadRequest)
			return
		}

		filter := bson.M{"apis.ident": fmt.Sprintf("%s.%s", ns, api)}
		var perPage int64 = 6
		page, err := strconv.ParseInt(ctx.Query("page"), 10, 64)
		if err != nil {
			log.Println(err.Error())
			ctx.Status(http.StatusBadRequest)
			return
		}
		skip := page*perPage - perPage
		opts := &options.FindOptions{
			Limit: &perPage,
			Skip:  &skip,
		}
		cur, err := mongoColl.Find(context.TODO(), filter, opts)
		if err != nil {
			log.Println(err.Error())
			ctx.Status(http.StatusInternalServerError)
			return
		}
		contribs := make([]bson.M, 0)
		if err := cur.All(context.TODO(), &contribs); err != nil {
			log.Println(err.Error())
			ctx.Status(http.StatusInternalServerError)
			return
		}

		contribsn, err := mongoColl.CountDocuments(context.TODO(), filter)
		if err != nil {
			log.Println(err.Error())
			ctx.Status(http.StatusInternalServerError)
			return
		}
		type pagination struct {
			Contribs []bson.M `json:"contribs"`
			Total    int64    `json:"total"`
			PerPage  int64    `json:"per_page"`
		}
		p := pagination{contribs, contribsn, perPage}

		ctx.JSON(http.StatusOK, p)
	})

	// Any other route
	router.NoRoute(func(ctx *gin.Context) {
		// Route non-API requests to website
		if !strings.HasPrefix(ctx.Request.RequestURI, "/api") {
			ctx.Status(http.StatusNotFound)
			ctx.File("./website/index.html")
		}
	})

	// Favicon
	router.StaticFile("/favicon.png", "./website/favicon.png")
	router.StaticFile("/sitemap.xml", "./website/sitemap.xml")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)
	router.Run(addr)
}

func mongoCollFromCtx(ctx *gin.Context, db string) (*mongo.Collection, error) {
	var mongoColl *mongo.Collection
	switch ctx.Param("tech") {
	case tech_go:
		mongoColl = mongoClient.Database(db).Collection("go")

	case tech_node:
		mongoColl = mongoClient.Database(db).Collection("node")

	case tech_python:
		return nil, errors.New("not implemented")

	default:
		return nil, errors.New("can't find tech")
	}

	return mongoColl, nil
}

func fromPtr[T any](v *T) T { return *v }

func toPtr[T any](v T) *T { return &v }
