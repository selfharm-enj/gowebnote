BINARY_NAME=noteapp
.DEFAULT_GOAL := run

build:
	GOARCH=amd64 GOOS=linux go build -o ./target/${BINARY_NAME}-linux ./cmd/main.go

run:
	go run cmd/main.go

test:
	go test ./tests/...

lint:
	golangci-lint run

docker_tests:
# 	sudo rm -rf ./db-data
	docker compose -f tests.compose.yaml --env-file test.env build
	docker compose -f tests.compose.yaml --env-file test.env up
	docker compose -f tests.compose.yaml --env-file test.env down -v 

docker_prod:
	docker compose -f prod.compose.yaml --env-file prod.env build
	docker compose -f prod.compose.yaml --env-file prod.env up

clean:
	go clean
	rm ./target/${BINARY_NAME}-linux