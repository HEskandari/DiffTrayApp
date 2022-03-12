package main

import (
	"net"
	"strconv"
	"testing"
	"time"
)

func TestSendingData_Integration(t *testing.T) {
	tracker := newTracker()
	tracker.Start()

	time.Sleep(2 * time.Second)

	sendMessage(t, "Test message")

	moveMsg := `{
"Type":"Move",
"CanKill":true,
"Exe":"theExe",
"Arguments":"TheArguments \"s\"",
"ProcessId":10
}`
	sendMessage(t, moveMsg)

	deleteMsg := `
{
"Type":"Move",
"Temp":"theTempFilePath",
"Target":"theTargetFilePath",
"CanKill":true,
"Exe":"theExePath",
"Arguments":"TheArguments",
"ProcessId":1000
}`
	sendMessage(t, deleteMsg)

	time.Sleep(30 * time.Second)
}

func sendMessage(t *testing.T, msg string) {
	conn, err := net.Dial("tcp", ":"+strconv.Itoa(TrackerPort))
	if err != nil {
		t.Fatalf("Failed to open connection: %s", err)
	}

	_, err = conn.Write([]byte(msg))
	if err != nil {
		t.Fatalf("Failed to write the message payload: %s", err)
	}
	_ = conn.Close()
}

func TestDeleteRequests(t *testing.T) {

	//tracker := payloadReceiver{}
	//tracker.Start()

	conn, err := net.Dial("tcp", ":"+strconv.Itoa(TrackerPort))
	if err != nil {
		t.Fatalf("Failed to open connection: %s", err)
	}

	defer conn.Close()
	message := "this is a test"
	_, err = conn.Write([]byte(message))
	if err != nil {
		t.Fatalf("Failed to write the message payload: %s", err)
	}
}