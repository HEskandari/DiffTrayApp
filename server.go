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

var TrackerPort = 4523

type MoveMessageHandler func(cmd *MovePayload)
type DeleteMessageHandler func(cmd *DeletePayload)

type server struct {
	processor     chan string
	moveHandler   MoveMessageHandler
	deleteHandler DeleteMessageHandler
}

func (t *server) Start() {
	go t.startReceiver()
	go t.startProcessor(t.processor)
}

func (t *server) Stop() {
	close(t.processor)
}

func (t *server) startProcessor(input <-chan string) {
	for {
		message := <-input

		if strings.Contains(message, "\"Type\":\"Move\"") {
			moveCommand := MovePayload{}
			deserialize(message, &moveCommand)
			t.moveFile(&moveCommand)
		} else if strings.Contains(message, "\"Type\":\"Delete\"") {
			deleteCommand := DeletePayload{}
			deserialize(message, &deleteCommand)
			t.deleteFile(&deleteCommand)
		} else {
			log.Printf("Unknown payload: %s", message)
		}
	}
}

func newServer() *server {
	srv := &server{
		processor: make(chan string, 1),
	}
	return srv
}

func (t *server) startReceiver() {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(TrackerPort))
	check(err, "Server is ready.")

	for {
		conn, err := listener.Accept()
		check(err, "Accepted connection.")

		go func(reader io.Reader) {
			bytes, err := ioutil.ReadAll(reader)
			if err != nil {
				log.Printf("Failed to read: %s", err)
				return
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

func (t *server) deleteFile(command *DeletePayload) {
	log.Printf("Delete: %+v", command)
	t.deleteHandler(command)
}

func (t *server) moveFile(command *MovePayload) {
	log.Printf("Move: %+v", command)
	t.moveHandler(command)
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