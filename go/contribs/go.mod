module contribs-go

go 1.23.0

require (
	apis-go v0.0.0
	github.com/google/go-github v17.0.0+incompatible
	go.mongodb.org/mongo-driver v1.16.1
	golang.org/x/oauth2 v0.23.0
	mongo v0.0.0
)

replace mongo => ../mongo

replace apis-go => ../apis

require (
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	golang.org/x/crypto v0.35.0 // indirect
	golang.org/x/mod v0.21.0 // indirect
	golang.org/x/sync v0.11.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	golang.org/x/tools v0.24.0 // indirect
)
