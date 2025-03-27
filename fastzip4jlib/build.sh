go build -buildmode=c-shared -o ../fastzip4j/src/main/resources/lib/libfastzip4j.so fastzip4j.go
go build -buildmode=c-shared -o ../fastzip4j/src/main/resources/lib/libfastzip4j.dll fastzip4j.go
go build -buildmode=c-shared -o ../fastzip4j/src/main/resources/lib/libfastzip4j.dylib fastzip4j.go