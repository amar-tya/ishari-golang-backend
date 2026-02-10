
.PHONY: run build test clean help migrate-up migrate-down migrate-force migrate-version migrate-create

# Variables
BINARY_NAME=api
MAIN_PATH=cmd/api/main.go
BUILD_DIR=bin
MIGRATIONS_DIR=migrations
MIGRATE_BIN=migrate
MIGRATE_DB_URL=postgres://postgres@localhost:5432/ishari?sslmode=disable

# Default target
help:
	@echo "Available commands:"
	@echo "  make run       - Run the application"
	@echo "  make build     - Build the application"
	@echo "  make test      - Run tests (concise)"
	@echo "  make test-v    - Run tests with verbose output"
	@echo "  make test-race - Run tests with race detection"
	@echo "  make coverage  - Run tests with coverage report"
	@echo "  make clean     - Clean build artifacts"
	@echo "  make deps      - Download dependencies"
	@echo "  make lint      - Run linter"

# Run the application
run:
	go run $(MAIN_PATH)

# Build the application
build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Binary created at $(BUILD_DIR)/$(BINARY_NAME)"

# Run tests
test:
	go test ./...

# Run tests with verbose output
test-v:
	go test -v ./...

# Run tests with race detection
test-race:
	go test -race -v ./...

# Run tests with coverage
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	@echo "Cleaned build artifacts"

# Download dependencies
deps:
	go mod download
	go mod tidy

# Run linter (requires golangci-lint)
lint:
	golangci-lint run

# Run the built binary
start: build
	./$(BUILD_DIR)/$(BINARY_NAME)

# Migration commands
migrate-up:
	$(MIGRATE_BIN) -path "$(MIGRATIONS_DIR)" -database "$(MIGRATE_DB_URL)" up

migrate-down:
	$(MIGRATE_BIN) -path "$(MIGRATIONS_DIR)" -database "$(MIGRATE_DB_URL)" down 1

migrate-force:
	$(MIGRATE_BIN) -path "$(MIGRATIONS_DIR)" -database "$(MIGRATE_DB_URL)" force $(VERSION)

migrate-version:
	$(MIGRATE_BIN) -path "$(MIGRATIONS_DIR)" -database "$(MIGRATE_DB_URL)" version

migrate-create:
	$(MIGRATE_BIN) create -ext sql -dir "$(MIGRATIONS_DIR)" -seq $(NAME)

migrate-drop:
	$(MIGRATE_BIN) -path "$(MIGRATIONS_DIR)" -database "$(MIGRATE_DB_URL)" drop -f

migrate-reset:
	$(MIGRATE_BIN) -path "$(MIGRATIONS_DIR)" -database "$(MIGRATE_DB_URL)" drop -f
	$(MIGRATE_BIN) -path "$(MIGRATIONS_DIR)" -database "$(MIGRATE_DB_URL)" up
