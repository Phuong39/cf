BUILD_ENV := CGO_ENABLED=0
LDFLAGS=-v -a -ldflags '-s -w' -gcflags="all=-trimpath=${PWD}" -asmflags="all=-trimpath=${PWD}"
VERSION=v0.3.1

TARGET_EXEC := cf

.PHONY: all setup build-linux build-osx build-windows

all: setup build-linux build-osx build-windows

setup:
	mkdir -p build/linux
	mkdir -p build/osx
	mkdir -p build/windows

build-linux:
	${BUILD_ENV} GOARCH=amd64 GOOS=linux go build ${LDFLAGS} -o build/linux/${TARGET_EXEC}-${VERSION}-linux-amd64

build-osx:
	${BUILD_ENV} GOARCH=amd64 GOOS=darwin go build ${LDFLAGS} -o build/osx/${TARGET_EXEC}-${VERSION}-darwin-amd64

build-windows:
	${BUILD_ENV} GOARCH=amd64 GOOS=windows go build ${LDFLAGS} -o build/windows/${TARGET_EXEC}-${VERSION}-windows-amd64.exe