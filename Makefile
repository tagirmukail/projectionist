# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GO111MODULE=auto
CGO_ENABLED=0
BINARY_NAME=projectionist
BINARY_UNIX=$(BINARY_NAME)_unix
GOPATH=$(HOME)/go
PROTO_DIR=./proto
SWAGGER_DIR=./apps/swagger

build:
	@echo "#Build"
	@GO111MODULE=$(GO111MODULE) $(GOBUILD) -o $(BINARY_NAME) -v
	@echo "#Build completed"

test-verbose-cover:
	GO111MODULE=$(GO111MODULE) $(GOTEST) -v -cover ./...

test-cover:
	GO111MODULE=$(GO111MODULE) $(GOTEST) -cover ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

deps:
	@echo "#Download dependencies started"
	$(GOMOD) download
	@echo "#Download dependencies finished"

initmodules:
	$(GOMOD) init projectionist

tidy:
	GO111MODULE=$(GO111MODULE) $(GOMOD) tidy

protogen:
	protoc -I/usr/local/include \
	 -I$(GOPATH)/src \
	 -I$(PROTO_DIR) \
	 -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	  --grpc-gateway_out=logtostderr=true:$(PROTO_DIR) \
	 --swagger_out=logtostderr=true:$(SWAGGER_DIR) --go_out=plugins=grpc:$(PROTO_DIR) projectionist.proto