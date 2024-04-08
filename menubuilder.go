package main

import (
	"fyne.io/fyne/v2"
	"github.com/VerifyTests/Verify.Go/utils"
	"log"
	"slices"
)

var MainMenu *fyne.Menu
var MenuOptions *fyne.MenuItem
var MenuLogs *fyne.MenuItem
var MenuIssues *fyne.MenuItem
var MenuQuitApp *fyne.MenuItem

//var StartSeparator =
//var EndSeparator = fyne.NewMenuItemSeparator()

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

func createMainMenu() {
	MainMenu = fyne.NewMenu("Main Menu", MenuOptions, MenuLogs, MenuIssues, fyne.NewMenuItemSeparator(), MenuQuitApp)
}

func updateMenuItems() {
	log.Printf("Updating menu")

	clearFileMenus()

	if track.trackingAny() {
		insertSeparator()

		addDiscardAllMenuItem()
		addAcceptAllMenuItem()

		insertSeparator()

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

		insertSeparator()
	}

	MainMenu.Refresh()
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
	if menuIssuesIndex == -1 {
		return
	}

	MainMenu.Items = append(MainMenu.Items[:menuIssuesIndex+1], MainMenu.Items[quitOption-2:]...)
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

func insertSeparator() {
	insertMenu(fyne.NewMenuItemSeparator())
}

//func addEndSeparator() {
//	endIndex := findIndex(MainMenu.Items, func(item *fyne.MenuItem) bool {
//		return item == EndSeparator
//	})
//
//	if endIndex == -1 {
//		insertMenu(EndSeparator)
//	}
//}

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