.DEFAULT_GOAL := build

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit:
	go mod verify
	go vet ./...

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## test: run all tests
.PHONY: test
test:
	GH_TOKEN= GITHUB_TOKEN= go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	GH_TOKEN= GITHUB_TOKEN= go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
ifndef CI
	go tool cover -html=/tmp/coverage.out
endif

## build: build the application
.PHONY: build
build:
	go build

.PHONY: clean
clean:
	rm -rf dist
	rm -f actions-dependency-graph
