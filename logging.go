package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
)

var LogFileName = "DiffTrayApp.Logs.txt"
var logFile *os.File

func initLogger() {
	logDir, err := appHomeDir()
	if err != nil {
		panic("Could not get the Application directory")
	}

	err = safeCreateDirectory(logDir)
	if err != nil {
		panic("Could not create logs directory: " + logDir)
	}
	logFile, err := os.OpenFile(logDir+"/"+LogFileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		panic("Failed to open log file: " + LogFileName + ". Error: " + err.Error())
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

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
	appHome, err := appHomeDir()
	if err != nil {
		return
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "darwin" {
		cmd = exec.Command("open", appHome)
	} else {
		cmd = exec.Command("explorer", appHome)
	}

	if cmd != nil {
		err := cmd.Start()
		if err != nil {
			log.Printf("Failed to open log directory: " + appHome)
		}
	}
}

func appHomeDir() (string, error) {
	home, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return path.Join(home, "DiffTray"), nil
}
