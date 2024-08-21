build-and-run-linux: build run-linux

build:
	go mod tidy
	go mod vendor
	go build ./cmd/app/main.go 

run-linux:
	./main -configpath="./cmd/app/config.yaml" -ipport="${HOSTNAME}:8008" -servicename="${HOSTNAME}service" -nodename="${HOSTNAME}node"
