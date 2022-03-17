package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type finderResult struct {
	Directory string
	Name      string
}

type solutionFinder struct {
	cache map[string]finderResult
}

func newProjectFinder() *solutionFinder {
	return &solutionFinder{
		cache: make(map[string]finderResult),
	}
}

func (s *solutionFinder) Find(filePath string) string {
	if val, ok := s.cache[filePath]; ok {
		return val.Name
	}

	result, found := s.find(filePath)
	if found {
		s.cache[filePath] = result
	}
	return result.Name
}

func (s *solutionFinder) find(filePath string) (finderResult, bool) {
	for _, val := range s.cache {
		if strings.HasPrefix(filePath, val.Directory) {
			return val, true
		}
	}

	currentDirectory := getDirectoryName(filePath)
	if len(currentDirectory) == 0 {
		return finderResult{}, false
	}

	for {
		projects, err := getDirectoriesFromRoot(currentDirectory, "go.mod")
		if err != nil {
			panic(fmt.Sprintf("failed to search for directory: %s", currentDirectory))
		}

		if len(projects) > 0 {
			return finderResult{
				Directory: currentDirectory,
				Name:      projects[0],
			}, true
		}

		parent := getParentDirectory(currentDirectory)
		if len(parent) == 0 {
			return finderResult{}, false
		}

		currentDirectory = parent
	}
}

func getDirectoryName(filePath string) string {
	dir := filepath.Dir(filePath)
	abs, err := filepath.Abs(dir)

	if err != nil {
		panic("failed to get the full file path")
	}
	return abs
}

func getParentDirectory(filePath string) string {
	dir := filepath.Dir(filePath)
	//base := filepath.Base(dir)
	abs, err := filepath.Abs(dir)

	if err != nil {
		panic("failed to get the full path")
	}
	return abs
}

func GetDirectoriesFromRoot(root string) ([]string, error) {
	matches := make([]string, 0)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func getDirectoriesFromRoot(root, segment string) ([]string, error) {
	matches := make([]string, 0)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if matched, _ := filepath.Match(segment, info.Name()); matched {
				matches = append(matches, path)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return matches, nil
}