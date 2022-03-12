package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"
)

type payloadReceiver struct {
}

type tracker struct {
	receiver  *payloadReceiver
	processor chan string
}

func (t *tracker) Start() {
	go t.startReceiver()
	go t.startProcessor(t.processor)
}

func (t *tracker) startProcessor(input <-chan string) {
	for {
		message := <-input

		if strings.Contains(message, "\"Type\":\"Move\"") {
			moveCommand := MovePayload{}
			deserialize(message, &moveCommand)
			moveFile(moveCommand)
		} else if strings.Contains(message, "\"Type\":\"Delete\"") {
			deleteCommand := DeletePayload{}
			deserialize(message, &deleteCommand)
			deleteFile(deleteCommand)
		} else {
			log.Printf("Unknown payload: %s", message)
		}
	}
}

func newTracker() tracker {
	tracker := tracker{
		processor: make(chan string, 1),
	}
	return tracker
}

var TrackerPort = 4523

func (t *tracker) startReceiver() {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(TrackerPort))
	check(err, "Server is ready.")

	for {
		conn, err := listener.Accept()
		check(err, "Accepted connection.")

		go func(reader io.Reader) {
			bytes, err := ioutil.ReadAll(reader)
			if err != nil {
				log.Printf("Failed to read: %s", err)
			}
			message := string(bytes)

			t.processor <- message
		}(conn)
	}
}

func check(err error, message string) {
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", message)
}

func deleteFile(command DeletePayload) {
	log.Printf("Delete: %+v", command)
}

func moveFile(command MovePayload) {
	log.Printf("Move: %+v", command)
}

func deserialize(payload string, obj interface{}) {
	err := json.Unmarshal([]byte(payload), obj)
	if err != nil {
		fmt.Println(err)
	}
}

type MovePayload struct {
	Temp      string
	Target    string
	Exe       string
	Arguments string
	CanKill   bool
	ProcessId int
}

type DeletePayload struct {
	File string
}