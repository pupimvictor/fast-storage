.PHONY: install build pack run stop dkr-run-redis

TAG?=$(shell date +%s)

export TAG

install:
		echo $(TAG)

build:  install
		GOOS=linux /usr/local/go1.11/bin/go build -ldflags="-s -w -X main.version=$(TAG)" -o ./cmd/fast-storage ./cmd/

pack:   build
		docker build -t fast-storage:$(TAG) .
		docker tag fast-storage:$(TAG) fast-storage:latest

run:
		docker run --name fast-storage-$(TAG)  --link some-redis:redis --rm redis sh -c 'exec redis-cli -h "$$REDIS_PORT_6379_TCP_ADDR" -p "$$REDIS_PORT_6379_TCP_PORT"' -d fast-storage:latest
		docker run fast-storage-$(TAG)
		docker ps
		docker logs fast-storage-$(TAG)

test:   pack run

stop:
		docker stop fast-storage
		docker rm fast-storage

dkr-run-redis:
	docker run --name some-redis -d redis


ship: build
