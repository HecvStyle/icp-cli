USER=hecvstyle
NAME=icp
VERSION=1.0.0
REGISTRY=registry.cn-hangzhou.aliyuncs.com

.PHONY: clear build start push

build: buildGo build-version

clear:
	docker stop ${NAME};docker rm ${NAME}
	docker rmi $(docker images -q -f dangling=true)

buildLocal:
	go build -o ${NAME}

buildGo:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${NAME}

build-version: buildGo
	docker build -t ${REGISTRY}/${USER}/${NAME}:${VERSION} .

tag-latest:
	docker tag ${REGISTRY}/${USER}/${NAME}:${VERSION} ${REGISTRY}/${USER}/${NAME}:latest

start:
	docker run -it -v ${PWD}/config.yaml:/app/config.yaml --rm --name ${NAME} ${REGISTRY}/${USER}/${NAME}:${VERSION} /bin/bash

push:   buildGo build-version tag-latest
	docker push ${REGISTRY}/${USER}/${NAME}:latest;
	docker push ${REGISTRY}/${USER}/${NAME}:${VERSION}

pull:
	docker pull  ${REGISTRY}/${USER}/${NAME}:latest

restart:
	docker run --name ${NAME} --restart=always -p 8899:8899 -v ${PWD}/config.yaml:/app/config.yaml -v ${PWD}/log:/app/log ${REGISTRY}/${USER}/${NAME}:latest