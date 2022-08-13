BUILD_ENV := CGO_ENABLED=0
LDFLAGS=-v -a -ldflags '-s -w' -gcflags="all=-trimpath=${PWD}" -asmflags="all=-trimpath=${PWD}"

TARGET_EXEC := cf

.PHONY: all setup build-linux build-osx build-windows

all: setup build-linux build-osx build-windows

setup:
	mkdir -p build

build-osx:
	${BUILD_ENV} GOARCH=amd64 GOOS=darwin go build ${LDFLAGS} -o build/${TARGET_EXEC}_darwin_amd64
	${BUILD_ENV} GOARCH=arm64 GOOS=darwin go build ${LDFLAGS} -o build/${TARGET_EXEC}_darwin_arm64

build-linux:
	${BUILD_ENV} GOARCH=amd64 GOOS=linux go build ${LDFLAGS} -o build/${TARGET_EXEC}_linux_amd64
	${BUILD_ENV} GOARCH=arm64 GOOS=linux go build ${LDFLAGS} -o build/${TARGET_EXEC}_linux_arm64
	${BUILD_ENV} GOARCH=386 GOOS=linux go build ${LDFLAGS} -o build/${TARGET_EXEC}_linux_386

build-windows:
	${BUILD_ENV} GOARCH=amd64 GOOS=windows go build ${LDFLAGS} -o build/${TARGET_EXEC}_windows_amd64.exe
	${BUILD_ENV} GOARCH=arm64 GOOS=windows go build ${LDFLAGS} -o build/${TARGET_EXEC}_windows_arm64.exe
	${BUILD_ENV} GOARCH=386 GOOS=windows go build ${LDFLAGS} -o build/${TARGET_EXEC}_windows_386.exe