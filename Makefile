
BINARY        := zbot-telegram
DOCKER_IMAGE  := zbot-telegram-build
DOCKER_BINARY := zbot-telegram-linux-amd64
DOCKER_PLATFORMS := linux/amd64,linux/arm64

VERSION   := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_HASH  := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ')

LDFLAGS := -X github.com/ssalvatori/zbot-telegram/zbot.version=$(VERSION) \
           -X github.com/ssalvatori/zbot-telegram/zbot.gitHash=$(GIT_HASH) \
           -X github.com/ssalvatori/zbot-telegram/zbot.buildTime=$(BUILD_TIME)

.DEFAULT_GOAL := build

.PHONY: build release release-snapshot build-docker build-docker-multiarch build-docker-push test coverage lint clean clean-docker help

build: ## Build the binary
	CGO_ENABLED=1 go build -ldflags "$(LDFLAGS)" -o $(BINARY)

release: ## Build and release all artifacts using goreleaser
	goreleaser release --clean

release-snapshot: ## Build a snapshot release locally (no publish)
	goreleaser release --snapshot --clean

build-docker: ## Build image for current platform
	docker build -t $(DOCKER_IMAGE) --build-arg OS=linux --build-arg ARCH=amd64 .

build-docker-multiarch: ## Build and load multi-arch image (linux/amd64 + linux/arm64) using buildx
	docker buildx build \
		--platform $(DOCKER_PLATFORMS) \
		--build-arg OS=linux \
		-t $(DOCKER_IMAGE):$(VERSION) \
		-t $(DOCKER_IMAGE):latest \
		--load \
		.

build-docker-push: ## Build and push multi-arch image to registry (set DOCKER_IMAGE to full registry path)
	docker buildx build \
		--platform $(DOCKER_PLATFORMS) \
		--build-arg OS=linux \
		-t $(DOCKER_IMAGE):$(VERSION) \
		-t $(DOCKER_IMAGE):latest \
		--push \
		.

test: ## Run all tests
	go test -race ./...

coverage: ## Run tests and show coverage report
	go test -race -coverpkg=./... -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -func=coverage.out | tail -1

coverage-html: coverage ## Open HTML coverage report in browser
	go tool cover -html=coverage.out

lint: ## Run go vet
	go vet ./...

clean: ## Remove build artifacts
	rm -f $(BINARY) coverage.out || true

clean-docker: ## Remove Docker image and containers
	docker rmi $(DOCKER_IMAGE) || true
	docker rm -f $(DOCKER_BINARY) || true

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'