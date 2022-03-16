package main

import (
	"fmt"
	"github.com/VerifyTests/Verify.Go/utils"
	"log"
	"sync"
	"time"
)

type Active func()
type Inactive func()

type trackedMove struct {
	Extension string
	Name      string
	Temp      string
	Target    string
	Exe       string
	Arguments string
	CanKill   bool
	Process   int
	Group     string
}

type trackedDelete struct {
	Name  string
	File  string
	Group string
}

type tracker struct {
	active       Active
	inactive     Inactive
	filesDeleted map[string]trackedDelete
	filesMoved   map[string]trackedMove
	locker       *sync.Mutex
}

func newTracker(active Active, inactive Inactive) *tracker {
	return &tracker{
		active:   active,
		inactive: inactive,
		locker:   &sync.Mutex{},
	}
}

func (t *tracker) Start() {
	time.AfterFunc(5*time.Second, t.updateLoop)
}

func (t *tracker) updateLoop() {
	fmt.Println("updated...")
}

func (t *tracker) moveFile(temp, target, exe, arguments string, canKill bool, processId int) {

}

func (t *tracker) addMove(temp, target, exe, arguments string, kill bool, processId int) {
	t.locker.Lock()
	defer t.locker.Unlock()

	exeFile := utils.File.GetFileName(exe)
	targetFile := utils.File.GetFileName(target)

	if processId > 0 {
		tryTerminateProcess(int32(processId))
	}

	t.filesMoved[target] = trackedMove{
		Exe:    exeFile,
		Target: targetFile,
	}
}

func (t *tracker) addDelete(filePath string) {
	t.locker.Lock()
	defer t.locker.Unlock()

	if _, ok := t.filesDeleted[filePath]; !ok {
		log.Printf("DeleteAdded. File: %s", filePath)
		t.filesDeleted[filePath] = trackedDelete{
			File: filePath,
			Name: utils.File.GetFileName(filePath),
		}
	} else {
		log.Printf("DeleteUpdated. File: %s", filePath)

	}
}