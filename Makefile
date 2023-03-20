.DEFAULT_GOAL := build
fmt:
	go fmt ./...
.PHONY:fmt
lint: fmt
	golint ./...
.PHONY:lint
vet: lint
	go vet ./...
.PHONY:vet
test: vet
	go test
.PHONY: test
build: test
	go build main.go
.PHONY:build
