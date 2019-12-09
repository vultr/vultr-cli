.PHONY: \
	build \
	remove \
	format \
	build_mac \
	build_linux_386 \
	build_linux_amd64 \
	build_windows_64 \
	build_windows_32
export CGO=0

DIR != echo $${DIR:-$$(pwd)/builds}

# Builds for mac os
GOOS = darwin
GOARCH = amd64
build_mac: _build

# Builds for linux 32bit
GOOS = linux
GOARCH = 386
build_linux_386: _build

# Builds for linux 64bit
GOOS = linux
GOARCH = amd64
build_linux_amd64: _build

# Builds for Windows 64bit
GOOS = windows
GOARCH = amd64
build_windows_64: _build_win

# Builds for Windows 32bit
GOOS = windows
GOARCH = 386
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

format:
	@go fmt ./...
