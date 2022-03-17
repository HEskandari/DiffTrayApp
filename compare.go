package main

import (
	"bytes"
	"io"
	"log"
	"os"
)

const chunkSize = 64000

type fileComparer struct {
}

func newFileComparer() *fileComparer {
	return &fileComparer{}
}

func (f *fileComparer) FilesAreEqual(file1, file2 string) bool {
	if f.FilesAreSameSize(file1, file2) {
		return false
	}

	return f.StreamsAreEqual(file1, file2)
}

func (f *fileComparer) StreamsAreEqual(file1, file2 string) bool {
	f1, err := os.Open(file1)
	if err != nil {
		log.Fatal(err)
	}
	defer f1.Close()

	f2, err := os.Open(file2)
	if err != nil {
		log.Fatal(err)
	}
	defer f2.Close()

	for {
		b1 := make([]byte, chunkSize)
		_, err1 := f1.Read(b1)

		b2 := make([]byte, chunkSize)
		_, err2 := f2.Read(b2)

		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return true
			} else if err1 == io.EOF || err2 == io.EOF {
				return false
			} else {
				log.Fatal(err1, err2)
			}
		}

		if !bytes.Equal(b1, b2) {
			return false
		}
	}
}

func (f *fileComparer) FilesAreSameSize(file1, file2 string) bool {
	s1, _ := getFileSize(file1)
	s2, _ := getFileSize(file2)

	return s1 == s2
}

func getFileSize(filePath string) (int64, error) {
	fi, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}

	size := fi.Size()
	return size, nil
}