# Application name.
APP = rproxy
# Current branch.
BR = `git name-rev --name-only HEAD`
# Build version from nearest git tag.
VER = `git describe --tags --abbrev=0`
# Docker image tag.
TAG = $(VER)-$(BR)

.PHONY: build
build:
	CGO_ENABLED=0 go build -o $(APP) main.go

.PHONY: build-image
build-image:
	CGO_ENABLED=0 go build -o $(APP) main.go && docker build -t $(APP):$(TAG) .
