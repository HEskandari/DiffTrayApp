package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFinder_FindsProjectRoot(t *testing.T) {
	finder := newSolutionFinder()
	path := finder.Find("./_testdata/dir1/TestFile.txt")

	assert.NotEmpty(t, path)
	assert.Contains(t, path, "go.mod")
}