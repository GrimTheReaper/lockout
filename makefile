# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=lockout
BINARY_UNIX=$(BINARY_NAME)_unix

all: test build package

test:
	$(GOTEST) ./...

build:
	env GOOS=linux GARCH=amd64 $(GOBUILD) -o $(BINARY_NAME) -ldflags "-linkmode external -extldflags -static" ./cmd/main.go

package:
	# Gonna copy this cert, cause I can't get the mmdb otherwise...
	cp /etc/ssl/certs/ca-certificates.crt .
	docker build -t awildtyphlosion/lockout .
	# Lets not let it linger.
	rm ca-certificates.crt

# Shouldn't need to redo the protofile, but might as well support it.
proto:
	protoc -I . --go_out=plugins=grpc:. pb/*.proto
