package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	filepat "path/filepath"

	goapis "contribs-go/api"

	"contribs-go/model"

	"github.com/google/go-github/github"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/oauth2"
)

func init() {
	log.SetFlags(0)
	log.Default().SetOutput(os.Stderr)
}

var gopkgs = make(map[string]struct{})

func init() {
	for _, api := range goapis.Get() {
		gopkgs[api.Ns] = struct{}{}
	}
}

var (
	tmpDir string

	mongoColl = goapis.MongoClient.Database("contribs").Collection("go")
)

func main() {
	log.Printf("creating temp dir...")
	dir, err := os.MkdirTemp("", "contribs-go")
	checkErr(err)
	tmpDir = dir

	defer checkErr(os.RemoveAll(tmpDir))

	githubAccessTok := os.Getenv("GITHUB_ACCESS_TOKEN_CONTRIBS")
	if githubAccessTok == "" {
		panic("can't find Github access token")
	}

	ghclient := newGithubClient(githubAccessTok)
	repos, err := getHandpickedRepos(context.TODO(), ghclient)
	checkErr(err)
	log.Printf("found %d handpicked repos", len(repos))

	var contribsn, filesn int
	for _, repo := range repos {
		repoOwner := repo.Owner.GetLogin()
		repoName := repo.GetName()

		log.Println("cleaning...")
		checkErr(clean(context.TODO(), repoOwner, repoName))

		log.Printf("cloning repo %s to %s...", repo.GetCloneURL(), tmpDir)

		if err := exec.Command("git",
			"clone",
			"-q",
			"--depth", "1",
			"--no-tags",
			"--filter=blob:limit=25k",
			*repo.CloneURL,
			tmpDir,
		).Run(); err != nil {
			log.Println(err.Error())
			continue
		}

		log.Println("cleaning repo files...")
		reductRepo()

		var repofilesn int
		go func() {
			err := filepat.Walk(tmpDir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				repofilesn++
				return nil
			})
			filesn += repofilesn
			logErr(err)
		}()

		var (
			gofilesn, apisn int

			contribs = make([]any, 0)
		)
		log.Println("looking for Go files...")
		for file := range findGoFiles(tmpDir) {
			gofilesn++

			fileBytes, err := os.ReadFile(file)
			logErr(err)
			apis, ok, _ := findGoAPIs(fileBytes)
			if !ok {
				continue
			}
			apisn += len(apis)

			pat := file[len(tmpDir):]
			code := string(fileBytes)
			filepath := filepat.Dir(pat)
			filename := filepat.Base(pat)
			contribs = append(contribs, model.Contrib{
				APIs:      apis,
				Code:      code,
				Filepath:  filepath,
				Filename:  filename,
				RepoOwner: repoOwner,
				RepoName:  repoName,
			})
			contribsn++
		}

		log.Printf("found %d contributions (%d Go apis)", len(contribs), apisn)
		log.Printf("found approx. %d cloned files (%d Go files)", repofilesn, gofilesn)
		docsn, err := saveContribs(context.TODO(), contribs)
		logErr(err)
		log.Printf("%d contributions saved", docsn)
	}

	log.Printf("found %d Go contributions", contribsn)
	log.Printf("found approx. %d files", filesn)

	log.Printf("saving catalogue...")
	checkErr(saveCatalogue(context.TODO(), contribsn, len(repos)))

	log.Printf("saving licenses...")
	checkErr(saveLicenses(context.TODO()))
}

func getHandpickedRepos(ctx context.Context, ghClient *github.Client) (repos []*github.Repository, err error) {
	for _, r := range handpickedRepos {
		owner, name := r[0], r[1]
		log.Printf("fetching repo %s/%s...", owner, name)
		repo, _, err := ghClient.Repositories.Get(ctx, owner, name)
		if err != nil {
			return repos, err
		}
		repos = append(repos, repo)
	}
	return
}

func findGoFiles(dir string) chan string {
	files := make(chan string, 100)
	go func() {
		defer close(files)

		const goext = ".go"
		err := filepat.WalkDir(dir, func(file string, dirEntry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if filepat.Ext(dirEntry.Name()) == goext {
				files <- file
			}
			return nil
		})
		checkErr(err)
	}()
	return files
}

func findGoAPIs(src []byte) ([]model.API, bool, error) {
	ex := NewExtractor(src)
	if ex.Error != nil {
		return []model.API{}, false, ex.Error
	}

	apis := ex.Extract()
	if ex.Error != nil {
		return []model.API{}, false, ex.Error
	}

	var ret = make([]model.API, 0)
	for api := range apis {
		ret = append(ret, api)
	}
	return ret, len(ret) > 0, nil
}

func saveContribs(ctx context.Context, contribs []any) (int, error) {
	if len(contribs) == 0 {
		return 0, nil
	}
	res, err := mongoColl.InsertMany(ctx, contribs)
	return len(res.InsertedIDs), err
}

func saveCatalogue(ctx context.Context, contribsn, reposn int) error {
	coll := goapis.MongoClient.Database(goapis.DB_CONTRIBS).Collection("go")
	if _, err := coll.DeleteOne(ctx, bson.M{"_id": model.CAT_ID}); err != nil {
		return err
	}
	_, err := coll.InsertOne(ctx, model.Cat{
		ID:        model.CAT_ID,
		NContribs: contribsn,
		NRepos:    reposn,
	})
	return err
}

// Removes directories like "vendor"
func reductRepo() {
	logErr(os.RemoveAll(fmt.Sprintf("%s/.git", tmpDir)))

	err := filepat.WalkDir(tmpDir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			if dir := filepat.Base(path); dir == "vendor" {
				logErr(os.RemoveAll(path))
			}
		}
		return nil
	})
	logErr(err)
}

func newGithubClient(githubAccessTok string) *github.Client {
	tokSrc := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubAccessTok},
	)
	httpClient := oauth2.NewClient(context.TODO(), tokSrc)
	return github.NewClient(httpClient)
}

func clean(ctx context.Context, repoOwner, repoName string) error {
	if err := os.RemoveAll(tmpDir); err != nil {
		return err
	}
	_, err := mongoColl.DeleteMany(ctx, bson.M{
		"repo_owner": repoOwner,
		"repo_name":  repoName,
	})
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
