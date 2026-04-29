BINARY := opencode-skeleton
BUILD_DIR := build

.PHONY: build install clean cross test

build:
	go build -o $(BUILD_DIR)/$(BINARY) ./cmd/skeleton

install: build
	cp $(BUILD_DIR)/$(BINARY) $(HOME)/.local/bin/$(BINARY)

clean:
	rm -rf $(BUILD_DIR)

cross:
	GOOS=linux   GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY)-linux-amd64   ./cmd/skeleton
	GOOS=darwin  GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY)-darwin-amd64  ./cmd/skeleton
	GOOS=darwin  GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY)-darwin-arm64  ./cmd/skeleton
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY)-windows-amd64.exe ./cmd/skeleton

test:
	go test ./...
