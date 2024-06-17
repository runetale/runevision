FROM golang:1.22.1 AS build-env

ARG GITHUB_TOKEN

ENV GOOS=linux
ENV CGO_ENABLED=0

WORKDIR /go/src/github.com/runetale/runevison

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v .

FROM kalilinux/kali-rolling

# Update
RUN apt -y update && DEBIAN_FRONTEND=noninteractive apt -y dist-upgrade && apt -y autoremove && apt clean

# Install common and useful tools
RUN apt -y install curl wget vim git net-tools whois netcat-traditional pciutils usbutils

COPY --from=build-env /go/src/github.com/runetale/runevison/runevison /runevison
RUN chmod u+x /runevison

COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod u+x ./docker-entrypoint.sh