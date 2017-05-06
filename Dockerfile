FROM golang:1.8.1-alpine
MAINTAINER Stefano Salvatori <ssalvatori@gmail.com>

#ADD https://github.com/ssalvatori/zbot-telegram-go/releases/download/v1.0.3/zbot-telegram-go-linux-amd64 /bin/zbot

#RUN apk add --update \
#    file \
#    build-base \
#    && rm -rf /var/cache/apk/*

COPY zbot-telegram-go-linux-amd64 /go/bin/zbot-telegram-go-linux-amd64
RUN chmod +x /go/bin/zbot-telegram-go-linux-amd64
CMD ["/go/bin/zbot-telegram-go-linux-amd64"]
