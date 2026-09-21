.PHONY: help run build test lint fmt vet tidy clean docker-up docker-down

help:
	@echo "Available commands:"
	@echo "  make run          Run the API"
	@echo "  make build        Build the API"
	@echo "  make test         Run tests"
	@echo "  make lint         Run linter"
	@echo "  make fmt          Format code"
	@echo "  make vet          Run go vet"
	@echo "  make tidy         Tidy dependencies"
	@echo "  make docker-up    Start Docker services"
	@echo "  make docker-down  Stop Docker services"
	@echo "  make clean        Clean build files"

run:
	go run ./cmd/api

build:
	go build -o bin/switchyard ./cmd/api

test:
	go test ./...

lint:
	golangci-lint run ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	go mod tidy

docker-up:
	docker compose up -d

docker-down:
	docker compose down

clean:
	rm -rf bin coverage.out

check:
	go fmt ./...
	go vet ./...
	golangci-lint run ./...
	go test ./...
