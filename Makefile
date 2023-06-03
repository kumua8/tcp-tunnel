BUILD_DIR := build
BINARY_NAME := tunnel

#environment
export GOPRIVATE=github.com/kukupa8
export GOPROXY=https://goproxy.io,direct

.PHONY: build

build-dirs:
	mkdir -p $(BUILD_DIR)

build:
	@make build-dirs
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./...
	cp -r conf build/

build-linux:
	@make build-dirs
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME) ./...
	cp -r conf build/

clean:
	rm -rf build