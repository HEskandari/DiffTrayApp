package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFinder_FindsProjectRoot(t *testing.T) {
	finder := newProjectFinder()
	path := finder.Find("./_testdata/dir1/TestFile.txt")

	assert.NotEmpty(t, path)
	assert.Contains(t, path, "go.mod")
	assert.Len(t, finder.cache, 1)
}