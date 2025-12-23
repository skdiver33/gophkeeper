APP_EXECUTABLE=client
VERSION=v1.0
BUILD_DATE=$(shell date +%Y/%m/%d-%H:%M:%S)
TARGET_DIR="./bin"

prepare:
		if [ ! -d ${TARGET_DIR} ]; then mkdir ${TARGET_DIR}; fi

linux:
		GOARCH=amd64 GOOS=linux go build -o ${TARGET_DIR}/${APP_EXECUTABLE}-linux -ldflags "-X github.com/skdiver33/gophkeeper/internal/client.Version=${VERSION} -X github.com/skdiver33/gophkeeper/internal/client.BuildDate=${BUILD_DATE}" client.go
windows:
		GOARCH=amd64 GOOS=windows go build -o ${TARGET_DIR}/${APP_EXECUTABLE}-windows -ldflags "-X github.com/skdiver33/gophkeeper/internal/client.Version=${VERSION} -X github.com/skdiver33/gophkeeper/internal/client.BuildDate=${BUILD_DATE}" client.go
darwin:
		GOARCH=amd64 GOOS=darwin go build -o ${TARGET_DIR}/${APP_EXECUTABLE}-darwin -ldflags "-X github.com/skdiver33/gophkeeper/internal/client.Version=${VERSION} -X github.com/skdiver33/gophkeeper/internal/client.BuildDate=${BUILD_DATE}" client.go

all: prepare linux windows darwin
clean:
		go clean
		rm -rf ${TARGET_DIR}
		