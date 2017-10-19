# Go parameters
GOCMD=go
GODEP=dep
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=kubeaudit
BINARY_UNIX=$(BINARY_NAME)_unix

all: setup test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

test:
	$(GOTEST) -cover ./... .

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
	docker run --rm -it -v "$(GOPATH)":/go -w /go/src/github.com/Shopify/kubeaudit golang:1.9 go build -o "$(BINARY_UNIX)" -v

.PHONY: build clean test check_version build-linux docker-build
