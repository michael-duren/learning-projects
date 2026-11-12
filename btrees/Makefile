BINARY := btrees
BIN_DIR := bin
PREFIX ?= $(HOME)/.local

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BINARY) .

run:
	go run . $(ARGS)

test:
	go test .

