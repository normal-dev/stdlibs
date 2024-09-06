package main

import (
	"context"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var ghclient *github.Client

func init() {
	accessTok := os.Getenv("GITHUB_ACCESS_TOKEN_CONTRIBS")
	if accessTok == "" {
		panic("can't find Github access token")
	}
	tokSrc := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessTok},
	)
	httpClient := oauth2.NewClient(context.TODO(), tokSrc)
	ghclient = github.NewClient(httpClient)
}
