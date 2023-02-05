PROJECT_NAME = "vbalancer"

.PHONY: all dep build test lint mocks

all: lint test race build

pre-push: fmt lint test race 

fmt:
  go fmt ./...

build-mocks:
  go get github.com/golang/mock/gomock
  go install github.com/golang/mock/mockgen
  
mocks:
  mockgen -destination=mocks/mock_peer.go -package=mocks -source=./internal/proxy/peer/peer.go Peer
  mockgen -destination=mocks/mock_vlog.go -package=mocks -source=./internal/vlog/vlog.go ILog

lint: 
  go vet ./...
  golangci-lint run -v ./...

test: 
  go test -v ./...

race: dep ## Run data race detector
  go test -race -v ./...

dep: ## Get the dependencies
  go mod tidy

build: 
  go build -o build/$(PROJECT_NAME) cmd/$(PROJECT_NAME)/$(PROJECT_NAME).go
