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
var track *tracker
var lastMenuIndex = 2

var menuOptions = fyne.NewMenuItem("Options", onOptionsClicked)
var menuLogs = fyne.NewMenuItem("Open logs", onOpenLogs)
var menuIssues = fyne.NewMenuItem("Raise issue", onRaiseIssue)

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
	track = newTracker(showActiveIcon, showInactiveIcon)

	serv = newServer(track.addMove, track.addDelete, updateMenuItems)
	serv.Start()
	track.Start()
}

func updateMenuItems() {
	log.Printf("Updating menu")

	clearFileMenus()

	if track.trackingAny() {
		addStartSeparator()

		//Add delete items
		for d := range track.filesDeleted {
			td := track.filesDeleted[d]
			addDeleteMenuItem(lastMenuIndex+1, td, func() {
				track.acceptDelete(td)
			})
		}

		//Add moved items
		for m := range track.filesMoved {
			tm := track.filesMoved[m]
			addMovedMenuItem(lastMenuIndex+1, tm, func() {
				track.acceptMove(tm)
			}, func() {
				track.discardMove(tm)
			})
		}

		addEndSeparator()
	}

	mainMenu.Refresh()
}

func clearFileMenus() {
	quitOption := findIndex(mainMenu.Items, func(item *fyne.MenuItem) bool {
		return item.IsQuit
	})

	if quitOption == -1 {
		return
	}

	if lastMenuIndex == quitOption-2 {
		//No file options are added
		return
	}
	mainMenu.Items = removeElementByRange(mainMenu.Items, lastMenuIndex+1, quitOption-2)
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
	startIndex := findIndex(mainMenu.Items, func(item *fyne.MenuItem) bool {
		return item == StartSeparator
	})

	if startIndex == -1 {
		mainMenu.Items = insertMenu(lastMenuIndex+1, StartSeparator)
		lastMenuIndex += 1
	}
}

func addEndSeparator() {
	endIndex := findIndex(mainMenu.Items, func(item *fyne.MenuItem) bool {
		return item == EndSeparator
	})

	if endIndex == -1 {
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
	track.Stop()
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
	mainMenu = fyne.NewMenu("Main Menu", menuOptions, menuLogs, menuIssues)
}

func onRaiseIssue() {
}

func onOpenLogs() {
}

func onOptionsClicked() {
}
