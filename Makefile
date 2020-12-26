
DOCKER_IMAGE_NAME=zbot-telegram-build
build-linux-amd64:
	docker build -t $(DOCKER_IMAGE_NAME) --build-arg OS=linux --build-arg ARCH=amd64 .
	docker create -ti --name zbot-telegram-linux-amd64 $(DOCKER_IMAGE_NAME) bash
	docker cp zbot-telegram-linux-amd64:/go/src/app/zbot-telegram-linux-amd64 .
	docker rm -f zbot-telegram-linux-amd64
.PHONY: clean
clean:
	docker rmi $(DOCKER_IMAGE_NAME) || true
	docker rm -f zbot-telegram-linux-amd64 || true