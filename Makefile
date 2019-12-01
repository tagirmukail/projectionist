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
