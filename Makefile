.PHONY: build test coverage clean

build:
	go build -ldflags="-extldflags=-static" -o bin/context 

test:
	go test -cover -coverprofile coverage.out

coverage: test
	go tool cover -html=coverage.out

clean:
	rm -rf ./bin
