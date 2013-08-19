default: test

build:
	go build -o ./bin/gocli

clean:
	rm -f bin/*

test: clean build
	go test -v
	# ./bin/gocli
	# ./bin/gocli co
	# ./bin/gocli co sta
