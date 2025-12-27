.PHONY: test test-cover lint fmt vet tidy build run clean check help

BIN_DIR=bin
CMD_DIR=cmd

test:
	go test -v ./...

test-cover:
	go test -v -cover ./...

lint:
	golangci-lint run

fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	go mod tidy

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/api ./$(CMD_DIR)/api

run:
	CONFIG_ENV_PATH=./$(CMD_DIR)/$(service)  go run ./$(CMD_DIR)/$(service)

clean:
	rm -rf $(BIN_DIR)

check: fmt vet lint test

help:
	@echo "Available targets:"
	@echo "  test        - Run all tests"
	@echo "  test-cover  - Run tests with coverage"
	@echo "  lint        - Run golangci-lint"
	@echo "  fmt         - Format Go code"
	@echo "  vet         - Run go vet"
	@echo "  tidy        - Clean go.mod"
	@echo "  build       - Build API binary"
	@echo "  run         - Run service (usage: make run service=api)"
	@echo "  clean       - Remove build artifacts"
	@echo "  check       - Run fmt, vet, lint, and test"
	@echo "  help        - Show this help"
