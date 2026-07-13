# Dev entrypoints. These mirror what the Nix devShell provides (BRIEFING.md §6);
# run `nix develop` first, or have the toolchain on PATH.

BINARY := terraform-provider-hcloudimage

.PHONY: default build test testacc lint fmt docs validate-examples clean

default: build

# Build the provider binary.
build:
	go build -o $(BINARY) .

# Unit tests (no network).
test:
	go test ./... -covermode=atomic -coverprofile=coverage.out

# Acceptance tests (real Hetzner, billable). Gated by TF_ACC (BRIEFING.md §8.3).
testacc:
	TF_ACC=1 go test ./... -v -timeout 120m

# Lint.
lint:
	golangci-lint run ./...

# Format Go sources.
fmt:
	gofmt -s -w .

# Generate provider docs from schema + examples (BRIEFING.md §11).
docs:
	tfplugindocs generate

# Validate the HCL examples under both terraform and tofu (BRIEFING.md §5).
# Builds the provider first so the mirror in the script has a binary.
validate-examples:
	nix build .#default
	bash scripts/validate-examples.sh

clean:
	rm -f $(BINARY) coverage.out
	rm -rf dist
