# Purtypics Makefile

# Variables
BINARY_NAME=purtypics
GO=go
GOFLAGS=-v
LDFLAGS=-s -w
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
PREFIX ?= /usr/local
DATADIR ?= $(PREFIX)/share

# Build variables
MAIN_PACKAGE=.
BUILD_DIR=build
DIST_DIR=dist

# Default target
.DEFAULT_GOAL := build

# Build targets
.PHONY: build
build: ## Build the binary
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS) -X main.version=$(VERSION)" -o $(BINARY_NAME) $(MAIN_PACKAGE)

.PHONY: build-all
build-all: ## Build for all platforms
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PACKAGE)
	GOOS=darwin GOARCH=arm64 $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PACKAGE)
	GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PACKAGE)
	GOOS=linux GOARCH=arm64 $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PACKAGE)
	GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PACKAGE)

.PHONY: install
install: build ## Install binary and themes
	$(GO) install $(GOFLAGS) $(MAIN_PACKAGE)
	@echo "Installing themes to $(DATADIR)/purtypics/themes..."
	install -d $(DATADIR)/purtypics/themes
	cp -r pkg/gallery/assets/themes/* $(DATADIR)/purtypics/themes/

.PHONY: uninstall
uninstall: ## Remove installed themes
	rm -rf $(DATADIR)/purtypics

# Testing targets
.PHONY: test
test: ## Run tests
	$(GO) test -v ./...

.PHONY: fmt
fmt: ## Format code
	$(GO) fmt ./...

.PHONY: vet
vet: ## Run go vet
	$(GO) vet ./...

# Dependency management
.PHONY: deps
deps: ## Download dependencies
	$(GO) mod download

.PHONY: deps-update
deps-update: ## Update dependencies
	$(GO) get -u ./...
	$(GO) mod tidy

# Clean targets
.PHONY: clean
clean: ## Clean build artifacts
	rm -f $(BINARY_NAME)
	rm -rf $(BUILD_DIR)
	rm -rf $(DIST_DIR)

# Distribution targets
.PHONY: dist
dist: clean build-all ## Create distribution packages
	@mkdir -p $(DIST_DIR)
	@for file in $(BUILD_DIR)/*; do \
		base=$$(basename $$file); \
		tar -czf $(DIST_DIR)/$$base.tar.gz -C $(BUILD_DIR) $$base; \
	done
	@echo "Distribution packages created in $(DIST_DIR)/"

# Release targets
.PHONY: release
release: ## Tag a new release, update README installer links, and push
	@if [ -z "$(TAG)" ]; then echo "Usage: make release TAG=v1.4.0"; exit 1; fi
	@VERSION=$$(echo "$(TAG)" | sed 's/^v//'); \
	sed -i'' -e '/<!-- begin:installer-links/,/<!-- end:installer-links/s/[0-9][0-9]*\.[0-9][0-9]*\.[0-9][0-9]*/'"$$VERSION"'/g' README.md; \
	git add README.md && \
	git commit -m "Update installer links to $(TAG)" && \
	git tag -a "$(TAG)" -m "Release $(TAG)" && \
	echo "Tagged $(TAG) with updated README links. Push with: git push && git push --tags"

# Help target
.PHONY: help
help: ## Show this help message
	@echo "Purtypics Makefile"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*##"; printf "\033[36m%-20s\033[0m %s\n", "Target", "Description"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)