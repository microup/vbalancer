FROM golang:1.23rc2

ENV ConfigFile="config.yaml"
ENV ProxyPort=8080

EXPOSE 8080:8080

LABEL maintainer="<contact@microup.ru>"

WORKDIR /vbalancer

COPY cmd/ cmd/
COPY internal/ internal/
COPY config/config.yaml .
COPY go.mod .
COPY go.sum .
COPY Makefile .
COPY LICENSE .
COPY readme.md .

RUN go mod download
RUN make init
RUN go build -o vbalancer cmd/vbalancer/vbalancer.go

RUN rm -rf cmd/
RUN rm -rf internal/
RUN rm -rf vendor/
RUN rm -rf mocks/

CMD ["./vbalancer"]
