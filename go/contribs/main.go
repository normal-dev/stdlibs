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
	"golang.org/x/oauth2"
)

var gopkgs = make(map[string]struct{})

func init() {
	for _, api := range goapis.Get() {
		gopkgs[api.Ns] = struct{}{}
	}
}

const workersn = 3

var mongoColl = mongo.MongoClient.Database("contribs").Collection("go")

func main() {
	ctx := context.TODO()

	githubAccessTok := os.Getenv("GITHUB_ACCESS_TOKEN_CONTRIBS")
	if githubAccessTok == "" {
		panic("can't find Github access token")
	}
	ghclient := newGithubClient(githubAccessTok)
	repos, err := getRepos(ctx, ghclient)
	checkErr(err)
	log.Printf("repos: %d", len(repos))

	var wg sync.WaitGroup
	reposchan := make(chan *github.Repository)
	var contribsn, filesn int
	for range workersn { // n workers
		go worker(
			ctx,
			&wg,
			reposchan,
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

	checkErr(saveCatalogue(ctx, contribsn, len(repos)))
	checkErr(saveLicenses(ctx))
}

func worker(
	ctx context.Context,
	wg *sync.WaitGroup,
	repos <-chan *github.Repository,
	contribsn,
	filesn *int,
) {
	var mu sync.Mutex
	for repo := range repos {
		wg.Add(1)

		repoOwner := repo.Owner.GetLogin()
		repoName := repo.GetName()

		logger := log.New(
			os.Stdout,
			fmt.Sprintf("%s/%s: ", repoOwner, repoName),
			log.Lmsgprefix,
		)

		repoDir, err := os.MkdirTemp("", fmt.Sprintf("%s_%s", repoOwner, repoName))
		checkErr(err)
		checkErr(rmRepo(ctx, repoDir, repoOwner, repoName))

		logger.Printf("cloning: %s", repo.GetCloneURL())

		if err := exec.Command("git",
			"clone",
			"-q",
			"--depth", "1",
			"--no-tags",
			"--filter=blob:limit=50k",
			*repo.CloneURL,
			repoDir,
		).Run(); err != nil {
			logErr(logger, err)
			logErr(logger, os.RemoveAll(repoDir))
			continue
		}

		stripeRepo(logger, repoDir)

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
			defer mu.Unlock()
			*filesn += repofilesn

			logErr(logger, err)
		}()

		var (
			gofilesn, locusn int

			contribs = make([]any, 0)
		)
		for file := range findGoFiles(repoDir) {
			log.Printf("file: %s", file)
			gofilesn++

			fileBytes, err := os.ReadFile(file)
			if err != nil {
				logErr(logger, err)
				continue
			}
			locus, ok, _ := findLocus(fileBytes)
			if !ok {
				continue
			}
			locusn += len(locus)
			log.Printf("locus: %d", len(locus))

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
			defer mu.Unlock()
			*contribsn += 1
		}

		logger.Printf("contribs: %d", len(contribs))
		logger.Printf("locus: %d", locusn)
		logger.Printf("files: %d", gofilesn)
		_, err = saveContribs(ctx, contribs)
		logErr(logger, err)

		go logErr(logger, os.RemoveAll(repoDir))

		wg.Done()
	}
}

func getRepos(ctx context.Context, ghClient *github.Client) (repos []*github.Repository, err error) {
	for _, repo := range [][2]string{
		{"cli", "cli"},
		{"traefik", "traefik"},
		// {"moby", "moby"},
		// {"docker", "compose"},
		// {"containers", "podman"},
		// {"helm", "helm"},
		// {"kubernetes", "kubernetes"},
		// {"minio", "minio"},
		// {"cloudflare", "cloudflared"},
		// {"cosmos", "cosmos-sdk"},
		// {"aws", "karpenter"},
		// {"cilium", "cilium"},
		// {"containerd", "containerd"},
		// {"containers", "buildah"},
		// {"hyperledger", "fabric"},
		// {"istio", "istio"},
		// {"pingcap", "tidb"},
		// {"vitessio", "vitess"},
		// {"go-delve", "delve"},
		// {"nektos", "act"},
		// {"slackhq", "nebula"},
		// {"go-gitea", "gitea"},
		// {"vmware-tanzu", "velero"},
		// {"vmware-tanzu", "sonobuoy"},
		// {"gravitational", "teleport"},
		// {"canonical", "lxd"},
		// {"eolinker", "apinto"},
		// {"portainer", "portainer"},
		// {"hyperledger", "firefly"},
		// {"gin-gonic", "gin"},
		// {"mattermost", "mattermost"},
		// {"beego", "beego"},
		// {"securego", "gosec"},
		// {"goreleaser", "goreleaser"},
		// {"caddyserver", "caddy"},
		// {"gopherjs", "gopherjs"},
		// {"v2ray", "v2ray-core"},
		// {"ollama", "ollama"},
		// {"spf13", "cobra"},
		// {"tailscale", "tailscale"},
		// {"rancher", "rancher"},
		// {"google", "syzkaller"},
		// {"goplus", "gop"},
		// {"ignite", "cli"},
		// {"apache", "incubator-devlake"},
		// {"rclone", "rclone"},
		// {"prometheus", "prometheus"},
		// {"benthosdev", "benthos"},
		// {"temporalio", "temporal"},
		// {"thanos-io", "thanos"},
		// {"envoyproxy", "envoy"},
		// {"ebitengine", "purego"},
		// {"goplus", "igop"},
		// {"alecthomas", "kong"},
		// {"alecthomas", "participle"},
		// {"go-critic", "go-critic"},
		// {"gohugoio", "hugo"},
		// {"harness", "gitness"},
		// {"aquasecurity", "trivy"},
		// {"cilium", "ebpf"},
		// {"uber-go", "zap"},
		// {"stackrox", "stackrox"},
		// {"fatedier", "frp"},
		// {"ava-labs", "avalanchego"},
		// {"etcd-io", "etcd"},
		// {"gonum", "plot"},
	} {
		owner, name := repo[0], repo[1]
		log.Printf("repo: %s/%s...", owner, name)
		repo, _, err := ghClient.Repositories.Get(ctx, owner, name)
		if err != nil {
			return repos, err
		}
		repos = append(repos, repo)
	}
	return
}

// Removes directories like "vendor"
func stripeRepo(logger *log.Logger, dir string) {
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

// Removes the repository directory and database entry
func rmRepo(ctx context.Context, repoDir, repoOwner, repoName string) error {
	if err := os.RemoveAll(repoDir); err != nil {
		return err
	}
	_, err := mongoColl.DeleteMany(ctx, bson.M{
		"repo_owner": repoOwner,
		"repo_name":  repoName,
	})
	return err
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

func saveContribs(ctx context.Context, contribs []any) (int, error) {
	if len(contribs) == 0 {
		return 0, nil
	}
	res, err := mongoColl.InsertMany(ctx, contribs)
	return len(res.InsertedIDs), err
}

func saveCatalogue(ctx context.Context, contribsn, reposn int) error {
	coll := mongo.MongoClient.Database(mongo.DB_CONTRIBS).Collection("go")
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

func newGithubClient(githubAccessTok string) *github.Client {
	tokSrc := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubAccessTok},
	)
	httpClient := oauth2.NewClient(context.TODO(), tokSrc)
	return github.NewClient(httpClient)
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
