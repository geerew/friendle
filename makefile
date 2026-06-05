
## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## dev: run API and UI dev servers together (single terminal)
.PHONY: dev
dev:
	@sh scripts/dev.sh

## test: run all tests
.PHONY: test
test:
	go test -tags dev -v ./...

## build: build the application
.PHONY: build
build:
	@echo "Building frontend..."
	(cd ui && pnpm install && pnpm run build)
	@echo "Building backend..."
	@VERSION=$$(git describe --tags --exact-match 2>/dev/null || echo "dev"); \
	COMMIT=$$(git rev-parse --short HEAD 2>/dev/null || echo "unknown"); \
	go build -ldflags "-X github.com/geerew/friendle/version.Version=$$VERSION -X github.com/geerew/friendle/version.Commit=$$COMMIT" -o friendle .
	@echo "Build complete: friendle"
