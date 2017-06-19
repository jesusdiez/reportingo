.PHONY: all build run

all: build run

build:
	go build -o reportingo main.go

run:
	./reportingo
