package model

// Catalogue with contributions
type Cat struct {
	ID        any `json:"_id" bson:"_id"`
	NContribs int `json:"n_contribs" bson:"n_contribs"`
	NRepos    int `json:"n_repos" bson:"n_repos"`
}
