.PHONY: build test clean install release

# Build variables
BINARY_NAME=envtool
# Fall back to "dev" when not in a git repo
GIT_VERSION=$(shell git describe --tags --always --dirty 2>/dev/null)
VERSION:=$(if $(strip $(GIT_VERSION)),$(GIT_VERSION),dev)
BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

.PHONY: tidy
tidy:
	go mod tidy
	gomod2nix

# Build the application
build:
	go build $(LDFLAGS) -o $(BINARY_NAME) main.go

# Run all tests
test:
	go test -v ./...

# Run integration tests
integration-test:
	INTEGRATION_TEST=true go test -v ./tests

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -rf dist/

# Install the application
install:
	go install $(LDFLAGS)

# Build for multiple platforms
build-all:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)_linux_amd64 main.go
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)_darwin_amd64 main.go
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)_windows_amd64.exe main.go

# Create a release (requires goreleaser)
release:
	goreleaser release --rm-dist