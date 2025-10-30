.PHONY: build install test clean fmt lint docs

BINARY=terraform-provider-pipeops
VERSION=0.1.0
OS_ARCH=darwin_arm64

# Build the provider
build:
	go build -o ${BINARY}

# Install the provider locally for testing
install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/PipeOpsHQ/pipeops/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/registry.terraform.io/PipeOpsHQ/pipeops/${VERSION}/${OS_ARCH}/

# Run tests
test:
	go test -v ./...

# Run acceptance tests
testacc:
	TF_ACC=1 go test -v -timeout 120m ./...

# Clean build artifacts
clean:
	rm -f ${BINARY}
	rm -rf dist/

# Format code
fmt:
	go fmt ./...
	terraform fmt -recursive ./examples/

# Run linter
lint:
	golangci-lint run

# Generate documentation
docs:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

# Download dependencies
deps:
	go mod download
	go mod tidy

# Run go mod tidy
tidy:
	go mod tidy

# Build for multiple platforms
build-all:
	GOOS=darwin GOARCH=amd64 go build -o dist/${BINARY}_${VERSION}_darwin_amd64
	GOOS=darwin GOARCH=arm64 go build -o dist/${BINARY}_${VERSION}_darwin_arm64
	GOOS=linux GOARCH=amd64 go build -o dist/${BINARY}_${VERSION}_linux_amd64
	GOOS=linux GOARCH=arm64 go build -o dist/${BINARY}_${VERSION}_linux_arm64
	GOOS=windows GOARCH=amd64 go build -o dist/${BINARY}_${VERSION}_windows_amd64.exe

# Help
help:
	@echo "Available targets:"
	@echo "  build      - Build the provider binary"
	@echo "  install    - Install the provider locally"
	@echo "  test       - Run unit tests"
	@echo "  testacc    - Run acceptance tests"
	@echo "  clean      - Remove build artifacts"
	@echo "  fmt        - Format code"
	@echo "  lint       - Run linter"
	@echo "  docs       - Generate documentation"
	@echo "  deps       - Download dependencies"
	@echo "  build-all  - Build for all platforms"
