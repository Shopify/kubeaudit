# Go parameters
GOCMD=go
GODEP=dep
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=kubeaudit
BINARY_UNIX=$(BINARY_NAME)_unix
LDFLAGS=$(shell build/ldflags.sh)

# kubernetes client won't build with go<1.10
GOVERSION:=$(shell go version | awk '{print $$3}')
GOVERSION_MIN:=go1.10
GOVERSION_CHECK=$(shell echo "$(GOVERSION)\n$(GOVERSION_MIN)" | sort -t. -k 1,1n -k 2,2n -k 3,3n -k 4,4n | head -n 1)

ifneq ($(GOVERSION_MIN), $(GOVERSION_CHECK))
$(error Detected Go version $(GOVERSION) < required version $(GOVERSION_MIN))
endif

all: setup test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v -ldflags=all="$(LDFLAGS)"

install:
	cp $(BINARY_NAME) $(GOPATH)/bin/kubeaudit

plugin:
	cp $(BINARY_NAME) $(GOPATH)/bin/kubectl-audit

test:
	./test.sh

show-coverage: test
	go tool cover -html=coverage.txt

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

check_version:
	$(GOBUILD) -o $(BINARY_NAME)
	./$(BINARY_NAME) version

setup:
	$(GOCMD) get -u github.com/golang/dep/cmd/dep
	$(GODEP) ensure

# Cross Compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

docker-build:
	docker run --rm -it -v "$(GOPATH)":/go -w /go/src/github.com/Shopify/kubeaudit golang:1.11 go build -o "$(BINARY_UNIX)" -v

.PHONY: build clean test check_version build-linux docker-build
