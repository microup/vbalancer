package core

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ArchiveFile(fileName string, extension string) {
	file := filepath.Base(fileName)                                 // + ".zip"
	file = strings.TrimSuffix(file, filepath.Ext(file)) + extension // ".zip"
	path := filepath.Dir(fileName)
	fs := filepath.Join(path, file)

	archive, err := os.Create(fs)
	defer func(archive *os.File) {
		err = archive.Close()
		if err != nil {
			log.Fatalf("archive close failed: %v", err)
		}
	}(archive)

	if err != nil {
		panic(err)
	}

	csvFileName, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	defer func(f1 *os.File) {
		err = f1.Close()
		if err != nil {
			log.Fatalf("file close failed: %v", err)
		}
	}(csvFileName)

	zipWriter := zip.NewWriter(archive)
	defer func(zipWriter *zip.Writer) {
		err = zipWriter.Close()
		if err != nil {
			log.Fatalf("zipWriter close failed: %v", err)
		}
	}(zipWriter)

	fc := filepath.Base(fileName)
	fileZip, err := zipWriter.Create(fc)

	if err != nil {
		panic(err)
	}

	if _, err := io.Copy(fileZip, csvFileName); err != nil {
		panic(err)
	}
}
