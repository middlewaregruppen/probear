BINARY=probear
GOARCH=amd64
VERSION=1.16
COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
GOVERSION=$(shell go version | awk -F\go '{print $$3}' | awk '{print $$1}')
GITHUB_USERNAME=middlewaregruppen
BUILD_DIR=.
PKG_LIST=$$(go list ./... | grep -v /vendor/)
# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-X main.VERSION=${VERSION} -X main.COMMIT=${COMMIT} -X main.BRANCH=${BRANCH} -X main.GOVERSION=${GOVERSION}"

# Build the project
all: build

test:
	cd ${BUILD_DIR}; \
	go test -v; \
	cd - >/dev/null

fmt:
	cd ${BUILD_DIR}; \
	go fmt ${PKG_LIST} ; \
	cd - >/dev/null

dep:
	go get -v -d ./... ;

linux: dep
	CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BUILD_DIR}/out/${BINARY}-linux-${GOARCH} cmd/probear/main.go

rpi: dep
	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build ${LDFLAGS} -o ${BUILD_DIR}/out/${BINARY}-linux-arm cmd/probear/main.go

darwin: dep
	CGO_ENABLED=0 GOOS=darwin GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BUILD_DIR}/out/${BINARY}-darwin-${GOARCH} cmd/probear/main.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o ${BUILD_DIR}/out/${BINARY}-darwin-arm cmd/probear/main.go 

windows: dep
	CGO_ENABLED=0 GOOS=windows GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BUILD_DIR}/out/${BINARY}-windows-${GOARCH}.exe cmd/probear/main.go

docker_build:
	docker run --rm -v "${PWD}":/go/src/github.com/middlewaregruppen/probear -w /go/src/github.com/middlewaregruppen/probear golang:${GOVERSION} make fmt test
	docker build -t ghcr.io/middlewaregruppen/probear:${VERSION} .
	docker tag ghcr.io/middlewaregruppen/probear:${VERSION} ghcr.io/middlewaregruppen/probear:latest

docker_push:
	docker push ghcr.io/middlewaregruppen/probear:${VERSION}
	docker push ghcr.io/middlewaregruppen/probear:latest

docker:  docker_build docker_push

build:  linux darwin rpi windows

clean:
	-rm -rf ${BUILD_DIR}/out/

.PHONY: linux darwin windows test fmt clean