VERSION=`git describe --abbrev=7 --always --tags`

default: build

build:
	cp ~/.ssh/id_rsa ./
	bash -c "VERSION=${VERSION} docker-compose build"
	rm -fr ./id_rsa

push: build
	bash -c "VERSION=${VERSION} docker-compose push"