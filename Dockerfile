FROM golang:1.17 AS build

RUN apt-get install git libsqlite3-0

ARG OS=linux
ARG ARCH=amd64

WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=1 GOOS=${OS} GOARCH=${ARCH} go build -ldflags "-X github.com/ssalvatori/zbot-telegram/zbot.version=`git describe --tags` -X github.com/ssalvatori/zbot-telegram/zbot.buildTime=`TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ'` -X github.com/ssalvatori/zbot-telegram/zbot.gitHash=`git rev-parse --short HEAD`" -o zbot-telegram-${OS}-${ARCH}

FROM debian:buster-slim

ARG OS=linux
ARG ARCH=amd64

WORKDIR /app
RUN apt-get update -y && \
    apt-get install ca-certificates -y
RUN update-ca-certificates --verbose

COPY --from=build /go/src/app/zbot-telegram-${OS}-${ARCH} /app/zbot-telegram

CMD ["/app/zbot-telegram"]