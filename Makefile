APP=./cmd/api

.PHONY: build test test-race test-int run up down tidy checks fmt

build:
	go build ./...

test:
	go test ./...

test-race:
	gotestsum --format testname -- -race ./...

# Run integration tests: starts compose, waits for PG, runs tests inside a Go container on compose network, then tears down
test-int:
	docker compose up --build
	@bash -lc 'until docker compose exec -T postgres pg_isready -U postgres -d jungle_gaming >/dev/null 2>&1; do echo waiting...; sleep 1; done; echo "PG ready"'
	docker run --rm -v "$$(pwd)":/src -w /src --network junglegaming-test_default -e INTEGRATION_TEST=1 -e PG_HOST=postgres golang:1.26 /usr/local/go/bin/go test ./... -run TestWalletRepository_Integration -v -count=1
	docker compose down -v

run:
	go run $(APP)

up:
	docker compose up -d

down:
	docker compose down -v

fmt:
	gofmt -s -w .

checks:
	gofmt -s -w .
	go vet ./...

tidy:
	go mod tidy
