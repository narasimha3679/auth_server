.PHONY: run build test clean

# Default target
all: build

# Build the application
build:
	go build -o bin/auth-server main.go

# Run the application
run:
	go run main.go

# Run tests
test:
	go test ./... -v

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Install dependencies
deps:
	go mod tidy

# Create database
create-db:
	psql -U postgres -c "CREATE DATABASE auth_db;"

# Drop database
drop-db:
	psql -U postgres -c "DROP DATABASE IF EXISTS auth_db;"

# Help
help:
	@echo "Available commands:"
	@echo "  make build      - Build the application"
	@echo "  make run        - Run the application"
	@echo "  make test       - Run tests"
	@echo "  make clean      - Clean build artifacts"
	@echo "  make deps       - Install dependencies"
	@echo "  make create-db  - Create PostgreSQL database"
	@echo "  make drop-db    - Drop PostgreSQL database"