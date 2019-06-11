.PHONY: build
export CGO=0

DIR = $(shell echo $${DIR:-$(shell pwd)/builds})


# Default values
GOOS = linux
GOARCH = amd64

UNAME_S := $(shell uname -s)
UNAME_P := $(shell uname -p)

# If machine is a mac
ifeq ($(UNAME_S),Darwin)
  GOOS = darwin
  GOARCH = 386
endif

# If machine is 32bit
ifneq ($(UNAME_P), x86_64)
  GOARCH = 386
endif

# Builds for current machine
build: _build

# Builds for mac os
build_mac: GOOS = darwin
build_mac: GOARCH = 386
build_mac: _build

# Builds for linux 32bit
build_linux_386: GOARCH = 386
build_linux_386: _build

# Builds for linux 64bit
build_linux_amd64: GOARCH = amd64
build_linux_amd64: _build

_build:
	@mkdir -p builds
	@echo "building vultr-cli"
	env GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(DIR)/vultr-cli_$(GOOS)_$(GOARCH)
	@echo "built $(DIR)/vultr-cli_$(GOOS)_$(GOARCH)"

remove:
	@rm -rf builds
