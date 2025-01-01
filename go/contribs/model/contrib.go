package model

type (
	Contrib struct {
		Locus     []Locus `json:"locus" bson:"locus"`
		Code      string  `json:"code" bson:"code"`
		Filename  string  `json:"filename" bson:"filename"`
		Filepath  string  `json:"filepath" bson:"filepath"`
		RepoName  string  `json:"repo_name" bson:"repo_name"`
		RepoOwner string  `json:"repo_owner" bson:"repo_owner"`
	}

	Locus struct {
		Ident string `json:"ident" bson:"ident"` // bytes.Buffer, time.Now
		Line  int    `json:"line" bson:"line"`   // 4
	}
)
