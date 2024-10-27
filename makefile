build:
	go build -o ./bin/ttm

install:
	go install

test.integ:
	go test main_test.go -v

format:
	go fmt ./...
