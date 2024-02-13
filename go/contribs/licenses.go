package main

import (
	"context"
	"contribs-go/model"

	"go.mongodb.org/mongo-driver/bson"
)

func saveLicenses(ctx context.Context) error {
	_, err := mongoColl.DeleteOne(ctx, bson.M{
		"_id": model.LICENSES_ID,
	})
	if err != nil {
		return err
	}

	var licenses = model.Licenses{
		ID: model.LICENSES_ID,
		Repos: []struct {
			Author string    "json:\"author\" bson:\"author\""
			Repo   [2]string "json:\"repo\" bson:\"repo\""
			Type   string    "json:\"type\" bson:\"type\""
		}{
			{
				Author: "GitHub Inc.",
				Repo:   [2]string{"cli", "cli"},
				Type:   "MIT",
			},
			{
				Author: "Traefik Labs",
				Repo:   [2]string{"traefik", "traefik"},
				Type:   "MIT",
			},
			{
				Author: "Docker, Inc.",
				Repo:   [2]string{"moby", "moby"},
				Type:   "Apache License 2.0",
			},
			{
				Author: "Docker, Inc.",
				Repo:   [2]string{"docker", "compose"},
				Type:   "Apache License 2.0",
			},
			{
				Author: "Podman",
				Repo:   [2]string{"containers", "podman"},
				Type:   "Apache License 2.0",
			},
			{
				Author: "The Kubernetes Authors",
				Repo:   [2]string{"helm", "helm"},
				Type:   "Apache License 2.0",
			},
			{
				Author: "The Kubernetes Authors",
				Repo:   [2]string{"kubernetes", "kubernetes"},
				Type:   "Apache License 2.0",
			},
			{
				Author: "MinIO, Inc.",
				Repo:   [2]string{"minio", "minio"},
				Type:   "GNU Affero General Public License v3.0",
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
				Type:   "MIT",
			},
			{
				Author: "nektos",
				Repo:   [2]string{"nektos", "act"},
				Type:   "MIT",
			},
			{
				Author: "Slack Technologies, Inc.",
				Repo:   [2]string{"slackhq", "nebula"},
				Type:   "MIT",
			},
			{
				Author: "The Gitea Authors, The Gogs Authors",
				Repo:   [2]string{"go-gitea", "gitea"},
				Type:   "MIT",
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
				Type:   "GNU Affero General Public License v3.0",
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
		},
	}
	_, err = mongoColl.InsertOne(context.TODO(), licenses)
	return err
}
