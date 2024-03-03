package main

import (
	"github.com/VerifyTests/Verify.Go/verifier"
	"github.com/google/uuid"
	"testing"
)

func TestMain_AddsItemsDynamically(t *testing.T) {
	//client := tray.NewClient()

	//for i := 0; i < 10; i++ {
	//	tmp := uuid.NewString()
	//	client.AddMove(fmt.Sprintf("./_testdata/dir1/%s.txt", tmp), "./_testdata/dir1/TestFile.txt", "Goland", nil, false, 0)
	//}

	verifier.NewVerifier(t).Configure(
		verifier.UseDirectory("./_testdata"),
		verifier.TestCase("TestingUI"),
	).Verify(uuid.NewString())
}
