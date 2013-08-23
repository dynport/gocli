default: test

build:
	go build -o ./bin/gocli

clean:
	rm -f bin/*

install:
	go install github.com/dynport/gocli

test: clean build
	VERBOSE=true go test
	# ./bin/gocli
	# ./bin/gocli co
	# ./bin/gocli co sta
