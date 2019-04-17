BUILD_DIR ?= build
PLUGIN_DIR ?= $(BUILD_DIR)/plugin
COMMIT = $(shell git rev-parse HEAD)
VERSION ?= $(shell git describe --always --tags --dirty)
ORG := github.com/hyperscale
PROJECT := hyperpush
REPOPATH ?= $(ORG)/$(PROJECT)
VERSION_PACKAGE = $(REPOPATH)/pkg/hyperpush/version

GO_LDFLAGS :="
GO_LDFLAGS += -X $(VERSION_PACKAGE).version=$(VERSION)
GO_LDFLAGS += -X $(VERSION_PACKAGE).buildDate=$(shell date +'%Y-%m-%dT%H:%M:%SZ')
GO_LDFLAGS += -X $(VERSION_PACKAGE).gitCommit=$(COMMIT)
GO_LDFLAGS += -X $(VERSION_PACKAGE).gitTreeState=$(if $(shell git status --porcelain),dirty,clean)
GO_LDFLAGS +="

GO_FILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: all
all: deps build test

.PHONY: deps
deps:
	@go mod download

.PHONY: clean
clean:
	@go clean -i ./...

$(BUILD_DIR)/coverage.out: $(GO_FILES)
	@CGO_ENABLED=0 go test -cover -coverprofile $(BUILD_DIR)/coverage.out.tmp ./...
	@cat $(BUILD_DIR)/coverage.out.tmp | grep -v '.pb.go' | grep -v 'mock_' > $(BUILD_DIR)/coverage.out
	@rm $(BUILD_DIR)/coverage.out.tmp

ci-test:
	@go test -race -cover -coverprofile ./coverage.out.tmp -v ./... | go2xunit -fail -output tests.xml
	@cat ./coverage.out.tmp | grep -v '.pb.go' | grep -v 'mock_' > ./coverage.out
	@rm ./coverage.out.tmp
	@echo ""
	@go tool cover -func ./coverage.out

.PHONY: lint
lint:
	@golangci-lint run ./...

.PHONY: test
test: $(BUILD_DIR)/coverage.out

.PHONY: coverage
coverage: $(BUILD_DIR)/coverage.out
	@echo ""
	@go tool cover -func ./$(BUILD_DIR)/coverage.out

.PHONY: coverage-html
coverage-html: $(BUILD_DIR)/coverage.out
	@go tool cover -html ./$(BUILD_DIR)/coverage.out

generate: $(GO_FILES)
	@go generate ./...

${BUILD_DIR}/hyperpush-server: $(GO_FILES)
	@echo "Building $@..."
	@go generate ./cmd/$(subst ${BUILD_DIR}/,,$@)/
	@CGO_ENABLED=0 go build -ldflags $(GO_LDFLAGS) -o $@ ./cmd/$(subst ${BUILD_DIR}/,,$@)/

.PHONY: run-hyperpush-server
run-hyperpush-server: ${BUILD_DIR}/hyperpush-server
	@echo "Running $<..."
	@./$<

run: run-hyperpush-server

.PHONY: build
build: ${BUILD_DIR}/hyperpush-server plugin


${PLUGIN_DIR}/authentication.jwt.so: $(shell find ./plugin/authentication/jwt -type f -name '*.go')
	@CGO_ENABLED=0 go build -buildmode=plugin -o $@ ./plugin/authentication/jwt

.PHONY: plugin
plugin: ${PLUGIN_DIR}/authentication.jwt.so
