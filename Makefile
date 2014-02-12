all: deps fmt docs build

deps:
	@gom install

docs:
	@gocco ./*.go

fmt:
	@go fmt ./...

build:
	@gom exec gox -osarch "darwin/amd64 linux/amd64" -output "./bin/shutter_{{.OS}}.{{.Arch}}"

setup:
	@go get github.com/mitchellh/gox
	@gox -build-toolchain
	@go get github.com/mattn/gom

.PHONY: all docs build fmt deps setup
