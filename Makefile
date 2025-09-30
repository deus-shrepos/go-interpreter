.SILENT:

GREEN := \033[32m
RED := \033[31m
RESET := \033[0m
ORANGE := \033[33m

# TODO: NEED TO ADD MORE STUFF
build:
	echo "Building the interpreter"
	go build -o bin/interpreter main.go
	echo "Build complete"

run:
	echo "Running the interpreter"
	go run cmd/repl/main.go $(ARGS)
