package main

import (
	"github.com/VerifyTests/Verify.Go/verifier"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var deleteMsg = `
{
"Type":"Move",
"Temp":"theTempFilePath",
"Target":"theTargetFilePath",
"CanKill":true,
"Exe":"theExePath",
"Arguments":"TheArguments",
"ProcessId":1000
}`

var moveMsg = `{
"Type":"Move",
"CanKill":true,
"Exe":"theExe",
"Arguments":"TheArguments \"s\"",
"ProcessId":10
}`

var testSettings verifier.VerifySettings

func init() {
	testSettings = verifier.NewSettings()
	testSettings.UseDirectory("./_testdata")
}

func TestTracker_ReceivingMove(t *testing.T) {

	var receivedDelete *DeletePayload
	var receivedMove *MovePayload
	var tracker = newTracker()

	tracker.deleteHandler = func(cmd *DeletePayload) {
		receivedDelete = cmd
	}
	tracker.moveHandler = func(cmd *MovePayload) {
		receivedMove = cmd
	}

	receiveChan := make(chan string)

	go func() {
		receiveChan <- moveMsg
	}()

	go func() {
		tracker.startProcessor(receiveChan)
	}()

	shouldReceive := func() bool {
		return receivedMove != nil && receivedDelete == nil
	}

	assert.Eventually(t, shouldReceive, 12*time.Second, 3*time.Second)

	verifier.VerifyWithSetting(t, testSettings, receivedMove)
}