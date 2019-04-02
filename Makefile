# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=v2web
BINARY_VERSION=0.0.1
BINARY_UNIX=$(BINARY_NAME)_$(BINARY_VERSION).unix
BINARY_DARWIN=$(BINARY_NAME)_$(BINARY_VERSION).darwin
BINARY_WINDOWS=$(BINARY_NAME)_$(BINARY_VERSION).windows

all: test build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v main/web.go
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME) $(BINARY_UNIX)_* $(BINARY_DARWIN)_* $(BINARY_WINDOWS)_*
run: build
	./$(BINARY_NAME)
deps:
	$(GOGET) github.com/go-redis/redis
	$(GOGET) gopkg.in/yaml.v2

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX)_amd64 -v
	CGO_ENABLED=0 GOOS=linux GOARCH=386 $(GOBUILD) -o $(BINARY_UNIX)_386 -v
build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_DARWIN)_amd64 -v
	CGO_ENABLED=0 GOOS=darwin GOARCH=386 $(GOBUILD) -o $(BINARY_DARWIN)_386 -v
build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_WINDOWS)_amd64.exe -v
	CGO_ENABLED=0 GOOS=windows GOARCH=386 $(GOBUILD) -o $(BINARY_WINDOWS)_386.exe -v
build-all: build-linux build-darwin build-windows
docker-build:
	docker run --rm -it -v "$(GOPATH)":/go -w /go/src/bitbucket.org/rsohlich/makepost golang:latest go build -o "$(BINARY_UNIX)" -v