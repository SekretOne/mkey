OUT_DIR := ./out

ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))
export GOBIN := $(ROOT_DIR)/gobin

.PHONY: clean
clean:
	rm -rf $(OUT_DIR)

dep:
	go mod tidy
	go mod download
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.2

lint:
	./gobin/golangci-lint run
