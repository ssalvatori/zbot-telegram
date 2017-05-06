clean:
	rm -rf dist/*

build:
	curl -L -O https://github.com/ssalvatori/zbot-telegram-go/releases/download/v1.0.3/zbot-telegram-go-linux-amd64
	docker build -t zbot-telegram-go .
