package main

import (
	"contribs-go/model"
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
	noClient := *flag.Bool("no-client", false, "")
	flag.Parse()

	router := gin.Default()

	// Compression
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	// Trust all proxies
	_ = router.SetTrustedProxies(nil)

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
	if wantsClient := noClient; wantsClient {
		router.Static("/assets", "./website/assets")
		router.LoadHTMLGlob("website/index.html")
		router.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", nil)
		})
	}

	// Cache pages
	store := persistence.NewInMemoryStore(time.Hour * 6)

	// SEO repositories
	router.GET("/api/seo/repositories", cache.CachePage(store, time.Hour*24, func(ctx *gin.Context) {
		mongoColl := mongoClient.Database("app").Collection("repos")

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
	// Generator, e. g. "/go/gen"
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

		// Shuffle contributions
		for i := range contribs {
			j := rand.Intn(i + 1)
			contribs[i], contribs[j] = contribs[j], contribs[i]
		}

		ctx.JSON(http.StatusOK, contribs)
	}))
	// Licenses, e. g. "/node/licenses"
	router.GET("/api/:tech/licenses", func(ctx *gin.Context) {
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

		type l struct {
			ID    any `json:"_id" bson:"_id"`
			Repos []struct {
				Author string    `json:"author" bson:"author"`
				Repo   [2]string `json:"repo" bson:"repo"`
				Type   string    `json:"type" bson:"type"`
			} `json:"repos" bson:"repos"`
		}
		var licenses l
		err = mongoColl.FindOne(ctx, bson.D{
			{Key: "_id", Value: model.LICENSES_ID},
		}).Decode(&licenses)
		if err != nil {
			log.Println(err.Error())
			ctx.Status(http.StatusInternalServerError)
			return
		}

		ctx.JSON(http.StatusOK, licenses)
	})
	// Catalogue/namespaces, e. g. "/node"
	router.GET("/api/:tech", cache.CachePage(store, time.Hour*12, func(ctx *gin.Context) {
		// Union type of API and contributions catalogue
		type cat struct {
			NAPIs     int               `json:"n_apis" bson:"n_apis"`
			NContribs int               `json:"n_contribs" bson:"n_contribs"`
			NNs       int               `json:"n_ns" bson:"n_ns"`
			NRepos    int               `json:"n_repos" bson:"n_repos"`
			Ns        []string          `json:"ns" bson:"ns"`
			Version   string            `json:"version" bson:"version"`
			Vids      map[string]string `json:"vids" bson:"vids"`
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
				{Key: "_id", Value: model.CAT_ID},
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
			{Key: "_id", Value: bson.D{primitive.E{Key: "$ne", Value: model.CAT_ID}}},
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
	// Ads
	router.StaticFile("/e4ed80318d98f275b8df.txt", "./website/e4ed80318d98f275b8df.txt")

	// Address
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

func toPtr[T any](v T) *T { return &v }
