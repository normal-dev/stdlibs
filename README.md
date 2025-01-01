# stdlibs.com

Hand-picked API examples for your favorite technology including Go, Node.js, Python and more.

>stdlibs.com employs static code analysis to identify and extract real-world examples of standard library usage from curated open-source repositories.

## Developing

### Prerequisities

The following dependencies are required:

- Docker
- MongoDB
- Go

Depending on what you would like to work on, these dependencies are optional:

- Node.js (and npm)
- Python (and pip)

If you would like to work on any contribution, a GitHub personal access token is
required, which needs to passed as `GITHUB_ACCESS_TOKEN_CONTRIBS` flag.

The easiest way to get started is to use the launch configurations for Visual
Studio Code, which can be found inside `.vscode/launch.json`. If you don't use
Visual Studio Code, consult the following sections.

### Web

Web contains the frontend (or client) and the API server. Make sure the MongoDB
deamon is running and serve the client first with:

```shell
cd app/web
npm install
npm run serve
```

Now start the server with:

```shell
cd app
go run . --no-client
```

You should now be able open the browser and see some interface at
`http://localhost:5173`. **Attention**: It's expected to receive the following
error:

> invalid argument to Int63n

That's because they are no contributions in your database. To add them, follow
the "contributions" section.

### Contributions

Every technology (Go, Node.js, etc.) has a Dockerfile named
`Dockerfile.contribs`. These are used to build production container. While it's
possible to use them for development as well, the recommended way is to use
local compiler. Each contribution is written into your MongoDB database.

#### Go

```shell
cd go/contribs
GITHUB_ACCESS_TOKEN_CONTRIBS=YOUR_PERSONAL_ACCESS_TOKEN go run .
```

#### Node.js

```shell
cd node/mongo
npm install
cd ../contribs
GITHUB_ACCESS_TOKEN_CONTRIBS=YOUR_PERSONAL_ACCESS_TOKEN node index.mjs
```

#### Python

```shell
cd python/contribs
pip install -r requirements.txt --break-system-packages
GITHUB_ACCESS_TOKEN_CONTRIBS=YOUR_PERSONAL_ACCESS_TOKEN python3 __main__.py
```

### APIs

Similar to contributions, each API has it's own Dockerfile named
`Dockerfile.apis`. Every API is written into your MongoDB database.

#### Go

```shell
cd go/apis
go run .
```

#### Node.js

```shell
cd node/apis
npm install
npm start
```

#### Python

You might need to install `Tkinter` first (assuming Linux):

```shell
sudo apt install python3-tk -y
```

Now you should be able to run:

```shell
cd python/apis
pip install -r requirements.txt --break-system-packages
python3 __main__.py
```

## Production

stdlibs.com is running on production using Google Cloud Run instances. The web
app is a service and each contribution and API is a manually invoked job.
`seo/repos` is build via buildpacks.
