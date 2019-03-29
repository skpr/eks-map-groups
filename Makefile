#!/usr/bin/make -f

export CGO_ENABLED=0

PROJECT=github.com/skpr/eks-map-groups
VERSION=$(shell git describe --tags --always)

# Builds the project.
build:
	gox -os='linux darwin' \
	    -arch='amd64' \
	    -output='bin/eks-map-groups_{{.OS}}_{{.Arch}}' \
	    -ldflags='-extldflags "-static"' \
	    $(PROJECT)

# Run all lint checking with exit codes for CI.
lint:
	golint -set_exit_status `go list ./... | grep -v /vendor/`

# Run tests with coverage reporting.
test:
	go test -cover ./...

IMAGE=skpr/eks-map-groups

# Releases the project Docker Hub.
release:
	docker build -t ${IMAGE}:${VERSION} -t ${IMAGE}:latest .
	docker push ${IMAGE}:${VERSION}
	docker push ${IMAGE}:latest

.PHONY: *