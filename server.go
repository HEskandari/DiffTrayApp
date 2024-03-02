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

var TrackerPort = 3492

type MoveMessageHandler func(cmd *MovePayload)
type DeleteMessageHandler func(cmd *DeletePayload)

type server struct {
	processor     chan string
	moveHandler   MoveMessageHandler
	deleteHandler DeleteMessageHandler
}

func newServer() *server {
	srv := &server{
		processor: make(chan string, 1),
	}
	return srv
}

func (s *server) Start() {
	go s.startReceiver()
	go s.startProcessor(s.processor)
}

func (s *server) Stop() {
	close(s.processor)
}

func (s *server) startProcessor(input <-chan string) {
	for {
		message := <-input

		if strings.Contains(message, "\"Type\":\"Move\"") {
			log.Printf("Move message received: %s", message)
			moveCommand := MovePayload{}
			deserialize(message, &moveCommand)
			s.moveFile(&moveCommand)
		} else if strings.Contains(message, "\"Type\":\"Delete\"") {
			log.Printf("Delete message received: %s", message)
			deleteCommand := DeletePayload{}
			deserialize(message, &deleteCommand)
			s.deleteFile(&deleteCommand)
		} else {
			log.Printf("Unknown message: %s", message)
		}
	}
}

func (s *server) startReceiver() {
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

			if len(message) > 0 {
				log.Printf("Sending message to processor: %s", message)
				s.processor <- message
			} else {
				log.Printf("Received empty message. Dropping")
			}
		}(conn)
	}
}

func check(err error, message string) {
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", message)
}

func (s *server) deleteFile(command *DeletePayload) {
	log.Printf("Deleting command received: %+v", command)
	s.deleteHandler(command)
}

func (s *server) moveFile(command *MovePayload) {
	log.Printf("Moving command received: %+v", command)
	s.moveHandler(command)
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
