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
	"slices"
)

var mainWindow fyne.Window
var deskApp desktop.App
var serv *server
var track *tracker

var MainMenu *fyne.Menu
var MenuOptions *fyne.MenuItem
var MenuLogs *fyne.MenuItem
var MenuIssues *fyne.MenuItem
var MenuQuitApp *fyne.MenuItem
var StartSeparator = fyne.NewMenuItemSeparator()
var EndSeparator = fyne.NewMenuItemSeparator()

var CurrentAppIcon *fyne.StaticResource

type Action = func()

var NoAction = func() {}

func init() {
	initLogger()

	initMenu()
}

func debugMenu() {
	for i, item := range MainMenu.Items {
		println(i, item.Label)
	}
}

func initMenu() {
	MenuOptions = fyne.NewMenuItem("Options", onOptionsClicked)
	MenuOptions.Icon = resourceCogsPng

	MenuLogs = fyne.NewMenuItem("Open logs", openLogDirectory)
	MenuLogs.Icon = resourceFolderPng

	MenuIssues = fyne.NewMenuItem("Raise issue", onRaiseIssue)
	MenuIssues.Icon = resourceLinkPng

	MenuQuitApp = fyne.NewMenuItem("Quit", nil)
	MenuQuitApp.IsQuit = true
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

	registerAndRun(application)
}

func createMainMenu() {
	MainMenu = fyne.NewMenu("Main Menu", MenuOptions, MenuLogs, MenuIssues, EndSeparator, MenuQuitApp)
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

		addDiscardAllMenuItem()
		addAcceptAllMenuItem()

		addStartSeparator()

		//Add delete items
		for d := range track.filesDeleted {
			td := track.filesDeleted[d]
			addDeleteMenuItem(td, func() {
				acceptDelete(td)
			})
		}

		//Add moved items
		for m := range track.filesMoved {
			tm := track.filesMoved[m]
			addMovedMenuItem(tm, func() {
				acceptMove(tm)
			}, func() {
				discardMove(tm)
			})
		}

		addEndSeparator()
	}

	MainMenu.Refresh()
}

func discardMove(tm *trackedMove) {
	track.discardMove(tm)
	updateMenuItems()
}

func acceptMove(tm *trackedMove) {
	track.acceptMove(tm)
	updateMenuItems()
}

func acceptDelete(td *trackedDelete) {
	track.acceptDelete(td)
	updateMenuItems()
}

func clearFileMenus() {

	if len(MainMenu.Items) == 5 {
		return
	}

	quitOption := findIndex(MainMenu.Items, func(item *fyne.MenuItem) bool {
		return item.IsQuit
	})

	menuIssuesIndex := findIndex(MainMenu.Items, func(item *fyne.MenuItem) bool {
		return item == MenuIssues
	})

	if quitOption == -1 {
		return
	}

	MainMenu.Items = append(MainMenu.Items[:menuIssuesIndex+1], MainMenu.Items[quitOption-2:]...)
}

func addAcceptAllMenuItem() {
	menu := fyne.NewMenuItem(fmt.Sprintf("Accept all (%d)", track.getCount()), track.acceptAll)
	menu.Icon = resourceDiscardPng
	insertMenu(menu)
}

func addDiscardAllMenuItem() {
	menu := fyne.NewMenuItem(fmt.Sprintf("Discard (%d)", track.getCount()), track.discard)
	menu.Icon = resourceDiscardPng
	insertMenu(menu)
}

func addMovedMenuItem(move *trackedMove, accept Action, discard Action) {
	tempName := utils.File.GetFileNameWithoutExtension(move.Temp)
	targetName := utils.File.GetFileNameWithoutExtension(move.Target)
	text := getMoveText(move, tempName, targetName)

	menu := fyne.NewMenuItem(text, NoAction)
	insertMenu(menu)

	menu.ChildMenu = fyne.NewMenu("",
		fyne.NewMenuItem("Accept move", accept),
		fyne.NewMenuItem("Discard", discard))

	if len(move.Exe) > 0 {
		menu.ChildMenu.Items = append(menu.ChildMenu.Items,
			fyne.NewMenuItem("Open diff tool", launchDiffTool))
	}
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
	insertMenu(StartSeparator)
}

func addEndSeparator() {
	endIndex := findIndex(MainMenu.Items, func(item *fyne.MenuItem) bool {
		return item == EndSeparator
	})

	if endIndex == -1 {
		insertMenu(EndSeparator)
	}
}

func insertMenu(menuItem *fyne.MenuItem) {
	//endIndex := findIndex(MainMenu.Items, func(item *fyne.MenuItem) bool {
	//	return item.IsQuit
	//})

	menuIssuesIndex := findIndex(MainMenu.Items, func(item *fyne.MenuItem) bool {
		return item == MenuIssues
	})

	//if len(MainMenu.Items) == index { // nil or empty slice or after last element
	//	return append(MainMenu.Items, menuItem)
	//}
	MainMenu.Items = slices.Insert(MainMenu.Items, menuIssuesIndex+1, menuItem)
	//MainMenu.Items = append(MainMenu.Items[:menuIssuesIndex+1], menuItem, MainMenu.Items[menuIssuesIndex+2:]...)
	//MainMenu.Items[index] = menuItem
	//return MainMenu.Items
}

//func insertMenu(index int, menuItem *fyne.MenuItem) []*fyne.MenuItem {
//	if len(MainMenu.Items) == index { // nil or empty slice or after last element
//		return append(MainMenu.Items, menuItem)
//	}
//	MainMenu.Items = append(MainMenu.Items[:index+1], MainMenu.Items[index:]...) // index < len(a)
//	MainMenu.Items[index] = menuItem
//	return MainMenu.Items
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

func registerAndRun(application fyne.App) {
	application.Lifecycle().SetOnStarted(func() {
		log.Println("Application: Started")
	})
	application.Lifecycle().SetOnStopped(func() {
		log.Println("Application: Stopped")
		appStopped()
	})
	application.Lifecycle().SetOnEnteredForeground(func() {
		log.Println("Application: Entered Foreground")
	})
	application.Lifecycle().SetOnExitedForeground(func() {
		log.Println("Application: Exited Foreground")
	})
	application.Run()
}

func appStopped() {
	serv.Stop()
	track.Stop()
	closeLogFile()
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
	deskApp.SetSystemTrayMenu(MainMenu)
	deskApp.SetSystemTrayIcon(CurrentAppIcon)
}

func onRaiseIssue() {
}

func onOptionsClicked() {
}
