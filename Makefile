.SILENT:

GREEN := \033[32m
RED := \033[31m
RESET := \033[0m
ORANGE := \033[33m

SHELL := /bin/sh


BINARY := gointer
CMD := ./cmd/gointer
INTER := ./cmd/gointer

.PHONY: init all build install
all: build install test clean

init:
	@echo "Using Shell: $(SHELL)"

build:
	mkdir -p bin
	go build -o bin/$(BINARY) $((CMD))
install:
	@echo "Building & Installing the binary..."
	mkdir -p $HOME/bin
	export GOBIN=$HOME/bin
	go install $(INTER)
	@echo "Install Done!"

update-env:
		@if echo "$$SHELL" | grep -q "zsh"; then \
		echo "Deteched ZSH"; \
		echo 'export PATH=$$HOME/bin:$$PATH' >> $$HOME/.zshrc; \
		zsh -ic "source $$HOME/.zshrc"; \
	elif echo "$$SHELL" | grep -q "bash"; then \
		echo "Deteched Bash";\
		echo 'export PATH=$$HOME/bin:$$PATH' >> ~/.bashrc; \
	fi; \
	echo "go interpreter has been added to your .rc file. Please source your .bashrc"

install-interpreter: install update-env

test:
	go test tests/
clean:
	rm -rf bin
