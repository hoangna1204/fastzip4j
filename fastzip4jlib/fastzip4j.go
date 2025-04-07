package main

/*
#include <stdlib.h>
*/

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/klauspost/compress/zip"
	"github.com/saracen/fastzip"
)

import "C"

const temporaryRootPath = ".fastzip4j/"
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func main() {}

//export ArchiveFile
func ArchiveFile(SourceFile *C.char, ZipDestination *C.char, TemporaryPath *C.char, CompressionLevel int) {
	srcFile := C.GoString(SourceFile)
	zipDest := C.GoString(ZipDestination)
	tempDir := C.GoString(TemporaryPath)

	if srcFile == "" {
		panic("Source File cannot be empty")
	}
	if zipDest == "" {
		panic("Destination File cannot be empty")
	}
	if CompressionLevel < 1 || CompressionLevel > 9 {
		panic("CompressionLevel must be between 1 and 9! Following zlib, levels range from 1 (BestSpeed) to 9 (BestCompression); higher levels typically run slower but compress more.")
	}
	if tempDir == "" {
		panic("Temporary Path cannot be empty")
	}

	_, err := os.Stat(zipDest)
	zipFileExists := !os.IsNotExist(err)

	if zipFileExists {
		_, err = os.Stat(tempDir)
		if os.IsNotExist(err) {
			err := os.MkdirAll(tempDir, os.ModePerm)
			if err != nil {
				panic(err)
			}
		}

		e, err := fastzip.NewExtractor(zipDest, tempDir)
		if err != nil {
			panic(err)
		}
		defer func(e *fastzip.Extractor) {
			err := e.Close()
			if err != nil {
				panic(err)
			}
		}(e)
		if err = e.Extract(context.Background()); err != nil {
			panic(err)
		}
	}

	err = copyFile(srcFile, tempDir)
	if err != nil {
		panic(err)
	}

	w, err := os.Create(zipDest)
	if err != nil {
		panic(err)
	}
	defer func(w *os.File) {
		err := w.Close()
		if err != nil {
			panic(err)
		}
	}(w)

	archiver, err := fastzip.NewArchiver(w, tempDir)
	if err != nil {
		panic(err)
	}
	defer func(archiver *fastzip.Archiver) {
		err := archiver.Close()
		if err != nil {
			panic(err)
		}
	}(archiver)

	archiver.RegisterCompressor(zip.Deflate, fastzip.FlateCompressor(CompressionLevel))

	files := make(map[string]os.FileInfo)
	err = filepath.Walk(tempDir, func(pathname string, info os.FileInfo, err error) error {
		files[pathname] = info
		return nil
	})

	if err = archiver.Archive(context.Background(), files); err != nil {
		panic(err)
	}

	if zipFileExists {
		removeFs(tempDir)
	}
}

//export ArchiveDir
func ArchiveDir(SourceDir *C.char, ZipFile *C.char, TemporaryPath *C.char, CompressionLevel int) {
	srcFile := C.GoString(SourceDir)
	zipDest := C.GoString(ZipFile)
	tempDir := C.GoString(TemporaryPath)

	if srcFile == "" {
		panic("Source Folder cannot be empty!")
	}
	if zipDest == "" {
		panic("ZipFile cannot be empty!")
	}
	if CompressionLevel < 1 || CompressionLevel > 9 {
		panic("CompressionLevel must be between 1 and 9! Following zlib, levels range from 1 (BestSpeed) to 9 (BestCompression); higher levels typically run slower but compress more.")
	}
	if tempDir == "" {
		panic("TemporaryPath cannot be empty!")
	}

	_, err := os.Stat(zipDest)
	zipFileExists := !os.IsNotExist(err)
	cleanDsStoreFile(srcFile)

	if zipFileExists {
		_, err := os.Stat(tempDir)
		if os.IsNotExist(err) {
			err := os.MkdirAll(tempDir, os.ModePerm)
			cleanDsStoreFile(tempDir)
			if err != nil {
				panic(err)
			}
		}

		e, err := fastzip.NewExtractor(zipDest, tempDir)
		if err != nil {
			panic(err)
		}
		defer func(e *fastzip.Extractor) {
			err := e.Close()
			if err != nil {
				panic(err)
			}
		}(e)
		if err = e.Extract(context.Background()); err != nil {
			panic(err)
		}

		err = os.CopyFS(tempDir, os.DirFS(srcFile))
		if err != nil {
			panic(err)
		}
	}
	w, err := os.Create(zipDest)
	if err != nil {
		panic(err)
	}
	defer func(w *os.File) {
		err := w.Close()
		if err != nil {
			panic(err)
		}
	}(w)

	sourcePath := srcFile
	if zipFileExists {
		sourcePath = tempDir
	}
	archiver, err := fastzip.NewArchiver(w, sourcePath)
	if err != nil {
		panic(err)
	}
	defer func(archiver *fastzip.Archiver) {
		err := archiver.Close()
		if err != nil {
			panic(err)
		}
	}(archiver)

	archiver.RegisterCompressor(zip.Deflate, fastzip.FlateCompressor(CompressionLevel))

	files := make(map[string]os.FileInfo)
	err = filepath.Walk(sourcePath, func(pathname string, info os.FileInfo, err error) error {
		files[pathname] = info
		return nil
	})

	if err = archiver.Archive(context.Background(), files); err != nil {
		panic(err)
	}

	if zipFileExists {
		removeFs(tempDir)
	}
}

//export Extract
func Extract(ZipFile *C.char, DestinationDirectory *C.char) {
	zip := C.GoString(ZipFile)
	dest := C.GoString(DestinationDirectory)

	if zip == "" {
		panic("ZipFile cannot be empty!")
	}

	if dest == "" {
		panic("DestinationDirectory cannot be empty!")
	}

	e, err := fastzip.NewExtractor(zip, dest)
	if err != nil {
		panic(err)
	}
	defer func(e *fastzip.Extractor) {
		err := e.Close()
		if err != nil {
			panic(err)
		}
	}(e)

	if err = e.Extract(context.Background()); err != nil {
		panic(err)
	}
}

func copyFile(sourcePath string, destDirectory string) error {
	// Open the source file
	srcFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer func(srcFile *os.File) {
		err := srcFile.Close()
		if err != nil {
			panic(err)
		}
	}(srcFile)

	// Ensure the destination directory exists
	err = os.MkdirAll(destDirectory, os.ModePerm)
	if err != nil {
		return err
	}

	dstPath := filepath.Join(destDirectory, filepath.Base(sourcePath))

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer func(dstFile *os.File) {
		err := dstFile.Close()
		if err != nil {
			panic(err)
		}
	}(dstFile)

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

func removeFs(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		panic(err)
	}
}

func cleanDsStoreFile(dirPath string) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if file.Name() == ".DS_Store" {
			err = os.Remove(filepath.Join(dirPath, file.Name()))
			if err != nil {
				panic(err)
			}
		}
	}
}
