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

var gopkgs = make(map[string]struct{})

func init() {
	for _, api := range goapis.Get() {
		gopkgs[api.Ns] = struct{}{}
	}
}

var mongoColl = goapis.MongoClient.Database("contribs").Collection("go")

func main() {
	ctx := context.TODO()
	githubAccessTok := os.Getenv("GITHUB_ACCESS_TOKEN_CONTRIBS")
	if githubAccessTok == "" {
		panic("can't find Github access token")
	}
	ghclient := newGithubClient(githubAccessTok)
	repos, err := getHandpickedRepos(ctx, ghclient)
	checkErr(err)
	log.Printf("found %d handpicked repos", len(repos))

	var contribsn, filesn int
	reposchan := make(chan *github.Repository)
	for idx := range 2 { // n workers
		go worker(idx+1, reposchan, &contribsn, &filesn)
	}
	for _, repo := range repos {
		reposchan <- repo
	}

	log.Printf("found %d Go contributions", contribsn)
	log.Printf("found approx. %d files", filesn)

	log.Printf("saving catalogue...")
	checkErr(saveCatalogue(ctx, contribsn, len(repos)))

	log.Printf("saving licenses...")
	checkErr(saveLicenses(ctx))
}

func worker(workerID int, repos <-chan *github.Repository, contribsn, filesn *int) {
	for repo := range repos {
		repoOwner := repo.Owner.GetLogin()
		repoName := repo.GetName()

		logger := log.New(
			os.Stdout,
			fmt.Sprintf("%s/%s[%d]", repoOwner, repoName, workerID),
			log.Lmsgprefix,
		)

		logger.Printf("creating temp repo dir...")
		repoDir, err := os.MkdirTemp("", fmt.Sprintf("%s/%s", repoOwner, repoName))
		checkErr(err)

		defer checkErr(os.RemoveAll(repoDir))

		ctx := context.TODO()

		logger.Println("cleaning...")
		checkErr(clean(ctx, repoDir, repoOwner, repoName))

		logger.Printf("cloning repo %s to %s...", repo.GetCloneURL(), repoDir)

		if err := exec.Command("git",
			"clone",
			"-q",
			"--depth", "1",
			"--no-tags",
			"--filter=blob:limit=40k",
			*repo.CloneURL,
			repoDir,
		).Run(); err != nil {
			logErr(logger, err)
			continue
		}

		logger.Println("cleaning repo files...")
		cleanRepo(logger, repoDir)

		var repofilesn int
		go func() {
			err := filepat.Walk(repoDir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				repofilesn++
				return nil
			})
			*filesn += repofilesn
			logErr(logger, err)
		}()

		var (
			gofilesn, apisn int

			contribs = make([]any, 0)
		)
		log.Println("looking for Go files...")
		for file := range findGoFiles(repoDir) {
			gofilesn++

			fileBytes, err := os.ReadFile(file)
			if err != nil {
				logErr(logger, err)
				continue
			}
			apis, ok, _ := findGoAPIs(fileBytes)
			if !ok {
				continue
			}
			apisn += len(apis)

			pat := file[len(repoDir):]
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
			*contribsn += 1
		}

		logger.Printf("found %d contributions (%d Go apis)", len(contribs), apisn)
		logger.Printf("found approx. %d cloned files (%d Go files)", repofilesn, gofilesn)
		_, err = saveContribs(ctx, contribs)
		logErr(logger, err)
	}
}

func getHandpickedRepos(ctx context.Context, ghClient *github.Client) (repos []*github.Repository, err error) {
	for _, repo := range [][2]string{
		{"cli", "cli"},
		{"traefik", "traefik"},
		{"moby", "moby"},
		{"docker", "compose"},
		{"containers", "podman"},
		{"helm", "helm"},
		{"kubernetes", "kubernetes"},
		{"minio", "minio"},
		{"cloudflare", "cloudflared"},
		{"cosmos", "cosmos-sdk"},
		{"aws", "karpenter"},
		{"cilium", "cilium"},
		{"containerd", "containerd"},
		{"containers", "buildah"},
		{"hyperledger", "fabric"},
		{"istio", "istio"},
		{"pingcap", "tidb"},
		{"vitessio", "vitess"},
		{"go-delve", "delve"},
		{"nektos", "act"},
		{"slackhq", "nebula"},
		{"go-gitea", "gitea"},
		{"vmware-tanzu", "velero"},
		{"vmware-tanzu", "sonobuoy"},
		{"gravitational", "teleport"},
		{"canonical", "lxd"},
		{"eolinker", "apinto"},
		{"portainer", "portainer"},
		{"hyperledger", "firefly"},
		{"gin-gonic", "gin"},
		{"mattermost", "mattermost"},
		{"beego", "beego"},
		{"securego", "gosec"},
		{"goreleaser", "goreleaser"},
		{"caddyserver", "caddy"},
		{"gopherjs", "gopherjs"},
		{"v2ray", "v2ray-core"},
		{"ollama", "ollama"},
		{"spf13", "cobra"},
		{"tailscale", "tailscale"},
		{"rancher", "rancher"},
		{"google", "syzkaller"},
		{"goplus", "gop"},
		{"ignite", "cli"},
		{"apache", "incubator-devlake"},
		{"rclone", "rclone"},
		{"prometheus", "prometheus"},
		{"benthosdev", "benthos"},
		{"temporalio", "temporal"},
		{"thanos-io", "thanos"},
		{"envoyproxy", "envoy"},
		{"ebitengine", "purego"},
		{"goplus", "igop"},
		{"alecthomas", "kong"},
		{"alecthomas", "participle"},
		{"go-critic", "go-critic"},
		{"gohugoio", "hugo"},
		{"harness", "gitness"},
		{"aquasecurity", "trivy"},
	} {
		owner, name := repo[0], repo[1]
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

	ret := make([]model.API, 0)
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
func cleanRepo(logger *log.Logger, dir string) {
	logErr(logger, os.RemoveAll(fmt.Sprintf("%s/.git", dir)))

	_ = filepat.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			if dir := filepat.Base(path); dir == "vendor" {
				logErr(logger, os.RemoveAll(path))
			}
		}
		return nil
	})
}

func newGithubClient(githubAccessTok string) *github.Client {
	tokSrc := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubAccessTok},
	)
	httpClient := oauth2.NewClient(context.TODO(), tokSrc)
	return github.NewClient(httpClient)
}

func clean(ctx context.Context, repoDir, repoOwner, repoName string) error {
	if err := os.RemoveAll(repoDir); err != nil {
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

func logErr(logger *log.Logger, err error) {
	logger.SetOutput(os.Stderr)
	defer logger.SetOutput(os.Stdout)
	if err != nil {
		logger.Println(err.Error())
	}
}
