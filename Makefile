all: deps fmt docs build

deps:
	@gom install

docs:
	@gocco ./*.go

fmt:
	@go fmt ./...

build:
	@gom exec gox -osarch "darwin/amd64 linux/amd64" -output "./bin/shutter_{{.OS}}.{{.Arch}}"

.PHONY: all docs build fmt deps
