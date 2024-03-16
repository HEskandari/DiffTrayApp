package main

import (
	"github.com/VerifyTests/Verify.Go/utils"
	"io"
	"log"
	"os"
)

func safeMoveFile(sourcePath, destPath string) bool {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		log.Printf("couldn't open source file: %s", err)
		return false
	}

	outputFile, err := os.Create(destPath)
	if err != nil {
		safeCloseFile(inputFile)
		log.Printf("couldn't open dest file: %s", err)
		return false
	}

	defer safeCloseFile(outputFile)
	_, err = io.Copy(outputFile, inputFile)
	safeCloseFile(inputFile)
	if err != nil {
		log.Printf("writing to output file failed: %s", err)
		return false
	}

	err = os.Remove(sourcePath)
	if err != nil {
		log.Printf("failed removing original file: %s", err)
		return false
	}

	return true
}

func safeDeleteDirectory(path string) {
	if exists, _ := utils.File.FileOrDirectoryExists(path); !exists {
		return
	}

	if isEmptyDirectory(path) {
		if !safeDelete(path) {
			log.Printf("Failed to delete '%s'", path)
		}
	}
}

func safeCloseFile(file *os.File) {
	_ = file.Close()
}

func isEmptyDirectory(path string) bool {
	content, err := utils.File.GetDirectoriesFromRoot(path)
	if err != nil {
		return false
	}

	return len(content) == 0
}

func safeDelete(path string) bool {
	if utils.File.Exists(path) {
		err := os.Remove(path)
		if err != nil {
			return false
		}
	}
	return true
}
