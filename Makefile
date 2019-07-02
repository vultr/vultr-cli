.PHONY: build
export CGO=0

DIR = $(shell echo $${DIR:-$(shell pwd)/builds})

# Builds for mac os
build_mac: GOOS = darwin
build_mac: GOARCH = amd64
build_mac: _build

# Builds for linux 32bit
build_linux_386: GOOS = linux
build_linux_386: GOARCH = 386
build_linux_386: _build

# Builds for linux 64bit
build_linux_amd64: GOOS = linux
build_linux_amd64: GOARCH = amd64
build_linux_amd64: _build

# Builds for Windows 64bit
build_windows_64: GOOS = windows
build_windows_64: GOARCH = amd64
build_windows_64: _build_win

# Builds for Windows 32bit
build_windows_32: GOOS = windows
build_windows_32: GOARCH = 386
build_windows_32: _build_win

_build:
	@mkdir -p builds
	@echo "building vultr-cli"
	env GOOS=$(GOOS) GOARCH=$(GOARCH) go build -gcflags=-trimpath=$$GOPATH -o $(DIR)/vultr-cli_$(GOOS)_$(GOARCH)
	@echo "built $(DIR)/vultr-cli_$(GOOS)_$(GOARCH)"

_build_win:
	@mkdir -p builds
	@echo "building vultr-cli"
	env GOOS=$(GOOS) GOARCH=$(GOARCH) go build -gcflags=-trimpath=$$GOPATH -o $(DIR)/vultr-cli_$(GOOS)_$(GOARCH).exe
	@echo "built $(DIR)/vultr-cli_$(GOOS)_$(GOARCH)"

remove:
	@rm -rf builds
