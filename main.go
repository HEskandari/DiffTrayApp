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
	"os"
	"slices"
)

var mainWindow fyne.Window
var mainMenu *fyne.Menu
var deskApp desktop.App
var serv *server
var track *tracker

//var lastMenuIndex = 2

var menuOptions *fyne.MenuItem
var menuLogs *fyne.MenuItem
var menuIssues *fyne.MenuItem

type Action = func()

var NoAction = func() {}
var StartSeparator = fyne.NewMenuItemSeparator()
var EndSeparator = fyne.NewMenuItemSeparator()
var CurrentAppIcon *fyne.StaticResource

func init() {
	initLogger()

	initMenu()
}

func initLogger() {
	file, err := os.OpenFile("Verify.Logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)
}

func initMenu() {
	menuOptions = fyne.NewMenuItem("Options", onOptionsClicked)
	menuOptions.Icon = resourceCogsPng

	menuLogs = fyne.NewMenuItem("Open logs", onOpenLogs)
	menuLogs.Icon = resourceFolderPng

	menuIssues = fyne.NewMenuItem("Raise issue", onRaiseIssue)
	menuIssues.Icon = resourceLinkPng
}

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
	track = newTracker(showActiveIcon, showInactiveIcon, updateMenuItems)
	serv = newServer(track.addMove, track.addDelete, updateMenuItems)

	serv.Start()
	track.Start()
}

func updateMenuItems() {
	log.Printf("Updating menu")

	clearFileMenus()

	if track.trackingAny() {
		addStartSeparator()

		addAcceptAllMenuItem()
		addClearAllMenuItem()

		addStartSeparator()

		//Add delete items
		for d := range track.filesDeleted {
			td := track.filesDeleted[d]
			addDeleteMenuItem(td, func() {
				track.acceptDelete(td)
			})
		}

		//Add moved items
		for m := range track.filesMoved {
			tm := track.filesMoved[m]
			addMovedMenuItem(tm, func() {
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

	if len(mainMenu.Items) == 5 {
		return
	}

	quitOption := findIndex(mainMenu.Items, func(item *fyne.MenuItem) bool {
		return item.IsQuit
	})

	menuIssuesIndex := findIndex(mainMenu.Items, func(item *fyne.MenuItem) bool {
		return item == menuIssues
	})

	if quitOption == -1 {
		return
	}

	//if lastMenuIndex == quitOption-2 {
	//	//No file options are added
	//	return
	//}
	//mainMenu.Items = removeElementByRange(mainMenu.Items, menuIssuesIndex+1, quitOption-2)
	mainMenu.Items = append(mainMenu.Items[:menuIssuesIndex+1], mainMenu.Items[quitOption-2:]...)
}

func addAcceptAllMenuItem() {
	menu := fyne.NewMenuItem(fmt.Sprintf("Accept all (%d)", track.getCount()), track.acceptAll)
	menu.Icon = resourceDiscardPng
	insertMenu(menu)
}

func addClearAllMenuItem() {
	menu := fyne.NewMenuItem(fmt.Sprintf("Discard (%d)", track.getCount()), track.clear)
	menu.Icon = resourceDiscardPng
	insertMenu(menu)
}

func addMovedMenuItem(move *trackedMove, accept Action, discard Action) {
	tempName := utils.File.GetFileNameWithoutExtension(move.Temp)
	targetName := utils.File.GetFileNameWithoutExtension(move.Target)
	text := getMoveText(move, tempName, targetName)

	menu := fyne.NewMenuItem(text, NoAction)
	insertMenu(menu)
}

func getMoveText(move *trackedMove, temp string, target string) string {
	if utils.File.GetFileNameWithoutExtension(temp) == utils.File.GetFileNameWithoutExtension(target) {
		return fmt.Sprintf("%s (%s)", move.Name, move.Extension)
	}
	return fmt.Sprintf("%s > %s (%s)", temp, target, move.Extension)
}

func addDeleteMenuItem(move *trackedDelete, action Action) {
	menu := fyne.NewMenuItem(move.Name, action)
	menu.Icon = resourceAcceptPng
	insertMenu(menu)
}

func addStartSeparator() {
	//endIndex := findIndex(mainMenu.Items, func(item *fyne.MenuItem) bool {
	//	return item.IsQuit
	//})
	//
	//startIndex := findIndex(mainMenu.Items, func(item *fyne.MenuItem) bool {
	//	return item == StartSeparator
	//})

	//if startIndex == -1 {
	//lastMenuIndex += 1
	insertMenu(StartSeparator)
	//}
}

func addEndSeparator() {
	endIndex := findIndex(mainMenu.Items, func(item *fyne.MenuItem) bool {
		return item == EndSeparator
	})

	if endIndex == -1 {
		insertMenu(EndSeparator)
	}
}

func insertMenu(menuItem *fyne.MenuItem) {
	//endIndex := findIndex(mainMenu.Items, func(item *fyne.MenuItem) bool {
	//	return item.IsQuit
	//})

	menuIssuesIndex := findIndex(mainMenu.Items, func(item *fyne.MenuItem) bool {
		return item == menuIssues
	})

	//if len(mainMenu.Items) == index { // nil or empty slice or after last element
	//	return append(mainMenu.Items, menuItem)
	//}
	mainMenu.Items = slices.Insert(mainMenu.Items, menuIssuesIndex+1, menuItem)
	//mainMenu.Items = append(mainMenu.Items[:menuIssuesIndex+1], menuItem, mainMenu.Items[menuIssuesIndex+2:]...)
	//mainMenu.Items[index] = menuItem
	//return mainMenu.Items
}

//func insertMenu(index int, menuItem *fyne.MenuItem) []*fyne.MenuItem {
//	if len(mainMenu.Items) == index { // nil or empty slice or after last element
//		return append(mainMenu.Items, menuItem)
//	}
//	mainMenu.Items = append(mainMenu.Items[:index+1], mainMenu.Items[index:]...) // index < len(a)
//	mainMenu.Items[index] = menuItem
//	return mainMenu.Items
//}

func showInactiveIcon() {
	if CurrentAppIcon != resourceDefaultPng {
		log.Println("Show inactive icon")
		CurrentAppIcon = resourceDefaultPng
		deskApp.SetSystemTrayIcon(CurrentAppIcon)
	}
}

func showActiveIcon() {
	if CurrentAppIcon != resourceActivePng {
		log.Println("Show active icon")
		CurrentAppIcon = resourceActivePng
		deskApp.SetSystemTrayIcon(CurrentAppIcon)
	}
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
	CurrentAppIcon = resourceDefaultPng
	deskApp = fyne.CurrentApp().(desktop.App)
	deskApp.SetSystemTrayMenu(mainMenu)
	deskApp.SetSystemTrayIcon(CurrentAppIcon)
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
