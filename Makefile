default: test

build:
	go build -o ./bin/gocli

test: build
	go test -v
	./bin/gocli
	./bin/gocli co
	./bin/gocli co sta
