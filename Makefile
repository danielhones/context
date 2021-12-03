.PHONY: build

build:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-extldflags=-static" -o bin/context context.go

clean:
	rm -rf ./bin
