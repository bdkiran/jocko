#This needs to get rebuilt....
#For now just build manually
BUILD_PATH := /nolan
DOCKER_TAG := latest

all: test

deps:
	@go mod download

vet:
	@go list ./... | grep -v vendor | xargs go vet

build: deps
	@go build -o $(BUILD_PATH) main.go

release:
	@which goreleaser 2>/dev/null || go get -u github.com/goreleaser/goreleaser
	@goreleaser

clean:
	@rm -rf dist

build-docker:
	@docker build -t travisjeffery/jocko:$(DOCKER_TAG) .

generate:
	@go generate

test:
	@go test -v ./...

test-race:
	@go test -v -race -p=1 ./...

.PHONY: test-race test build-docker clean release build deps vet all

hi:
	echo "hi"
