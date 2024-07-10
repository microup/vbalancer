PROJECT_NAME=vbalancer

.PHONY: all dep build test lint mocks

all: lint test race build

pre-push: lint test race

prepare:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/golang/mock/mockgen@v1.6.0

mocks:
	mockgen -destination=mocks/mock_peer.go -package=mocks -source=./internal/proxy/peer/peer.go Peer
	mockgen -destination=mocks/mock_vlog.go -package=mocks -source=./internal/vlog/vlog.go ILog

init:
	go mod tidy
	go mod vendor

fmt:
	go fmt ./...

lint:
	go vet ./...
	golangci-lint run -v ./...

test:
	go test -v ./...

race:
	go test -race -v ./...

build:
	go build -o build/$(PROJECT_NAME) cmd/$(PROJECT_NAME)/$(PROJECT_NAME).go

build-win:
	go build -o build/$(PROJECT_NAME).exe cmd/$(PROJECT_NAME)/$(PROJECT_NAME).go

docker-create:
	docker build --tag vbalancer . -f Dockerfile

docker-run:
	docker run --restart=always -p 8080:8080 vbalancer

docker-delete:
	docker rmi vbalancer
