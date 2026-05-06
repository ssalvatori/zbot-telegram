FROM --platform=$BUILDPLATFORM golang:1.26 AS build

ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

RUN apt-get update -y && apt-get install -y git libsqlite3-dev gcc-aarch64-linux-gnu

WORKDIR /go/src/app
COPY . .

RUN if [ "$TARGETARCH" = "arm64" ]; then \
      CC=aarch64-linux-gnu-gcc; \
    else \
      CC=gcc; \
    fi && \
    CC=$CC CGO_ENABLED=1 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build \
      -ldflags "-X github.com/ssalvatori/zbot-telegram/zbot.version=$(git describe --tags 2>/dev/null || echo dev) \
                -X github.com/ssalvatori/zbot-telegram/zbot.buildTime=$(TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ') \
                -X github.com/ssalvatori/zbot-telegram/zbot.gitHash=$(git rev-parse --short HEAD 2>/dev/null || echo unknown)" \
      -o zbot-telegram

FROM debian:bookworm-slim

WORKDIR /app
RUN apt-get update -y && \
    apt-get install -y ca-certificates libsqlite3-0 && \
    update-ca-certificates --verbose && \
    rm -rf /var/lib/apt/lists/*

COPY --from=build /go/src/app/zbot-telegram /app/zbot-telegram

CMD ["/app/zbot-telegram"]
