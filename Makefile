DEPS = $(go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)

all: fmt docs build

docs:
	@gocco ./*.go

fmt:
	@go fmt ./...

build:
	@gox -osarch "darwin/amd64 linux/amd64" -output "./bin/shutter_{{.OS}}.{{.Arch}}"

.PHONY: all docs build fmt
