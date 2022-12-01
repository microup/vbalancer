FROM golang:latest

EXPOSE 8010
EXPOSE 8010
EXPOSE 8010

ENV PORT=8010
ENV JWT_SECRET=drsgserg_esr%$SdSDfvsd_sdfbsdfb_DGVS_#S

LABEL maintainer="<contact@microup.ru>"

WORKDIR /app

COPY go.mod .
COPY go.sum .
COPY . .

RUN apt update && apt upgrade -y && apt install mc -y
RUN go mod download
RUN go build
RUN find . -name "*.go" -type f -delete

CMD ["./golang-gin-poc"]
