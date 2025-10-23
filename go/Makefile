.PHONY: build clean test bench install

# Build shared library for PHP FFI
build:
	go build -buildmode=c-shared -o libboxpacker.so .

# Build for different platforms
build-linux:
	GOOS=linux GOARCH=amd64 go build -buildmode=c-shared -o libboxpacker-linux-amd64.so .

build-macos:
	GOOS=darwin GOARCH=amd64 go build -buildmode=c-shared -o libboxpacker-darwin-amd64.dylib .
	GOOS=darwin GOARCH=arm64 go build -buildmode=c-shared -o libboxpacker-darwin-arm64.dylib .

build-windows:
	GOOS=windows GOARCH=amd64 go build -buildmode=c-shared -o libboxpacker-windows-amd64.dll .

# Run Go tests
test:
	go test -v .

# Run benchmarks
bench:
	go test -bench=. -benchmem .

# Clean build artifacts
clean:
	rm -f *.so *.dylib *.dll *.h
	go clean

# Install (copy to system library path - optional)
install: build
	@echo "Library built: libboxpacker.so"
	@echo "You can copy it to your PHP project or system library path"

# Quick build and test
all: clean build test

# Development: build with race detector (Linux only)
build-dev:
	go build -race -buildmode=c-shared -o libboxpacker-dev.so .
