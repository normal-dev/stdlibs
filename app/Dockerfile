FROM golang:1.23.0

# Client
WORKDIR /
ENV NODE_VERSION=23.4.0
RUN apt install -y curl
RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.1/install.sh | bash
ENV NVM_DIR=/root/.nvm
RUN . "$NVM_DIR/nvm.sh" && nvm install ${NODE_VERSION}
RUN . "$NVM_DIR/nvm.sh" && nvm use v${NODE_VERSION}
RUN . "$NVM_DIR/nvm.sh" && nvm alias default v${NODE_VERSION}
ENV PATH="/root/.nvm/versions/node/v${NODE_VERSION}/bin/:${PATH}"
COPY web /web
COPY website /website
WORKDIR /web
RUN npm ci
RUN npm run build

# Server
WORKDIR /
COPY go.mod go.sum /
RUN go mod download
COPY model /model
COPY *.go /
ENV GIN_MODE=release
RUN CGO_ENABLED=0 GOOS=linux go build -o ./app

EXPOSE 8080

CMD ["./app"]
