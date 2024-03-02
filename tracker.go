package main

import (
	"fmt"
	"github.com/VerifyTests/Verify.Go/diff"
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

func newTrackedMove(temp, exe, arguments, target string, canKill bool, process int) *trackedMove {
	finder := newProjectFinder()
	project := finder.Find(target)
	extension := utils.File.GetFileExtension(target)

	if len(exe) == 0 {
		toolFinder := diff.NewTools()
		tool, found := toolFinder.TryFindByPath(exe)
		if found {
			canKill = !tool.IsMdi
		} else {
			canKill = false
		}
	}

	return &trackedMove{
		Extension: extension,
		Name:      utils.File.GetFileNameWithoutExtension(target),
		Temp:      temp,
		Target:    target,
		Exe:       exe,
		Arguments: arguments,
		CanKill:   canKill,
		Process:   process,
		Group:     project,
	}
}

type trackedDelete struct {
	Name  string
	File  string
	Group string
}

type tracker struct {
	lastCount    int
	active       Active
	inactive     Inactive
	filesDeleted map[string]*trackedDelete
	filesMoved   map[string]*trackedMove
	locker       *sync.Mutex
	processor    *processCleaner
	comparer     *fileComparer
}

func newTracker(active Active, inactive Inactive) *tracker {
	return &tracker{
		active:       active,
		inactive:     inactive,
		lastCount:    0,
		locker:       &sync.Mutex{},
		processor:    newProcessCleaner(),
		comparer:     newFileComparer(),
		filesMoved:   map[string]*trackedMove{},
		filesDeleted: map[string]*trackedDelete{},
	}
}

func (t *tracker) Start() {
	time.AfterFunc(5*time.Second, t.scanFiles)
}

func (t *tracker) scanFiles() {
	fmt.Println("updated...")
	for _, deleted := range t.filesDeleted {
		if !utils.File.Exists(deleted.File) {
			delete(t.filesDeleted, deleted.File)
		}
	}

	var count = len(t.filesMoved) + len(t.filesDeleted)
	if t.lastCount != count {
		t.toggleActive()
	}

	t.lastCount = count

	for _, moved := range t.filesMoved {
		t.handleScanMove(moved)
	}
}

func (t *tracker) moveFile(temp, target, exe, arguments string, canKill bool, processId int) {

}

func (t *tracker) addMove(temp, target, exe, arguments string, canKill bool, processId int) {
	t.locker.Lock()
	defer t.locker.Unlock()

	t.lastCount += 1

	log.Println("Tracked received move command:", temp, target, exe, arguments, canKill, processId)

	exeFile := utils.File.GetFileName(exe)
	targetFile := utils.File.GetFileName(target)

	if processId > 0 {
		t.processor.TryTerminateProcess(int32(processId))
	}

	moved := newTrackedMove(temp, exe, arguments, target, canKill, processId)

	if len(exeFile) == 0 {
		log.Printf("Move Added. Target:%s, CanKill:%t, Process:%d", targetFile, moved.CanKill, processId)
	} else {
		log.Printf("Move Added. Target:%s, CanKill:%t, Process:%d, Command: %s %s", targetFile, moved.CanKill, processId, exeFile, arguments)
	}

	if existing, ok := t.filesMoved[target]; ok {
		if processId == 0 {
			processId = existing.Process
		} else {
			//find the actual process with the Id
		}
	} else {
		t.filesMoved[target] = moved
	}
}

func (t *tracker) addDelete(filePath string) {
	t.locker.Lock()
	defer t.locker.Unlock()

	t.lastCount += 1

	log.Println("Tracked received delete command:", filePath)

	if _, ok := t.filesDeleted[filePath]; ok {
		log.Printf("DeleteUpdated. File: %s", filePath)
	} else {
		log.Printf("DeleteAdded. File: %s", filePath)
		t.filesDeleted[filePath] = &trackedDelete{
			File: filePath,
			Name: utils.File.GetFileName(filePath),
		}
	}
}

func (t *tracker) handleScanMove(moved *trackedMove) {
	if !utils.File.Exists(moved.Temp) {
		t.removeAndKill(moved)
		return
	}

	if !utils.File.Exists(moved.Target) {
		return
	}

	if !t.comparer.FilesAreEqual(moved.Temp, moved.Target) {
		return
	}

	t.removeAndKill(moved)
}

func (t *tracker) removeAndKill(moved *trackedMove) {
	if _, ok := t.filesMoved[moved.Target]; ok {
		delete(t.filesMoved, moved.Target)
		t.killProcess(moved)
	}
}

func (t *tracker) killProcess(moved *trackedMove) {
	if !moved.CanKill {
		log.Printf("did not kill for %s since CanKill=false", moved.Name)
		return
	}

	if moved.Process == 0 {
		log.Printf("no processes to kill for %s", moved.Name)
		return
	}

	t.processor.TryTerminateProcess(int32(moved.Process))
}

func (t *tracker) toggleActive() {
	if len(t.filesMoved) > 0 || len(t.filesDeleted) > 0 {
		t.active()
	} else {
		t.inactive()
	}
}

func (t *tracker) acceptDelete(del *trackedDelete) {
	log.Printf("Accepted deleted file: %s", del.Name)
}

func (t *tracker) acceptMove(mov *trackedMove) {
	log.Printf("Accepted moved file: %s", mov.Name)
}

func (t *tracker) discardMove(mov *trackedMove) {
	log.Printf("Discarded moved file: %s", mov.Name)
}

func (t *tracker) discardDelete(del *trackedDelete) {
	log.Printf("Discarded deleted file: %s", del.Name)
}
