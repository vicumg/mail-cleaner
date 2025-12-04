.PHONY: all build clean test run analyze help

# Build all binaries
all: build

# Build both commands
build:
	@echo "Building mail-cleaner..."
	@go build -o mail-cleaner ./cmd/mail-cleaner
	@echo "Building analyze..."
	@go build -o analyze ./cmd/analyze
	@echo "Done!"

# Build only mail-cleaner
mail-cleaner:
	@go build -o mail-cleaner ./cmd/mail-cleaner

# Build only analyze
analyze-tool:
	@go build -o analyze ./cmd/analyze

# Run tests
test:
	@go test ./...

# Clean binaries
clean:
	@rm -f mail-cleaner analyze
	@echo "Cleaned binaries"

# Run mail-cleaner (example: make run SERVICE=ukrnet RULES=rules.json)
run:
	@go run ./cmd/mail-cleaner $(SERVICE) $(RULES)

# Run analyze (example: make analyze LOG=spam_classification.log)
analyze:
	@go run ./cmd/analyze $(LOG)

# Install both commands to GOPATH/bin
install:
	@go install ./cmd/mail-cleaner
	@go install ./cmd/analyze

# Show help
help:
	@echo "Available targets:"
	@echo "  make build           - Build both mail-cleaner and analyze"
	@echo "  make mail-cleaner    - Build only mail-cleaner"
	@echo "  make analyze-tool    - Build only analyze"
	@echo "  make test            - Run tests"
	@echo "  make clean           - Remove binaries"
	@echo "  make run SERVICE=<name> RULES=<file> - Run mail-cleaner"
	@echo "  make analyze LOG=<file>              - Run analyze"
	@echo "  make install         - Install both commands to GOPATH/bin"
