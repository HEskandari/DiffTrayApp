package main

import (
	"github.com/VerifyTests/Verify.Go/verifier"
	"github.com/google/uuid"
	"testing"
)

func TestMain_AddsItemsDynamically(t *testing.T) {
	verifier.NewVerifier(t).Configure(
		verifier.UseDirectory("./_testdata"),
		verifier.TestCase("TestingUI"),
	).Verify(uuid.NewString())
}
