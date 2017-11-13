FROM golang:1.8.1-alpine
MAINTAINER Stefano Salvatori <ssalvatori@gmail.com>

#ADD https://github.com/ssalvatori/zbot-telegram-go/releases/download/v1.0.3/zbot-telegram-go-linux-amd64 /bin/zbot

RUN apk add --update \
    wget \
    sqlite \
    && rm -rf /var/cache/apk/*

RUN apk --no-cache add ca-certificates \
    && wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://raw.githubusercontent.com/sgerrand/alpine-pkg-glibc/master/sgerrand.rsa.pub \
    && wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.25-r0/glibc-2.25-r0.apk \
    && apk add glibc-2.25-r0.apk \
    && rm -rf /var/cache/apk/*

COPY zbot-telegram-go-linux-amd64 /go/bin/zbot-telegram-go-linux-amd64
RUN chmod +x /go/bin/zbot-telegram-go-linux-amd64
CMD ["/go/bin/zbot-telegram-go-linux-amd64"]
