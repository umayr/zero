.PHONY: default deps clean fmt pretest lint lint-test list vet build test all
SHELL := /bin/bash
BINARY_CLIENT=zero-client
BINARY_SERVER=zero-server

VERSION=0.1.0
BUILD_TIME=`date +%FT%T%z`

BRANCH=`git rev-parse --abbrev-ref HEAD`
COMMIT=`git rev-parse --short HEAD`

LDFLAGS_CLIENT="-X ${BINARY_CLIENT}.version ${VERSION} -X ${BINARY_CLIENT}.buildtime ${BUILD_TIME} -X ${BINARY_CLIENT}.branch ${BRANCH} -X ${BINARY_CLIENT}.commit ${COMMIT}"
LDFLAGS_SERVER="-X ${BINARY_SERVER}.version ${VERSION} -X ${BINARY_SERVER}.buildtime ${BUILD_TIME} -X ${BINARY_SERVER}.branch ${BRANCH} -X ${BINARY_SERVER}.commit ${COMMIT}"

GLIDE := $(shell glide --version 2>/dev/null)
REFLEX := $(shell reflex --version 2>/dev/null)

default: build

deps:
ifdef GLIDE
	@glide install
else
	@echo "Glide is not installed."
	@echo 'Reflex is not installed'
    @echo 'Installing Glide...'
    @go get github.com/Masterminds/glide
    @glide install
endif

clean:
	@if [ -d ./build ] ; then rm -rf ./build ; fi

pretest:
	@gofmt -d $$(find . -type f -name '*.go' -not -path "./vendor/*") 2>&1 | read; [ $$? == 1 ]

lint-test:
	@go get -v github.com/golang/lint/golint
	@golint ./... | grep -v vendor/ 2>&1 || true

vet:
	@go vet $(go list -f '{{ .ImportPath }}' ./... | grep -v vendor/)

test: pretest vet lint-test
	@go test -v $$(go list -f '{{ .ImportPath }}' ./... | grep -v vendor/) -p=1

build: clean deps test
	@go build -x -ldflags ${LDFLAGS_SERVER} -o ./build/${BINARY_SERVER} ./cmd/server
	@go build -x -ldflags ${LDFLAGS_CLIENT} -o ./build/${BINARY_CLIENT} ./cmd/client

fmt:
	@gofmt -w $$(find . -type f -name '*.go' -not -path "./vendor/*")

lint:
	@go get -v github.com/golang/lint/golint
	@golint ./... | grep -v vendor/

list:
	@go list -f '{{ .ImportPath }}' ./... | grep -v vendor/

watch:
ifdef REFLEX
	@reflex -r '\.go$$' -R 'app/|vendor/|\.idea/' -s -v -- go run cmd/server/main.go
else
	@echo 'Reflex is not installed'
	@echo 'Installing Reflex...'
	@go get github.com/cespare/reflex
	@reflex -r '\.go$$' -R 'app/|vendor/|\.idea/' -s -v -- go run cmd/server/main.go
endif