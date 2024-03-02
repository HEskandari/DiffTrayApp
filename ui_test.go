package main

import (
	"github.com/VerifyTests/Verify.Go/tray"
	"testing"
)

func TestMain_AddsItemsDynamically(t *testing.T) {
	client := tray.NewClient()
	client.AddMove("./_testdata/dir1/TempFile.txt", "./_testdata/dir1/TestFile.txt", "Goland", nil, false, 0)
}
