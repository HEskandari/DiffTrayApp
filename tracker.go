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

func (t *tracker) newTrackedMove(temp, exe, arguments, target string, canKill bool, process int) *trackedMove {
	project := t.finder.Find(target)
	extension := utils.File.GetFileExtension(target)

	tracked := &trackedMove{
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

	if len(exe) == 0 {
		toolFinder := diff.NewTools()
		tool, found := toolFinder.TryFindForExtension(extension)
		if found {
			tracked.CanKill = !tool.IsMdi
		} else {
			tracked.CanKill = false
		}
	}

	return tracked
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
	finder       *solutionFinder
	locker       *sync.Mutex
	processor    *processCleaner
	comparer     *fileComparer
	ticker       *time.Ticker
	stop         chan struct{}
}

func newTracker(active Active, inactive Inactive) *tracker {
	return &tracker{
		active:       active,
		inactive:     inactive,
		lastCount:    0,
		locker:       &sync.Mutex{},
		processor:    newProcessCleaner(),
		comparer:     newFileComparer(),
		finder:       newProjectFinder(),
		filesMoved:   map[string]*trackedMove{},
		filesDeleted: map[string]*trackedDelete{},
	}
}

func (t *tracker) Start() {
	t.ticker = time.NewTicker(5 * time.Second)
	t.stop = make(chan struct{})

	go func() {
		for {
			select {
			case <-t.ticker.C:
				t.scanFiles()
			case <-t.stop:
				return
			}
		}
	}()
}

func (t *tracker) Stop() {
	t.ticker.Stop()
	close(t.stop)
}

func (t *tracker) scanFiles() {
	fmt.Println("Scanning...")
	for _, deleted := range t.filesDeleted {
		if !utils.File.Exists(deleted.File) {
			delete(t.filesDeleted, deleted.File)
		}
	}

	t.lastCount = len(t.filesMoved) + len(t.filesDeleted)
	t.toggleActive()

	for _, moved := range t.filesMoved {
		t.handleScanMove(moved)
	}
}

func (t *tracker) trackingAny() bool {
	return len(t.filesMoved) > 0 || len(t.filesDeleted) > 0
}

func (t *tracker) addMove(move *MovePayload) {
	t.locker.Lock()
	defer t.locker.Unlock()

	log.Println("Tracked received move command:", move.Temp, move.Target, move.Exe, move.Arguments, move.CanKill, move.ProcessId)

	exeFile := utils.File.GetFileName(move.Exe)
	targetFile := utils.File.GetFileName(move.Target)

	if _, ok := t.filesMoved[move.Target]; !ok {

		//if processId == 0 {
		//	processId = existing.Process
		//} else {
		//	//find the actual process with the Id
		//}

		moved := t.newTrackedMove(move.Temp, move.Exe, move.Arguments, move.Target, move.CanKill, move.ProcessId)

		if len(exeFile) == 0 {
			log.Printf("Move Added. Target:%s, CanKill:%t, Process:%d", targetFile, moved.CanKill, move.ProcessId)
		} else {
			log.Printf("Move Added. Target:%s, CanKill:%t, Process:%d, Command: %s %s", targetFile, moved.CanKill, move.ProcessId, exeFile, move.Arguments)
		}

		t.filesMoved[move.Target] = moved
	} else {
		//update
		//t.filesMoved[target] = moved
	}

	//if processId > 0 {
	//	t.processor.TryTerminateProcess(int32(processId))
	//}

	//return t.filesMoved[move.Target]
}

func (t *tracker) addDelete(delete *DeletePayload) {
	t.locker.Lock()
	defer t.locker.Unlock()

	log.Println("Tracked received delete command:", delete.File)

	if _, ok := t.filesDeleted[delete.File]; !ok {
		log.Printf("DeleteUpdated. File: %s", delete.File)
		solution := t.finder.Find(delete.File)
		deleted := &trackedDelete{
			File:  delete.File,
			Name:  utils.File.GetFileName(delete.File),
			Group: solution,
		}
		t.filesDeleted[delete.File] = deleted
	}
	//return t.filesDeleted[filePath]
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
	if t.trackingAny() {
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
