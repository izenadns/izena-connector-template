.DEFAULT_GOAL := help

WASM_COMMAND := $(patsubst %/main.go,%,$(wildcard ./cmd/*/main.go))
WASM_ARTIFACT ?= dist/$(notdir $(WASM_COMMAND)).wasm

RELEASE_ARCHIVES := $(wildcard dist/*.tar.gz)
RELEASE_SBOMS := $(wildcard dist/*.sbom.json)
SHA256 := $(shell command -v sha256sum 2>/dev/null || printf 'shasum -a 256')

.PHONY: help bootstrap
.PHONY: fmt fmt-check lint security module-check
.PHONY: test test-integration test-integration-fixtures ci
.PHONY: build artifact-path force
.PHONY: release-check release-artifacts release-snapshot release-checksums

help: ## Show the connector development commands.
	@awk 'BEGIN {FS = ":.*## "; print "Connector development commands:"} /^[a-zA-Z0-9_-]+:.*## / {printf "  %-26s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

bootstrap: ## Check the Go toolchain.
	@go_version="$$(go env GOVERSION)"; \
	test "$$go_version" = go1.27.4 || { echo "Go 1.27.4 is required (found $$go_version)" >&2; exit 1; }

fmt: ## Format Go source.
	@go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.15.0 fmt

fmt-check: ## Check Go formatting.
	@go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.15.0 fmt --diff

lint: ## Run configured Go linters.
	@go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.15.0 run

security: ## Scan Go dependencies.
	@go run github.com/google/osv-scanner/v2/cmd/osv-scanner@v2.7.1 scan source \
		--config osv-scanner.toml --no-call-analysis=go --lockfile go.mod

module-check: ## Check standalone module dependencies.
	@GOWORK=off go mod tidy -diff

test: ## Run unit tests with the race detector.
	@go test -race ./...

test-integration: ## Run integration tests with the configured upstream.
	@go test -count=1 -tags=integration ./...

test-integration-fixtures: ## Run integration tests with the repository's disposable fixtures.
	@bash dev/test-integration.sh $(MAKE) test-integration

build: $(WASM_ARTIFACT).sha256 ## Build the CGO-disabled WASI reactor and its checksum.

# Go tracks the sources and module dependencies used to build the reactor.
$(WASM_ARTIFACT): force
	@test "$(words $(WASM_COMMAND))" -eq 1 || { \
		echo "Expected exactly one cmd/*/main.go" >&2; exit 1; \
	}
	@CGO_ENABLED=0 GOOS=wasip1 GOARCH=wasm go build -mod=readonly -trimpath \
		-buildmode=c-shared -ldflags="-s -w" -o "$(WASM_ARTIFACT)" "$(WASM_COMMAND)"

force:

artifact-path: ## Print the Wasm artifact path.
	@printf '%s\n' "$(WASM_ARTIFACT)"

release-check: ## Validate the GoReleaser configuration.
	@go run github.com/goreleaser/goreleaser/v2@v2.20.1 check

release-artifacts: ## Package a tagged release locally.
	@go run github.com/goreleaser/goreleaser/v2@v2.20.1 release --clean --skip=publish
	@$(MAKE) build release-checksums

release-snapshot: ## Package a local release snapshot.
	@go run github.com/goreleaser/goreleaser/v2@v2.20.1 release --snapshot --clean
	@$(MAKE) build release-checksums

release-checksums: $(addsuffix .sha256,$(RELEASE_ARCHIVES) $(RELEASE_SBOMS)) ## Write checksums for the source archives and SBOMs.
	@test -n "$(RELEASE_ARCHIVES)" && test -n "$(RELEASE_SBOMS)" || { \
		echo "Expected source archives and SBOMs in dist/" >&2; exit 1; \
	}

%.sha256: %
	@cd "$(@D)" && checksum="$$($(SHA256) -- "$(<F)")" && printf '%s\n' "$$checksum" > "$(@F)"

ci: bootstrap fmt-check lint security module-check test test-integration-fixtures ## Run the connector CI checks.
