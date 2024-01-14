EXECUTABLE=marino
WINDOWS=build/$(EXECUTABLE)_amd64.exe
LINUX=build/$(EXECUTABLE)_linux
DARWIN=build/$(EXECUTABLE)
VERSION=$(shell cat VERSION)

.PHONY: all test clean web

all: test build ## Build and run tests

test: ## Run unit tests
	go test ./...

web: ## Run unit tests
	cd web && env REACT_APP_VERSION=$(VERSION) npm run build-local

build: web windows linux darwin ## Build binaries
	@echo version: $(VERSION)

windows: $(WINDOWS) ## Build for Windows

linux: $(LINUX) ## Build for Linux

darwin: $(DARWIN) ## Build for Darwin (macOS)

$(WINDOWS):
	go generate
	env GOOS=windows GOARCH=amd64 go build -v -o $(WINDOWS) -ldflags="-s -w -X dashboard/version.Version=$(VERSION)" main.go

$(LINUX):
	env GOOS=linux GOARCH=amd64 go build -v -o $(LINUX) -ldflags="-s -w -X dashboard/version.Version=$(VERSION)"  main.go

$(DARWIN):
	env GOOS=darwin GOARCH=amd64 go build -v -o $(DARWIN)_amd64 -ldflags="-s -w -X dashboard/version.Version=$(VERSION)" main.go
	env GOOS=darwin GOARCH=arm64 go build -v -o $(DARWIN)_arm64 -ldflags="-s -w -X dashboard/version.Version=$(VERSION)" main.go
	lipo -create -output $(DARWIN) $(DARWIN)_amd64 $(DARWIN)_arm64
	rm -f $(DARWIN)_amd64 $(DARWIN)_arm64

clean: ## Remove previous build
	rm -f $(WINDOWS) $(LINUX) $(DARWIN)

help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
