VERSION := 0.1.0
TARGET  := nsnitro
TEST    ?= ./...

default: test build

deps:
	go get -v -u ./...

test:
	go test -v -cover -run=$(RUN) $(TEST)

build: clean
	go build -v -o bin/$(TARGET)

clean:
	rm -rf bin/
