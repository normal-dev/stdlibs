# Contributing

## APIs

See [glossary](GLOSSARY.md#api) for an explanation.

### General process

Ideally, every programming language (or technology) handler is written in its
respective programming language. Each programming language or technology shall
share the same algorithm:

1. Get all importable symbols from the (standard) library
2. Remove internal/non-exported symbols
3. Iterate over each symbol and collect:
   1. Symbol name, e. g. `Buffer`
   2. Namespace, e. g. `bytes`
   3. Type, e. g. `int`
   4. Value, e. g. `512`
4. Save all APIs to database
5. Save a catalogue which contains the amount of APIs, the amount of namespaces,
   the namespaces and programming language or technology version

### Schemas

#### API

```json
{
    "_id": "",
    "doc": "",
    "name": "",
    "type": "",
    "ns": "",
    "value": ""
}
```

Example:

```json
{
    "_id": "go/token.AND_ASSIGN",
    "doc": "",
    "name": "AND_ASSIGN",
    "type": "int",
    "ns": "go/token",
    "value": "28"
}
```

#### Catalogue

```json
{
    "_id": "_cat",
    "n_apis": 0,
    "n_ns": 0,
    "ns": [],
    "version": ""
}
```

Example:

```json
{
    "_id": "_cat",
    "n_apis": 8025,
    "n_ns": 161,
    "ns": [
        "go/constant",
        "fmt",
        "regexp",
        "crypto/rand",
        "encoding"
    ],
    "version": "1.22.0"
}
```

## Contributions

See [glossary](GLOSSARY.md#contribution) for an explanation.

### General process

To clone multiple repositories, a GitHub personal access token needs to be 
passed as `GITHUB_ACCESS_TOKEN_CONTRIBS` environmental variable.

Ideally, every programming language (or technology) handler is written in its
respective programming language. Each programming language or technology shall
share the same algorithm:

1. Fetch repositories which contain the programming language or technology and
   (standard) library usages
2. Iterate over each repository
   1. Clone repository to a temporary directory, e. g.
      `/tmp/traefik_traefik`
      1. (Optional): Remove extraneous directories like `.git`
   2. Find files which qualify for a contribution, e. g. `.go` files
   3. Iterate over each file and collect:
      1. Locus (see [glossary](GLOSSARY.md#locus)) for each
      file
      2. Amount of files
      3. Amount of contributions
      4. Source code, e. g. `var x int\nx = 5`
      5. File path, e. g. `/src/models`
      6. File name, `user.go`
      7. Repository owner `traefik`
      8. Repository name `traefik`
   4. Delete every contribution of the repository in the database
   5. Save all contributions to database
3. Save a catalogue which contains the amount of contributions and
   repositories to the database
4. Save license information for each repository to database

### Schemas

#### Locus

```json
{
    "ident": "",
    "line": 0
}
```

#### Contribution

```json
{
    "locus": [],
    "code": "",
    "filename": "",
    "filepath": "",
    "repo_name": "",
    "repo_owner": ""
}
```

Example:

```json
{
    "locus": [
        {
            "ident": "io.ReadAll",
            "line": 10
        },
        {
            "ident": "io.ReadCloser",
            "line": 8
        },
        {
            "ident": "os.ReadFile",
            "line": 15
        }
    ],
    "code": "package cmdutil\n\nimport (\n\t\"io\"\n\t\"os\"\n)\n\nfunc ReadFile(filename string, stdin io.ReadCloser) ([]byte, error) {\n\tif filename == \"-\" {\n\t\tb, err := io.ReadAll(stdin)\n\t\t_ = stdin.Close()\n\t\treturn b, err\n\t}\n\n\treturn os.ReadFile(filename)\n}\n",
    "filename": "file_input.go",
    "filepath": "/pkg/cmdutil",
    "repo_name": "cli",
    "repo_owner": "cli"
}
```

#### License

```json
{
    "_id": "_licenses",
    "repos": []
}
```

Example:

```json
{
    "_id": "_licenses",
    "repos": [
        {
            "author": "GitHub Inc.",
            "repo": [
                "cli",
                "cli"
            ],
            "type": "MIT license"
        },
        {
            "author": "Traefik Labs",
            "repo": [
                "traefik",
                "traefik"
            ],
            "type": "MIT license"
        }
    ]
}
```

#### Catalogue

```json
{
    "_id": "_cat",
    "n_contribs": 0,
    "n_repos": 0
}
```

Example:

```json
{
    "_id": "_cat",
    "n_contribs": 21459,
    "n_repos": 89
}
```

### Database

MongoDB is used as database. On production, a contribution will be saved into a
MongoDB database at `MONGO_DB_URI` (during development it falls back to:
`mongodb://localhost:27017`).
