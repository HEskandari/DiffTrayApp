package main

import (
	"github.com/VerifyTests/Verify.Go/verifier"
	"github.com/google/uuid"
	"os"
	"strconv"
	"testing"
)

func TestMain_AddsItemsDynamically(t *testing.T) {
	runTests, _ := strconv.ParseBool(os.Getenv("RUN_INTEGRATION_TESTS"))
	if !runTests {
		t.Skip("Skipping integration tests")
	}

	verifier.NewVerifier(t).Configure(
		verifier.UseDirectory("./_testdata"),
		verifier.TestCase("TestingUI"),
	).Verify(uuid.NewString())
}
