package main

import (
	"context"
	"fmt"
	"log"
	"slices"

	"github.com/google/go-github/github"
	"go.mongodb.org/mongo-driver/bson"
)

type license struct {
	Author string    "json:\"author\" bson:\"author\""
	Repo   [2]string "json:\"repo\" bson:\"repo\""
	Type   string    "json:\"type\" bson:\"type\""
}

var (
	licenses = []license{
		{
			Author: "GitHub Inc.",
			Repo:   [2]string{"cli", "cli"},
			Type:   "MIT license",
		},
		{
			Author: "Traefik Labs",
			Repo:   [2]string{"traefik", "traefik"},
			Type:   "MIT license",
		},
		{
			Author: "Docker, Inc.",
			Repo:   [2]string{"moby", "moby"},
			Type:   "Apache license 2.0",
		},
		{
			Author: "Docker, Inc.",
			Repo:   [2]string{"docker", "compose"},
			Type:   "Apache license 2.0",
		},
		{
			Author: "Podman",
			Repo:   [2]string{"containers", "podman"},
			Type:   "Apache license 2.0",
		},
		{
			Author: "The Kubernetes Authors",
			Repo:   [2]string{"helm", "helm"},
			Type:   "Apache license 2.0",
		},
		{
			Author: "The Kubernetes Authors",
			Repo:   [2]string{"kubernetes", "kubernetes"},
			Type:   "Apache license 2.0",
		},
		{
			Author: "MinIO, Inc.",
			Repo:   [2]string{"minio", "minio"},
			Type:   "GNU Affero general public license v3.0",
		},
		{
			Author: "Cloudflare, Inc.",
			Repo:   [2]string{"cloudflare", "cloudflared"},
			Type:   "Apache-2.0 license",
		},
		{
			Author: "Interchain Foundation",
			Repo:   [2]string{"cosmos", "cosmos-sdk"},
			Type:   "Apache-2.0 license",
		},
		{
			Author: "Amazon.com, Inc. or its affiliates",
			Repo:   [2]string{"aws", "karpenter"},
			Type:   "Apache-2.0 license",
		},
		{
			Author: "The Cilium Authors",
			Repo:   [2]string{"cilium", "cilium"},
			Type:   "Apache-2.0 license",
		},
		{
			Author: "The containerd Authors",
			Repo:   [2]string{"containerd", "containerd"},
			Type:   "Apache-2.0 license",
		},
		{
			Author: "containers",
			Repo:   [2]string{"containers", "buildah"},
			Type:   "Apache-2.0 license",
		},
		{
			Author: "Hyperledger Foundation",
			Repo:   [2]string{"hyperledger", "fabric"},
			Type:   "Apache-2.0 license",
		},
		{
			Author: "the Istio Authors",
			Repo:   [2]string{"istio", "istio"},
			Type:   "Apache-2.0 license",
		},
		{
			Author: "PingCAP",
			Repo:   [2]string{"pingcap", "tidb"},
			Type:   "Apache-2.0 license",
		},
		{
			Author: "The Linux Foundation",
			Repo:   [2]string{"vitessio", "vitess"},
			Type:   "Apache-2.0 license",
		},
		{
			Author: "Derek Parker",
			Repo:   [2]string{"go-delve", "delve"},
			Type:   "MIT license",
		},
		{
			Author: "nektos",
			Repo:   [2]string{"nektos", "act"},
			Type:   "MIT license",
		},
		{
			Author: "Slack Technologies, Inc.",
			Repo:   [2]string{"slackhq", "nebula"},
			Type:   "MIT license",
		},
		{
			Author: "The Gitea Authors, The Gogs Authors",
			Repo:   [2]string{"go-gitea", "gitea"},
			Type:   "MIT license",
		},
		{
			Author: "Broadcom",
			Repo:   [2]string{"vmware-tanzu", "velero"},
			Type:   "Apache-2.0 license",
		},
		{
			Author: "Broadcom",
			Repo:   [2]string{"vmware-tanzu", "sonobuoy"},
			Type:   "Apache-2.0 license",
		},
		{
			Author: "Gravitational, Inc.",
			Repo:   [2]string{"gravitational", "teleport"},
			Type:   "GNU Affero general public license v3.0",
		},
		{
			Author: "Canonical Ltd.",
			Repo:   [2]string{"canonical", "lxd"},
			Type:   "Apache-2.0 license",
		},
		{
			Author: "Shenzhen Silver Cloud Information Technology Co., Ltd.",
			Repo:   [2]string{"eolinker", "apinto"},
			Type:   "Apache-2.0 license",
		},
		{
			Author: "Portainer.io",
			Repo:   [2]string{"portainer", "portainer"},
			Type:   "Zlib license",
		},
		{
			Author: "Hyperledger Foundation",
			Repo:   [2]string{"hyperledger", "firefly"},
			Type:   "Apache-2.0 license",
		},
		{
			Author: "Manuel Martínez-Almeida",
			Repo:   [2]string{"gin-gonic", "gin"},
			Type:   "MIT license",
		},
		{
			Author: "Mattermost, Inc.",
			Repo:   [2]string{"mattermost", "mattermost"},
			Type:   "Apache-2.0 license",
		},
		{
			Author: "Beego",
			Repo:   [2]string{"beego", "beego"},
			Type:   "Apache-2.0 license",
		},
		{
			Author: "Carlos Alexandro Becker",
			Repo:   [2]string{"goreleaser", "goreleaser"},
			Type:   "MIT license",
		},
		{
			Author: "Grant Murphy",
			Repo:   [2]string{"securego", "gosec"},
			Type:   "Apache-2.0 license",
		},
		{
			Author: "ZeroSSL",
			Repo:   [2]string{"caddyserver", "caddy"},
			Type:   "Apache-2.0 license",
		},
		{
			Author: "Richard Musiol",
			Repo:   [2]string{"gopherjs", "gopherjs"},
			Type:   "BSD-2-Clause license",
		},
		{
			Author: "V2Fly Community",
			Repo:   [2]string{"v2ray", "v2ray-core"},
			Type:   "MIT license",
		},
		{
			Author: "Ollama",
			Repo:   [2]string{"ollama", "ollama"},
			Type:   "MIT license",
		},
		{
			Author: "spf13",
			Repo:   [2]string{"spf13", "cobra"},
			Type:   "Apache-2.0 license",
		},
		{
			Author: "Tailscale Inc & AUTHORS",
			Repo:   [2]string{"tailscale", "tailscale"},
			Type:   "BSD 3-Clause license",
		},

		{
			Author: "Rancher Labs, Inc.",
			Repo:   [2]string{"rancher", "rancher"},
			Type:   "Apache-2.0 license",
		},

		{
			Author: "syzkaller project authors",
			Repo:   [2]string{"google", "syzkaller"},
			Type:   "Apache-2.0 license",
		},
		{
			Author: "The GoPlus Authors",
			Repo:   [2]string{"goplus", "gop"},
			Type:   "Apache-2.0 license",
		},
		{
			Author: "All in Bits, Inc.",
			Repo:   [2]string{"ignite", "cli"},
			Type:   "Apache-2.0 license",
		},
		{
			Author: "Apache DevLake, DevLake, Apache",
			Repo:   [2]string{"apache", "incubator-devlake"},
			Type:   "Apache-2.0 license",
		},
		{
			Author: "Nick Craig-Wood",
			Repo:   [2]string{"rclone", "rclone"},
			Type:   "MIT license",
		},
		{
			Author: "Prometheus Authors, The Linux Foundation",
			Repo:   [2]string{"prometheus", "prometheus"},
			Type:   "Apache-2.0 license",
		},
		{
			Author: "Ashley Jeffs",
			Repo:   [2]string{"benthosdev", "benthos"},
			Type:   "MIT license",
		},
		{
			Author: "Temporal Technologies Inc., Uber Technologies, Inc.",
			Repo:   [2]string{"temporalio", "temporal"},
			Type:   "MIT license",
		},
		{
			Author: "Fabian Reinartz @fabxc and Bartłomiej Płotka @bwplotka",
			Repo:   [2]string{"thanos-io", "thanos"},
			Type:   "Apache-2.0 license",
		},
		{
			Author: "Envoy Project Authors",
			Repo:   [2]string{"envoyproxy", "envoy"},
			Type:   "Apache-2.0 license",
		},
		{
			Author: "The Ebitengine Authors",
			Repo:   [2]string{"ebitengine", "purego"},
			Type:   "Apache-2.0 license",
		},
		{
			Author: "The GoPlus Authors",
			Repo:   [2]string{"goplus", "igop"},
			Type:   "Apache-2.0 license",
		},
		{
			Author: "Alec Thomas",
			Repo:   [2]string{"alecthomas", "kong"},
			Type:   "MIT license",
		},
		{
			Author: "Alec Thomas",
			Repo:   [2]string{"alecthomas", "participle"},
			Type:   "MIT license",
		},
		{
			Author: "go-critic team",
			Repo:   [2]string{"go-critic", "go-critic"},
			Type:   "MIT license",
		},
		{
			Author: "The Hugo Authors",
			Repo:   [2]string{"gohugoio", "hugo"},
			Type:   "Apache-2.0 license",
		},
		{
			Author: "Harness, Inc.",
			Repo:   [2]string{"harness", "gitness"},
			Type:   "Apache-2.0 license",
		},
		{
			Author: "Aqua Security Software Ltd.",
			Repo:   [2]string{"aquasecurity", "trivy"},
			Type:   "Apache-2.0 license",
		},
		{
			Author: "Nathan Sweet, Cloudflare, Authors of Cilium",
			Repo:   [2]string{"cilium", "ebpf"},
			Type:   "MIT license",
		},
		{
			Author: "Uber Technologies, Inc.",
			Repo:   [2]string{"uber-go", "zap"},
			Type:   "MIT license",
		},
		{
			Author: "StackRox",
			Repo:   [2]string{"stackrox", "stackrox"},
			Type:   "Apache-2.0 license",
		},
		{
			Author: "The frp Authors",
			Repo:   [2]string{"fatedier", "frp"},
			Type:   "Apache-2.0 license",
		},
		{
			Author: "Ava Labs, Inc,",
			Repo:   [2]string{"ava-labs", "avalanchego"},
			Type:   "BSD-3-Clause license",
		},
		{
			Author: "The etcd Authors",
			Repo:   [2]string{"etcd-io", "etcd"},
			Type:   "Apache-2.0 license",
		},
		{
			Author: "The Gonum Authors",
			Repo:   [2]string{"gonum", "plot"},
			Type:   "BSD-3-Clause license",
		},
		{
			Author: "Jakob Borg",
			Repo:   [2]string{"syncthing", "syncthing"},
			Type:   "MPL-2.0 license",
		},
	}

	repos = [][2]string{
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
		{"cilium", "ebpf"},
		{"uber-go", "zap"},
		{"stackrox", "stackrox"},
		{"fatedier", "frp"},
		{"ava-labs", "avalanchego"},
		{"etcd-io", "etcd"},
		{"gonum", "plot"},
		{"syncthing", "syncthing"},
	}
)

func init() {
	for _, repo := range repos {
		ok := slices.ContainsFunc(licenses, func(l license) bool {
			return l.Repo == repo
		})
		if !ok {
			panic(fmt.Sprintf("%s/%s: license or repo missmatch", repo[0], repo[1]))
		}
	}
}

func findRepos(ctx context.Context, ghClient *github.Client) (r []*github.Repository, err error) {
	for _, repo := range repos {
		owner, name := repo[0], repo[1]
		log.Printf("repo: %s/%s...", owner, name)
		repo, _, err := ghClient.Repositories.Get(ctx, owner, name)
		if err != nil {
			return r, err
		}
		r = append(r, repo)
	}
	return
}

func insertLicenses(ctx context.Context) error {
	_, err := mongoColl.DeleteOne(ctx, bson.M{
		"_id": licenses_id,
	})
	if err != nil {
		return err
	}

	doc := bson.D{
		bson.E{Key: "_id", Value: licenses_id},
		bson.E{Key: "repos", Value: licenses},
	}
	_, err = mongoColl.InsertOne(ctx, doc)
	return err
}
