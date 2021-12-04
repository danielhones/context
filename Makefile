.PHONY: build

build:
	env go build -ldflags="-extldflags=-static" -o bin/context 

clean:
	rm -rf ./bin
