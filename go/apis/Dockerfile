FROM golang:1.22

WORKDIR /

COPY api /api
COPY model /model
COPY go.mod go.sum ./
COPY *.go ./

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o apis_go
RUN export PATH=$PATH:/usr/local/go/bin

CMD ["./apis_go"]
