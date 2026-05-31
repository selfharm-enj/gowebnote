FROM golang:1.26

WORKDIR /app

COPY go.* .
COPY ./cmd ./cmd/
COPY ./internal ./internal/
COPY ./utils ./utils/
COPY ./entrypoint_prod.sh ./

RUN go mod download
RUN env GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o main ./cmd/main.go

LABEL appname="noteapp"

EXPOSE 8080

ENTRYPOINT ["./entrypoint_prod.sh"]