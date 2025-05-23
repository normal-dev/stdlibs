package main

import (
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

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/normal-dev/stdlibs/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	log.SetFlags(0)
}

const (
	tech_go     = "go"
	tech_node   = "node"
	tech_python = "python"
)

func main() {
	// Flags
	// Don't use built client. This should be used during development
	noClient := flag.Bool("no-client", false, "")
	flag.Parse()

	router := gin.Default()

	// HTTP compression
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	// Trust all proxies
	_ = router.SetTrustedProxies(nil)

	// CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"https://www.stdlibs.com",
			"https://stdlibs.com",
		},
		AllowMethods:     []string{"GET", "OPTIONS", "HEAD"},
		AllowCredentials: false,
		MaxAge:           24 * time.Hour,
	}))

	// Router path defaults
	router.UseRawPath = true
	router.UnescapePathValues = false

	// Eventually serve client
	if ok := fromPtr(noClient); !ok {
		router.Static("/assets", "./website/assets")
		router.LoadHTMLGlob("./website/index.html")
		router.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", nil)
		})
	}

	// Cache pages
	store := persistence.NewInMemoryStore(time.Hour * 1)

	// SEO repositories
	router.GET("/api/seo/repositories", cache.CachePage(store, time.Hour*48, func(ctx *gin.Context) {
		mongoColl := mongoClient.Database("seo").Collection("repos")

		var repos []bson.M
		cur, err := mongoColl.Find(ctx, bson.D{})
		if err != nil {
			log.Println(err.Error())
			ctx.Status(http.StatusInternalServerError)
			return
		}

		if err := cur.All(ctx, &repos); err != nil {
			log.Println(err.Error())
			ctx.Status(http.StatusInternalServerError)
			return
		}

		ctx.JSON(http.StatusOK, repos)
	}))

	// Generator returns random contributions, e. g. "/go/gen"
	router.GET("/api/gen", cache.CachePage(store, time.Hour*6, func(ctx *gin.Context) {
		const maxcontribs = 1

		contribs := make([]bson.M, 0)
		// Go
		{
			mongoColl := mongoClient.Database(db_contribs).Collection("go")

			size := rand.Intn(6-3) + 3 // 3-6
			filter := bson.M{"locus": bson.M{"$size": size}}
			ncontribs, err := mongoColl.CountDocuments(ctx, filter)
			if err != nil {
				log.Println(err.Error())
				ctx.Status(http.StatusInternalServerError)
				return
			}

			skip := rand.Int63n(ncontribs)
			cur, err := mongoColl.Find(ctx, filter, &options.FindOptions{
				Limit: toPtr[int64](maxcontribs),
				Skip:  toPtr(skip),
			})
			if err != nil {
				log.Println(err.Error())
				ctx.Status(http.StatusInternalServerError)
				return
			}
			for cur.Next(ctx) {
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

			size := rand.Intn(6-3) + 3 // 3-6
			filter := bson.M{"locus": bson.M{"$size": size}}
			ncontribs, err := mongoColl.CountDocuments(ctx, filter)
			if err != nil {
				log.Println(err.Error())
				ctx.Status(http.StatusInternalServerError)
				return
			}

			skip := rand.Int63n(ncontribs)
			cur, err := mongoColl.Find(ctx, filter, &options.FindOptions{
				Limit: toPtr[int64](maxcontribs),
				Skip:  toPtr[int64](skip),
			})
			if err != nil {
				log.Println(err.Error())
				ctx.Status(http.StatusInternalServerError)
				return
			}
			for cur.Next(ctx) {
				var contrib primitive.M
				if err := cur.Decode(&contrib); err != nil {
					log.Println(err.Error())
					ctx.Status(http.StatusInternalServerError)
					return
				}
				contribs = append(contribs, contrib)
			}
		}

		// TODO: Add Python

		// Shuffle contributions
		for i := range contribs {
			j := rand.Intn(i + 1)
			contribs[i], contribs[j] = contribs[j], contribs[i]
		}

		ctx.JSON(http.StatusOK, contribs)
	}))

	// Licenses, e. g. "/node/licenses"
	router.GET("/api/:tech/licenses", cache.CachePage(store, time.Hour*3, func(ctx *gin.Context) {
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

		var lic model.License
		err = mongoColl.FindOne(ctx, bson.D{
			{Key: "_id", Value: licenses_id},
		}).Decode(&lic)
		if err != nil {
			log.Println(err.Error())
			ctx.Status(http.StatusInternalServerError)
			return
		}

		ctx.JSON(http.StatusOK, lic)
	}))

	// Catalogue/namespaces, e. g. "/node"
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

			err = mongoColl.FindOne(ctx, bson.D{
				{Key: "_id", Value: catalogue_id},
			}).Decode(&c)
			if err != nil {
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

			err = mongoColl.FindOne(ctx, bson.D{
				{Key: "_id", Value: catalogue_id},
			}).Decode(&c)
			if err != nil {
				log.Println(err.Error())
				ctx.Status(http.StatusInternalServerError)
				return
			}
		}

		ctx.JSON(http.StatusOK, c)
	}))

	// APIs, e. g. "/go/context"
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
			{
				Key: "_id", Value: bson.D{
					primitive.E{Key: "$ne", Value: catalogue_id},
				},
			},
		}
		cur, err := mongoColl.Find(ctx, filter)
		if err != nil {
			log.Println(err.Error())
			ctx.Status(http.StatusInternalServerError)
			return
		}
		apis := make([]bson.M, 0)
		if err := cur.All(ctx, &apis); err != nil {
			log.Println(err.Error())
			ctx.Status(http.StatusInternalServerError)
			return
		}
		ctx.JSON(http.StatusOK, apis)
	}))

	// Contributions, e. g. "/node/crypto/verify"
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

		filter := bson.M{"locus.ident": fmt.Sprintf("%s.%s", ns, api)}
		page, err := strconv.ParseInt(ctx.Query("page"), 10, 64)
		if err != nil {
			log.Println(err.Error())
			ctx.Status(http.StatusBadRequest)
			return
		}
		var perPage int64 = 6
		opts := &options.FindOptions{
			Limit: toPtr(perPage),
			Skip:  toPtr(page*perPage - perPage),
		}
		cur, err := mongoColl.Find(ctx, filter, opts)
		if err != nil {
			log.Println(err.Error())
			ctx.Status(http.StatusInternalServerError)
			return
		}
		contribs := make([]bson.M, 0)
		if err := cur.All(ctx, &contribs); err != nil {
			log.Println(err.Error())
			ctx.Status(http.StatusInternalServerError)
			return
		}

		contribsn, err := mongoColl.CountDocuments(ctx, filter)
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
	// Sitemap
	router.StaticFile("/sitemap.xml", "./website/sitemap.xml")

	// Address and port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)

	_ = router.Run(addr)
}

func mongoCollFromCtx(ctx *gin.Context, db string) (*mongo.Collection, error) {
	return mongoCollFromTech(ctx.Param("tech"), db)
}

func fromPtr[T any](v *T) T { return *v }

func toPtr[T any](v T) *T { return &v }
