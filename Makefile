SOURCEDIR="."
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=ffprogress
BINARY_RELEASE=bin/${BINARY}_${VERSION}

VERSION=$(shell cat VERSION)

.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(SOURCES)
	go build -o bin/${BINARY}

static_linux: bin_dir
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ${BINARY_RELEASE}_linux_amd64

static_darwin: bin_dir
	CGO_ENABLED=0 GOOS=darwin go build -a -installsuffix cgo -o ${BINARY_RELEASE}_darwin_amd64

.PHONY: bin_dir
bin_dir:
	mkdir -p bin

.PHONY: run
run:
	go run main.go

.PHONY: clean
clean:
	rm -f bin/${BINARY} bin/${BINARY}_*
