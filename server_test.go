package main

import (
	"github.com/VerifyTests/Verify.Go/verifier"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var deleteMsg = `
{
"Type":"Delete",
"File":"Foo"
}`

var moveMsg = `{
"Type":"Move",
"CanKill":true,
"Exe":"theExe",
"Arguments":"TheArguments \"s\"",
"ProcessId":10
}`

func initVerifier(t *testing.T) verifier.Verifier {
	return verifier.NewVerifier(t).Configure(verifier.UseDirectory("./_testdata"))
}

func TestTracker_ReceivingMove(t *testing.T) {
	var ver = initVerifier(t)

	var receivedDelete *DeletePayload
	var receivedMove *MovePayload
	var tracker = newServer(func(cmd *MovePayload) {
		receivedMove = cmd
	}, func(cmd *DeletePayload) {
		receivedDelete = cmd
	}, func() {
		//No-Op
	})

	receiveChan := make(chan string)
	tracker.processor = receiveChan

	go func() {
		receiveChan <- moveMsg
	}()

	go func() {
		tracker.startProcessor()
	}()

	shouldReceive := func() bool {
		return receivedMove != nil && receivedDelete == nil
	}

	assert.Eventually(t, shouldReceive, 5*time.Second, 2*time.Second)

	ver.Verify(receivedMove)
}

func TestTracker_ReceivingDelete(t *testing.T) {
	var ver = initVerifier(t)

	var receivedDelete *DeletePayload
	var receivedMove *MovePayload
	var tracker = newServer(func(cmd *MovePayload) {
		receivedMove = cmd
	}, func(cmd *DeletePayload) {
		receivedDelete = cmd
	}, func() {
		//No-Op
	})

	receiveChan := make(chan string)
	tracker.processor = receiveChan

	go func() {
		receiveChan <- deleteMsg
	}()

	go func() {
		tracker.startProcessor()
	}()

	shouldReceive := func() bool {
		return receivedDelete != nil && receivedMove == nil
	}

	assert.Eventually(t, shouldReceive, 5*time.Second, 2*time.Second)

	ver.Verify(receivedDelete)
}
