package main

import (
	"github.com/VerifyTests/Verify.Go/utils"
	"log"
	"os"
	"os/exec"
	"runtime"
)

var LogsDirectory string = "Logs"
var LogFileName = "Verify.Logs.txt"
var logFile *os.File

func initLogger() {
	err := utils.File.CreateDirectory(LogsDirectory)
	if err != nil {
		panic("Could not create logs directory: " + LogsDirectory)
	}
	logFile, err := os.OpenFile(LogsDirectory+"/"+LogFileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Panic("Failed to open log file: " + LogFileName + ". Error: " + err.Error())
	}

	// Set log out put and enjoy :)
	log.SetOutput(logFile)

	// optional: log date-time, filename, and line number
	log.SetFlags(log.LstdFlags)
}

func closeLogFile() {
	if logFile != nil {
		err := logFile.Close()
		if err != nil {
			log.Panic("Failed to close log file: " + LogFileName + ". Error: " + err.Error())
		}
	}
}

func openLogDirectory() {
	var cmd *exec.Cmd
	if runtime.GOOS == "darwin" {
		cmd = exec.Command("open", LogsDirectory)
	} else {
		cmd = exec.Command("explorer", LogsDirectory)
	}

	if cmd != nil {
		err := cmd.Start()
		if err != nil {
			log.Printf("Failed to open log directory: " + LogsDirectory)
		}
	}
}
