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

## lint: use golanci-lint to lint the app
.PHONY: lint
lint:
	golangci-lint run

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
ifndef CI
	go tool cover -html=/tmp/coverage.out
endif

## build: build the application
.PHONY: build
build:
	CGO_ENABLED=1 go build

.PHONY: clean
clean:
	rm -rf dist
	rm -f actions-dependency-graph
