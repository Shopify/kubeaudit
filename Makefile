# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
BINARY_NAME=kubeaudit
BINARY_UNIX=$(BINARY_NAME)_unix
LDFLAGS=$(shell build/ldflags.sh)

# kubernetes client won't build with go<1.10
GOVERSION:=$(shell go version | awk '{print $$3}')
GOVERSION_MIN:=go1.17
GOVERSION_CHECK=$(shell printf "%s\n%s\n" "$(GOVERSION)" "$(GOVERSION_MIN)" | sort -t. -k 1,1n -k 2,2n -k 3,3n -k 4,4n | head -n 1)

# Test parameters
TEST_CLUSTER_NAME="kubeaudit-test"

export GO111MODULE=on

ifneq ($(GOVERSION_MIN), $(GOVERSION_CHECK))
$(error Detected Go version $(GOVERSION) < required version $(GOVERSION_MIN))
endif

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v -ldflags=all="$(LDFLAGS)" cmd/main.go

install:
	cp $(BINARY_NAME) $(GOPATH)/bin/kubeaudit

plugin:
	cp $(BINARY_NAME) $(GOPATH)/bin/kubectl-audit

test:
	./test.sh

test-setup:
	kind create cluster --name ${TEST_CLUSTER_NAME} --image kindest/node:v1.20.15@sha256:6f2d011dffe182bad80b85f6c00e8ca9d86b5b8922cdf433d53575c4c5212248

test-teardown:
	kind delete cluster --name ${TEST_CLUSTER_NAME}

show-coverage: test
	go tool cover -html=coverage.txt

setup:
	$(GOMOD) download
	$(GOMOD) tidy

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

# Cross Compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

docker-build:
	docker run --rm -it -v "$(GOPATH)":/go -w /go/src/github.com/Shopify/kubeaudit golang:1.12 go build -o "$(BINARY_UNIX)" -v

.PHONY: all build install plugin test test-setup test-teardown show-coverage clean build-linux docker-build
