# Build for Linux
env GOOS=linux GOARCH=amd64 go build -buildmode=c-shared -o ../fastzip4j/src/main/resources/lib/libfastzip4j.so fastzip4j.go

# Build for Windows
env GOOS=windows GOARCH=amd64 go build -buildmode=c-shared -o ../fastzip4j/src/main/resources/lib/libfastzip4j.dll fastzip4j.go

# Build for macOS
env GOOS=darwin GOARCH=amd64 go build -buildmode=c-shared -o ../fastzip4j/src/main/resources/lib/libfastzip4j.dylib fastzip4j.go