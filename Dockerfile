FROM --platform=$BUILDPLATFORM golang:1.26 AS build

ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH
ARG VERSION=dev
ARG GIT_HASH=unknown
ARG BUILD_TIME=unknown

RUN apt-get update -y && apt-get install -y libsqlite3-dev gcc-aarch64-linux-gnu

WORKDIR /go/src/app
COPY . .

RUN if [ "$TARGETARCH" = "arm64" ]; then \
      CC=aarch64-linux-gnu-gcc; \
    else \
      CC=gcc; \
    fi && \
    CC=$CC CGO_ENABLED=1 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build \
      -ldflags "-X github.com/ssalvatori/zbot-telegram/zbot.version=${VERSION} \
                -X github.com/ssalvatori/zbot-telegram/zbot.buildTime=${BUILD_TIME} \
                -X github.com/ssalvatori/zbot-telegram/zbot.gitHash=${GIT_HASH}" \
      -o zbot-telegram

FROM debian:bookworm-slim

RUN apt-get update -y && \
    apt-get install -y --no-install-recommends ca-certificates libsqlite3-0 && \
    update-ca-certificates --verbose && \
    rm -rf /var/lib/apt/lists/* && \
    useradd -r -s /usr/sbin/nologin appuser

WORKDIR /app
COPY --from=build /go/src/app/zbot-telegram /app/zbot-telegram

USER appuser

CMD ["/app/zbot-telegram"]
