FROM alpine:3.8
LABEL maintainer="ssalvatori@gmail.com"

RUN apk add --update \
    wget \
    sqlite \
    jq \
    && rm -rf /var/cache/apk/*

RUN apk --no-cache add ca-certificates \
    && wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub \
    && wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.28-r0/glibc-2.28-r0.apk \
    && apk add glibc-2.28-r0.apk \
    && rm -rf /var/cache/apk/*

WORKDIR /app

COPY zbot-telegram-go-linux-amd64 /app/zbot-telegram-go-linux-amd64
RUN chmod +x /app/zbot-telegram-go-linux-amd64
CMD ["/app/zbot-telegram-go-linux-amd64"]
