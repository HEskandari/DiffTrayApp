package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/cmd/fyne_demo/tutorials"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/VerifyTests/Verify.Go/utils"
	"log"
)

var mainWindow fyne.Window
var mainMenu *fyne.Menu
var deskApp desktop.App
var serv *server
var lastMenuIndex = 3

type Action = func()

var NoAction = func() {}
var StartSeparator = fyne.NewMenuItemSeparator()
var EndSeparator = fyne.NewMenuItemSeparator()

func main() {
	application := app.NewWithID("Verify.DiffTrayApp")
	application.SetIcon(resourceCogsPng)

	createMainMenu()
	createTrayIcon()
	startServer()

	mainWindow = application.NewWindow("Main")
	mainWindow.Resize(fyne.NewSize(640, 460))
	mainWindow.SetCloseIntercept(func() {
		mainWindow.Hide() //prevent main window from closing
	})
	application.Run()
}

func startServer() {
	tracker := newTracker(showActiveIcon, showInactiveIcon)

	serv = newServer()
	serv.moveHandler = func(cmd *MovePayload) {
		tracker.addMove(cmd.Temp, cmd.Target, cmd.Exe, cmd.Arguments, cmd.CanKill, cmd.ProcessId)
		updateMenuItems(tracker)
	}
	serv.deleteHandler = func(cmd *DeletePayload) {
		tracker.addDelete(cmd.File)
		updateMenuItems(tracker)
	}

	serv.Start()
}

func updateMenuItems(t *tracker) {
	log.Printf("Updating menu")

	//separator := -1
	//for i, item := range mainMenu.Items {
	//	if item.IsSeparator {
	//		separator = i
	//		break
	//	}
	//}
	//
	//if separator > -1 {
	//	mainMenu.Items = append(mainMenu.Items[:separator+1], mainMenu.Items[separator+t.lastCount:]...)
	//}

	if t.lastCount > 0 {
		addStartSeparator()

		//Add delete items
		for d := range t.filesDeleted {
			td := t.filesDeleted[d]
			addDeleteMenuItem(lastMenuIndex+1, td, func() {
				t.acceptDelete(td)
			})
		}

		//Add moved items
		for m := range t.filesMoved {
			tm := t.filesMoved[m]
			addMovedMenuItem(lastMenuIndex+1, tm, func() {
				t.acceptMove(tm)
			}, func() {
				t.discardMove(tm)
			})
		}

		addEndSeparator()
	}

	mainMenu.Refresh()
}

func addMovedMenuItem(index int, move *trackedMove, accept Action, discard Action) {
	tempName := utils.File.GetFileNameWithoutExtension(move.Temp)
	targetName := utils.File.GetFileNameWithoutExtension(move.Target)
	text := getMoveText(move, tempName, targetName)

	menu := fyne.NewMenuItem(text, NoAction)
	mainMenu.Items = insertMenu(index, menu)
}

func getMoveText(move *trackedMove, temp string, target string) string {
	if utils.File.GetFileNameWithoutExtension(temp) == utils.File.GetFileNameWithoutExtension(target) {
		return fmt.Sprintf("%s (%s)", move.Name, move.Extension)
	}
	return fmt.Sprintf("%s > %s (%s)", temp, target, move.Extension)
}

func addDeleteMenuItem(index int, move *trackedDelete, action Action) {
	menu := fyne.NewMenuItem(move.Name, action)
	menu.Icon = resourceAcceptPng
	mainMenu.Items = insertMenu(index, menu)
}

func addStartSeparator() {
	for _, item := range mainMenu.Items {
		if item == StartSeparator {
			return
		}
	}

	if !mainMenu.Items[lastMenuIndex+1].IsSeparator {
		mainMenu.Items = insertMenu(lastMenuIndex+1, StartSeparator)
	}
}

func addEndSeparator() {
	for _, item := range mainMenu.Items {
		if item == EndSeparator {
			return
		}
	}

	if !mainMenu.Items[len(mainMenu.Items)-1].IsSeparator {
		mainMenu.Items = insertMenu(len(mainMenu.Items)-1, EndSeparator)
	}
}

func insertMenu(index int, menuItem *fyne.MenuItem) []*fyne.MenuItem {
	if len(mainMenu.Items) == index { // nil or empty slice or after last element
		return append(mainMenu.Items, menuItem)
	}
	mainMenu.Items = append(mainMenu.Items[:index+1], mainMenu.Items[index:]...) // index < len(a)
	mainMenu.Items[index] = menuItem
	return mainMenu.Items
}

func showInactiveIcon() {
	log.Println("Show inactive icon")
	deskApp.SetSystemTrayIcon(resourceDefaultPng)
}

func showActiveIcon() {
	log.Println("Show active icon")
	deskApp.SetSystemTrayIcon(resourceActivePng)
}

func logLifecycle(a fyne.App) {
	a.Lifecycle().SetOnStarted(func() {
		log.Println("Lifecycle: Started")
	})
	a.Lifecycle().SetOnStopped(func() {
		log.Println("Lifecycle: Stopped")
		appStopped()
	})
	a.Lifecycle().SetOnEnteredForeground(func() {
		log.Println("Lifecycle: Entered Foreground")
	})
	a.Lifecycle().SetOnExitedForeground(func() {
		log.Println("Lifecycle: Exited Foreground")
	})
}

func appStopped() {
	serv.Stop()
}

func makeNav(setTutorial func(tutorial tutorials.Tutorial), loadPrevious bool) fyne.CanvasObject {
	a := fyne.CurrentApp()
	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return tutorials.TutorialIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := tutorials.TutorialIndex[uid]
			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := tutorials.Tutorials[uid]
			if !ok {
				fyne.LogError("Missing tutorial panel: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.Title)
			obj.(*widget.Label).TextStyle = fyne.TextStyle{}
		},
		OnSelected: func(uid string) {
			if _, ok := tutorials.Tutorials[uid]; ok {
				return
				//if unsupportedTutorial(t) {
				//	return
				//}
				//a.Preferences().SetString(preferenceCurrentTutorial, uid)
				//setTutorial(t)
			}
		},
	}
	if loadPrevious {
		//currentPref := a.Preferences().StringWithFallback(preferenceCurrentTutorial, "welcome")
		//tree.Select(currentPref)
	}
	themes := container.NewGridWithColumns(2,
		widget.NewButton("Dark", func() {
			a.Settings().SetTheme(theme.DarkTheme())
		}),
		widget.NewButton("Light", func() {
			a.Settings().SetTheme(theme.LightTheme())
		}),
	)
	return container.NewBorder(nil, themes, nil, nil, tree)
}

func shortcutFocused(s fyne.Shortcut, w fyne.Window) {
	if focused, ok := w.Canvas().Focused().(fyne.Shortcutable); ok {
		focused.TypedShortcut(s)
	}
}

func createTrayIcon() {
	deskApp = fyne.CurrentApp().(desktop.App)
	deskApp.SetSystemTrayMenu(mainMenu)
	deskApp.SetSystemTrayIcon(resourceCogsPng)
}

func createMainMenu() {
	mainMenu = fyne.NewMenu("Main Menu",
		fyne.NewMenuItem("Options", onOptionsClicked),
		fyne.NewMenuItem("Open logs", onOpenLogs),
		fyne.NewMenuItem("Raise issue", onRaiseIssue))
}

func onRaiseIssue() {
}

func onOpenLogs() {
}

func onOptionsClicked() {
}
