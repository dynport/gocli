default: build test

build:
	go get ./...

clean:
	rm -f bin/*

test: clean build
	go test
