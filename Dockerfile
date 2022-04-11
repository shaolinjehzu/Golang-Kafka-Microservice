FROM golang:latest

WORKDIR /app

COPY . /app

RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon
RUN go mod download

ENTRYPOINT CompileDaemon --build="go build -o main cmd/main.go" --command=./main