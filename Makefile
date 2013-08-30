default: test install

build:
	go build -o ./bin/gocli

clean:
	rm -f bin/*

install:
	go install github.com/dynport/gocli

test: clean build
	go test -v
	# ./bin/gocli
	# ./bin/gocli co
	# ./bin/gocli co sta
