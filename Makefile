.PHONY: all build run clean

all: build

build:
	go build -o tun ./cmd/tunShiffer/main.go

run: build
	sudo ./tun

clean:
	rm -f tun