#!/bin/bash

# Create output directory if it doesn't exist
mkdir -p ../fastzip4j/src/main/resources/lib

# Install Linux CGO dependencies if on macOS
if [[ "$OSTYPE" == "darwin"* ]]; then
    echo "Installing Linux CGO dependencies..."
    brew install FiloSottile/musl-cross/musl-cross
    export PATH="/opt/homebrew/opt/musl-cross/bin:$PATH"
fi

# Build for Linux
echo "Building for Linux..."
if [[ "$OSTYPE" == "darwin"* ]]; then
    env CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++ go build -buildmode=c-shared -o ../fastzip4j/src/main/resources/lib/libfastzip4j.so fastzip4j.go
else
    env CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -buildmode=c-shared -o ../fastzip4j/src/main/resources/lib/libfastzip4j.so fastzip4j.go
fi

# Build for Windows
echo "Building for Windows..."
env CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ go build -buildmode=c-shared -o ../fastzip4j/src/main/resources/lib/libfastzip4j.dll fastzip4j.go

# Build for macOS
echo "Building for macOS..."
env CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -buildmode=c-shared -o ../fastzip4j/src/main/resources/lib/libfastzip4j.dylib fastzip4j.go

# Ensure the libraries have proper permissions
chmod 755 ../fastzip4j/src/main/resources/lib/*.so
chmod 755 ../fastzip4j/src/main/resources/lib/*.dylib
chmod 755 ../fastzip4j/src/main/resources/lib/*.dll