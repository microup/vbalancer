FROM golang:latest

ENV ConfigFile="./config/config.yaml"

ENV ProxyPort 8080
EXPOSE 8080:8080

LABEL maintainer="<contact@microup.ru>"

WORKDIR /app

COPY go.mod .
COPY go.sum .
COPY . .

RUN go mod download
RUN go build -o vbalancer cmd/vbalancer/vbalancer.go
RUN find . -name "*.go" -type f -delete

CMD ["./vbalancer"]
