.PHONY: all clean deps fmt vet test

EXECUTABLE ?= hyperpush
PACKAGES = $(shell go list ./... | grep -v /vendor/)
VERSION ?= $(shell git describe --match 'v[0-9]*' --dirty='-dev' --always)
COMMIT ?= $(shell git rev-parse --short HEAD)
LDFLAGS = -X "main.Revision=$(COMMIT)" -X "main.Version=$(VERSION)"

all: deps build test

clean:
	@go clean -i ./...

deps:
	@glide install

fmt:
	@go fmt $(PACKAGES)

vet:
	@go vet $(PACKAGES)

test:
	@for PKG in $(PACKAGES); do go test -cover -coverprofile $$GOPATH/src/$$PKG/coverage.out $$PKG || exit 1; done;

cover: test
	@echo ""
	@for PKG in $(PACKAGES); do go tool cover -func $$GOPATH/src/$$PKG/coverage.out; echo ""; done;


$(EXECUTABLE): $(shell find . -type f -print | grep -v vendor | grep "\.go")
	@echo "Building $(EXECUTABLE)..."
	@go build -ldflags '-s -w $(LDFLAGS)' -o $(EXECUTABLE) cmd/main.go

build: $(EXECUTABLE)

run: build
	@./$(EXECUTABLE)
