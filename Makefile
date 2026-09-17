APP=./cmd/api

.PHONY: build test run up down tidy

build:
	go build ./...

test:
	go test ./...

run:
	go run $(APP)

up:
	docker compose up -d

down:
	docker compose down -v

tidy:
	go mod tidy
