GOPATH:=$(shell pwd)/gopath

all:
	go run app/main.go

build:
	go build app/main.go

deps:
	go install github.com/gorilla/mux
	go install github.com/gorilla/context
	go install github.com/mattn/go-sqlite3

fmt:
	gofmt -e -s -w ./app

clean:
	rm gopath/pkg gopath/bin -rf
