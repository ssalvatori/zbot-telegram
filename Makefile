
DOCKER_IMAGE_NAME=zbot-telegram-build

build:
	CGO_ENABLED=1 go build -ldflags "-X github.com/ssalvatori/zbot-telegram/zbot.version=`git describe --tags` -X github.com/ssalvatori/zbot-telegram/zbot.gitHash=`git rev-parse --short HEAD` -X github.com/ssalvatori/zbot-telegram/zbot.buildTime=`TZ=UTC date -u '+%Y-%m-%dTd%H:%M:%SZ'`" -o zbot-telegram

build-docker:
	docker build -t $(DOCKER_IMAGE_NAME) --build-arg OS=linux --build-arg ARCH=amd64 .
	docker create -ti --name zbot-telegram-linux-amd64 $(DOCKER_IMAGE_NAME) bash
	docker cp zbot-telegram-linux-amd64:/go/src/app/zbot-telegram-linux-amd64 .
	docker rm -f zbot-telegram-linux-amd64

.PHONY: clean clean-docker

clean:
	rm -rf zbot-telegram || true

clean-docker:
	docker rmi $(DOCKER_IMAGE_NAME) || true
	docker rm -f zbot-telegram-linux-amd64 || true