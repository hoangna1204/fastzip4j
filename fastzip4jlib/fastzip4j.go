package main

/*
#include <stdlib.h>
*/

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/klauspost/compress/zip"
	"github.com/saracen/fastzip"
)

import "C"

const temporaryRootPath = ".fastzip4j/"

func main() {}

//export ArchiveFile
func ArchiveFile(SourceFile *C.char, ZipDestination *C.char, CompressionLevel int) {
	srcFile := C.GoString(SourceFile)
	zipDest := C.GoString(ZipDestination)

	if srcFile == "" {
		panic("Source File cannot be empty")
	}
	if zipDest == "" {
		panic("Destination File cannot be empty")
	}
	if CompressionLevel < 1 || CompressionLevel > 9 {
		panic("CompressionLevel must be between 1 and 9! Following zlib, levels range from 1 (BestSpeed) to 9 (BestCompression); higher levels typically run slower but compress more.")
	}
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	temporaryPath := temporaryRootPath + timestamp + "/"

	_, err := os.Stat(zipDest)
	zipFileExists := !os.IsNotExist(err)

	if zipFileExists {
		_, err = os.Stat(temporaryPath)
		if os.IsNotExist(err) {
			err := os.MkdirAll(temporaryPath, os.ModePerm)
			if err != nil {
				panic(err)
			}
		}

		e, err := fastzip.NewExtractor(zipDest, temporaryPath)
		if err != nil {
			os.RemoveAll(temporaryPath)
			panic(err)
		}
		defer func(e *fastzip.Extractor) {
			err := e.Close()
			if err != nil {
				os.RemoveAll(temporaryPath)
				panic(err)
			}
		}(e)
		if err = e.Extract(context.Background()); err != nil {
			os.RemoveAll(temporaryPath)
			panic(err)
		}
	}

	err = copyFile(srcFile, temporaryPath)
	if err != nil {
		os.RemoveAll(temporaryPath)
		panic(err)
	}

	w, err := os.Create(zipDest)
	if err != nil {
		os.RemoveAll(temporaryPath)
		panic(err)
	}
	defer func(w *os.File) {
		err := w.Close()
		if err != nil {
			os.RemoveAll(temporaryPath)
			panic(err)
		}
	}(w)

	archiver, err := fastzip.NewArchiver(w, temporaryPath)
	if err != nil {
		os.RemoveAll(temporaryPath)
		panic(err)
	}
	defer func(archiver *fastzip.Archiver) {
		err := archiver.Close()
		if err != nil {
			os.RemoveAll(temporaryPath)
			panic(err)
		}
	}(archiver)

	archiver.RegisterCompressor(zip.Deflate, fastzip.FlateCompressor(CompressionLevel))

	files := make(map[string]os.FileInfo)
	err = filepath.Walk(temporaryPath, func(pathname string, info os.FileInfo, err error) error {
		files[pathname] = info
		return nil
	})

	if err = archiver.Archive(context.Background(), files); err != nil {
		os.RemoveAll(temporaryPath)
		panic(err)
	}

	if zipFileExists {
		removeFs(temporaryPath)
	}
}

//export ArchiveDir
func ArchiveDir(SourceDir *C.char, ZipFile *C.char, CompressionLevel int) {
	srcFile := C.GoString(SourceDir)
	zipDest := C.GoString(ZipFile)

	if srcFile == "" {
		panic("Source Folder cannot be empty!")
	}
	if zipDest == "" {
		panic("ZipFile cannot be empty!")
	}
	if CompressionLevel < 1 || CompressionLevel > 9 {
		panic("CompressionLevel must be between 1 and 9! Following zlib, levels range from 1 (BestSpeed) to 9 (BestCompression); higher levels typically run slower but compress more.")
	}
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	temporaryPath := temporaryRootPath + timestamp + "/"

	_, err := os.Stat(zipDest)
	zipFileExists := !os.IsNotExist(err)
	cleanDsStoreFile(srcFile)

	if zipFileExists {
		_, err := os.Stat(temporaryPath)
		if os.IsNotExist(err) {
			err := os.MkdirAll(temporaryPath, os.ModePerm)
			cleanDsStoreFile(temporaryPath)
			if err != nil {
				panic(err)
			}
		}

		e, err := fastzip.NewExtractor(zipDest, temporaryPath)
		if err != nil {
			os.RemoveAll(temporaryPath)
			panic(err)
		}
		defer func(e *fastzip.Extractor) {
			err := e.Close()
			if err != nil {
				os.RemoveAll(temporaryPath)
				panic(err)
			}
		}(e)
		if err = e.Extract(context.Background()); err != nil {
			os.RemoveAll(temporaryPath)
			panic(err)
		}

		err = os.CopyFS(temporaryPath, os.DirFS(srcFile))
		if err != nil {
			os.RemoveAll(temporaryPath)
			panic(err)
		}
	}
	w, err := os.Create(zipDest)
	if err != nil {
		os.RemoveAll(temporaryPath)
		panic(err)
	}
	defer func(w *os.File) {
		err := w.Close()
		if err != nil {
			os.RemoveAll(temporaryPath)
			panic(err)
		}
	}(w)

	sourcePath := srcFile
	if zipFileExists {
		sourcePath = temporaryPath
	}
	archiver, err := fastzip.NewArchiver(w, sourcePath)
	if err != nil {
		os.RemoveAll(temporaryPath)
		panic(err)
	}
	defer func(archiver *fastzip.Archiver) {
		err := archiver.Close()
		if err != nil {
			os.RemoveAll(temporaryPath)
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
		os.RemoveAll(temporaryPath)
		panic(err)
	}

	if zipFileExists {
		removeFs(temporaryPath)
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
