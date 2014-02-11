DEPS = $(go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)

all: deps fmt docs build

deps:
	gom install

docs:
	@gocco ./*.go

fmt:
	@go fmt ./...

build:
	@gox -osarch "darwin/amd64 linux/amd64" -output "./bin/shutter_{{.OS}}.{{.Arch}}"

.PHONY: all docs build fmt deps
