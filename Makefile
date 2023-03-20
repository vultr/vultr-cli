.PHONY: remove format

export CGO=0
export GOFLAGS=-trimpath

DIR?=builds

$(DIR):
	mkdir -p $(DIR)

$(DIR)/vultr-cli_darwin_amd64: $(DIR)
	env GOOS=darwin GOARCH=amd64 go build -o $@

$(DIR)/vultr-cli_darwin_arm64: $(DIR)
	env GOOS=darwin GOARCH=arm64 go build -o $@

$(DIR)/vultr-cli_linux_386: $(DIR)
	env GOOS=linux GOARCH=386 go build -o $@

$(DIR)/vultr-cli_linux_amd64: $(DIR)
	env GOOS=linux GOARCH=amd64 go build -o $@

$(DIR)/vultr-cli_linux_arm64: $(DIR)
	env GOOS=linux GOARCH=arm64 go build -o $@

$(DIR)/vultr-cli_windows_386.exe: $(DIR)
	env GOOS=windows GOARCH=386 go build -o $@

$(DIR)/vultr-cli_windows_amd64.exe: $(DIR)
	env GOOS=windows GOARCH=amd64 go build -o $@

$(DIR)/vultr-cli_linux_arm: $(DIR)
	env GOOS=linux GOARCH=arm go build -o $@

remove:
	@rm -rf builds

format:
	@go fmt ./...

docker:
	docker build . -t vultr/vultr-cli
	docker push vultr/vultr-cli
