package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/VerifyTests/Verify.Go/utils"
	"log"
)

var application fyne.App
var mainWindow fyne.Window
var optionsWindow fyne.Window
var deskApp desktop.App
var serv *server
var track *tracker

var CurrentAppIcon *fyne.StaticResource

type Action = func()

var NoAction = func() {}

func init() {
	initLogger()

	initMenu()
}

func main() {
	application = app.NewWithID("Verify.DiffTrayApp")

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

func startServer() {
	track = newTracker(showActiveIcon, showInactiveIcon, updateMenuItems)
	serv = newServer(track.addMove, track.addDelete, updateMenuItems)

	serv.Start()
	track.Start()
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

func addAcceptAllMenuItem() {
	menu := fyne.NewMenuItem(fmt.Sprintf("Accept all (%d)", track.getCount()), func() {
		track.acceptAll()
		updateMenuItems()
	})
	menu.Icon = resourceDiscardPng
	insertMenu(menu)
}

func addDiscardAllMenuItem() {
	menu := fyne.NewMenuItem(fmt.Sprintf("Discard all (%d)", track.getCount()), func() {
		track.discardAll()
		updateMenuItems()
	})
	menu.Icon = resourceDiscardPng
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

//func insertMenu(index int, menuItem *fyne.MenuItem) []*fyne.MenuItem {
//	if len(MainMenu.Items) == index { // nil or empty slice or after last element
//		return append(MainMenu.Items, menuItem)
//	}
//	MainMenu.Items = append(MainMenu.Items[:index+1], MainMenu.Items[index:]...) // index < len(a)
//	MainMenu.Items[index] = menuItem
//	return MainMenu.Items
//}

func showInactiveIcon() {
	//if CurrentAppIcon.Name() !=  {
	log.Println("Show inactive icon")
	CurrentAppIcon = resourceDefaultPng
	deskApp.SetSystemTrayIcon(CurrentAppIcon)
	//}
}

func showActiveIcon() {
	//log.Printf("Current icon name: %s", CurrentAppIconName)
	//if CurrentAppIconName != "Active" {
	log.Println("Show active icon")
	CurrentAppIcon = resourceDefaultPng
	deskApp.SetSystemTrayIcon(CurrentAppIcon)
	//}
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

func createTrayIcon() {
	CurrentAppIcon = resourceDefaultPng
	deskApp = fyne.CurrentApp().(desktop.App)
	deskApp.SetSystemTrayMenu(MainMenu)
	deskApp.SetSystemTrayIcon(CurrentAppIcon)
}

func onRaiseIssue() {
}

func onOptionsClicked() {
	log.Printf("Options menu clicked")

	if optionsWindow == nil {
		optionsWindow = fyne.CurrentApp().NewWindow("Options")
		label1 := widget.NewLabel("Version: ")
		value1 := widget.NewLabel("v1.0.0")
		grid := container.New(layout.NewFormLayout(), label1, value1)

		optionsWindow.SetContent(grid)
		optionsWindow.Resize(fyne.NewSize(480, 480))
		optionsWindow.SetCloseIntercept(func() {
			optionsWindow.Hide()
		})
	}

	optionsWindow.Show()
}