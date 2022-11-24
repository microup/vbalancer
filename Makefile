PROJECT_NAME = "vbalancer"

.PHONY: all dep build clean test lint

all: lint test race build

lint: ## Lint the files
  golangci-lint run -v ./...

test: ## Run unittests
  @go test -short ./...

race: dep ## Run data race detector
  @go test -race -short ./...

dep: ## Get the dependencies
  @go get -v -d ./...

build: 
   @go build -o bin/$(PROJECT_NAME) cmd/$(PROJECT_NAME)/main.go  
