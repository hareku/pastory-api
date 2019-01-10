.PHONY: deps clean build

deps:
	go get -u ./...

clean:
	rm -rf ./hello-world/hello-world

build:
	cross-env GOOS=linux GOARCH=amd64 go build -o main main.go
