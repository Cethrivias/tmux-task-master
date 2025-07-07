VERSION=v0.3.0
BUILD=`date +%FT%T%z`
BUILD_FLAGS=-ldflags "-X main.version=${VERSION} -X main.build=${BUILD}"

build:
	go build -ldflags "-X main.version=${VERSION} -X main.build=${BUILD}" -o ./bin/ttm

install:
	go install ${BUILD_FLAGS}

test:
	go test ./... -v

test.integ:
	go test main_test.go -v

format:
	go fmt ./...
