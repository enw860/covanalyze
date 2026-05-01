.PHONY: build test clean install

# Build the covanalyze binary
build:
	@mkdir -p bin
	go build -o bin/covanalyze ./cmd

# Run all unit tests with coverage
test:
	go test -v -cover ./...

# Remove bin directory and build artifacts
clean:
	rm -rf bin/

# Install target (builds to bin/covanalyze, no system-wide install)
install: build