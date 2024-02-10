package model

const (
	CAT_ID = "_cat"
)

// Catalogue with namespaces
type Cat struct {
	ID      any      `json:"_id" bson:"_id"`
	NAPIs   int      `json:"n_apis" bson:"n_apis"`
	NNs     int      `json:"n_ns" bson:"n_ns"`
	Ns      []string `json:"ns" bson:"ns"`
	Version string   `json:"version" bson:"version"`
}
