FROM golang:1.26

WORKDIR /app

COPY go.* .
RUN go mod download

COPY ./cmd ./cmd/
COPY ./internal ./internal/
COPY ./utils ./utils/
COPY ./tests ./tests/
COPY ./entrypoint_tests.sh ./

# ENABLE CGO
ENTRYPOINT ["./entrypoint_tests.sh"]