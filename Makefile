BINARY=fakecas
VERSION=$(shell git describe master 2&> /dev/null || echo "$${VERSION:-Unknown}")
BUILD_TIME=`date +%FT%T%z`
LDFLAGS=-ldflags "-X main.Version=${VERSION}"
BUILD_DIR=./build/
BINARY_TMPL=${BINARY}-{{.OS}}-{{.Arch}}
ARCH=386 amd64
OS=darwin linux windows

build:
	go build ${LDFLAGS} -o ${BINARY}

build-all:
	gox ${LDFLAGS} -output="${BUILD_DIR}${BINARY_TMPL}" -arch="${ARCH}" -os="${OS}"

.PHONY: build

clean:
	git clean -Xdf
