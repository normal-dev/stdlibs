package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"mongo"
	"os"
	"os/exec"
	filepat "path/filepath"
	"sync"

	goapis "apis-go/api"

	"contribs-go/model"

	"github.com/google/go-github/github"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	gopkgs = make(map[string]struct{})

	mu sync.Mutex
)

func init() {
	for _, api := range goapis.Get() {
		gopkgs[api.Ns] = struct{}{}
	}
}

const workersn = 3

func main() {
	ctx := context.Background()

	repos, err := findRepos(ctx, ghclient)
	checkErr(err)
	reposn := len(repos)
	log.Printf("repos: %d", reposn)

	var (
		reposchan = make(chan *github.Repository)
		wg        sync.WaitGroup

		contribsn, filesn int
	)
	defer func() {
		checkErr(saveCatalogue(ctx, contribsn, reposn))
		checkErr(insertLicenses(ctx))
	}()
	for range workersn {
		go worker(
			ctx,
			reposchan,
			&wg,
			&contribsn,
			&filesn,
		)
	}
	for _, repo := range repos {
		reposchan <- repo
	}

	wg.Wait()

	log.Printf("contribs: %d", contribsn)
	log.Printf("files: %d", filesn)
}

func worker(
	ctx context.Context,
	repos <-chan *github.Repository,
	wg *sync.WaitGroup,
	contribsn,
	filesn *int,
) {
	f := func(repo *github.Repository) {
		wg.Add(1)
		defer wg.Done()

		repoOwner := repo.Owner.GetLogin()
		repoName := repo.GetName()

		logger := log.New(
			os.Stdout,
			fmt.Sprintf("%s/%s: ", repoOwner, repoName),
			log.Lmsgprefix,
		)

		repoDir, err := os.MkdirTemp("", fmt.Sprintf("%s_%s", repoOwner, repoName))
		if err != nil {
			logErr(logger, err)
			return
		}

		logger.Printf("cloning: %s", repo.GetCloneURL())
		if err := exec.Command("git",
			"clone",
			"-q",
			"--depth", "1",
			"--no-tags",
			"--filter=blob:limit=75k",
			*repo.CloneURL,
			repoDir,
		).Run(); err != nil {
			logErr(logger, err)
			logErr(logger, os.RemoveAll(repoDir))
			return
		}

		rmExtraneous(logger, repoDir)

		var repofilesn int
		go func() {
			err := filepat.Walk(repoDir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				repofilesn++
				return nil
			})
			mu.Lock()
			*filesn += repofilesn
			mu.Unlock()

			logErr(logger, err)
		}()

		var (
			gofilesn, locusn int

			contribs = make([]any, 0)
		)
		for file := range findGoFiles(repoDir) {
			logger.Printf("file: %s", file)
			gofilesn++

			fileBytes, err := os.ReadFile(file)
			if err != nil {
				logErr(logger, err)
				continue
			}
			locus, ok, err := findLocus(fileBytes)
			if !ok {
				logErr(logger, err)
				continue
			}
			locusn += len(locus)
			logger.Printf("locus: %d", len(locus))

			pat := file[len(repoDir):]
			code := string(fileBytes)
			filepath := filepat.Dir(pat)
			filename := filepat.Base(pat)
			contribs = append(contribs, model.Contrib{
				Locus:     locus,
				Code:      code,
				Filepath:  filepath,
				Filename:  filename,
				RepoOwner: repoOwner,
				RepoName:  repoName,
			})

			mu.Lock()
			*contribsn += 1
			mu.Unlock()
		}

		// Remove temporary repository directory
		checkErr(os.RemoveAll(repoDir))

		logger.Printf("contribs: %d", len(contribs))
		logger.Printf("locus: %d", locusn)
		logger.Printf("files: %d", gofilesn)

		if len(contribs) == 0 {
			return
		}

		// Delete existing contributions
		_, err = mongoColl.DeleteMany(ctx, bson.M{
			"repo_owner": repoOwner,
			"repo_name":  repoName,
		})
		checkErr(err)
		// Save new contributions
		_, err = mongoColl.InsertMany(ctx, contribs)
		checkErr(err)
	}
	for repo := range repos {
		f(repo)
	}
}

func findLocus(src []byte) ([]model.Locus, bool, error) {
	ex := newExtractor(src)
	if ex.Error != nil {
		return []model.Locus{}, false, ex.Error
	}

	locus := ex.Extract()
	if ex.Error != nil {
		return []model.Locus{}, false, ex.Error
	}

	ret := make([]model.Locus, 0)
	for api := range locus {
		ret = append(ret, api)
	}
	return ret, len(ret) > 0, nil
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

// Removes directories like "vendor"
func rmExtraneous(logger *log.Logger, dir string) {
	logErr(logger, os.RemoveAll(fmt.Sprintf("%s/.git", dir)))

	_ = filepat.WalkDir(dir, func(path string, dirEntry fs.DirEntry, err error) error {
		if dirEntry.IsDir() {
			if dir := filepat.Base(path); dir == "vendor" {
				logErr(logger, os.RemoveAll(path))
			}
		}
		return nil
	})
}

func saveCatalogue(ctx context.Context, contribsn, reposn int) error {
	coll := mongo.Client.Database(mongo.DB_CONTRIBS).Collection("go")
	if _, err := coll.DeleteOne(ctx, bson.M{"_id": catalogue_id}); err != nil {
		return err
	}
	_, err := coll.InsertOne(ctx, model.Cat{
		ID:        catalogue_id,
		NContribs: contribsn,
		NRepos:    reposn,
	})
	return err
}

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func logErr(logger *log.Logger, err error) {
	logger.SetOutput(os.Stderr)
	defer logger.SetOutput(os.Stdout)
	if err != nil {
		logger.Println(err.Error())
	}
}
