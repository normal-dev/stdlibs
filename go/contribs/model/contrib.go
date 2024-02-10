package model

type (
	Contrib struct {
		APIs      []API  `json:"apis" bson:"apis"`
		Code      string `json:"code" bson:"code"`
		Filename  string `json:"filename" bson:"filename"`
		Filepath  string `json:"filepath" bson:"filepath"`
		RepoName  string `json:"repo_name" bson:"repo_name"`
		RepoOwner string `json:"repo_owner" bson:"repo_owner"`
	}

	API struct {
		Ident string `json:"ident" bson:"ident"` // bytes.Buffer, time.Now
		Line  int    `json:"line" bson:"line"`   // 4
	}
)
