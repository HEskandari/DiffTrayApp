package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/driver/desktop"
	"net/url"
)

func makeMenu() *fyne.MainMenu {
	newItem := fyne.NewMenuItem("New", nil)
	checkedItem := fyne.NewMenuItem("Checked", nil)
	checkedItem.Checked = true
	disabledItem := fyne.NewMenuItem("Disabled", nil)
	disabledItem.Disabled = true
	otherItem := fyne.NewMenuItem("Other", nil)
	otherItem.ChildMenu = fyne.NewMenu("",
		fyne.NewMenuItem("Project", func() { fmt.Println("Menu New->Other->Project") }),
		fyne.NewMenuItem("Mail", func() { fmt.Println("Menu New->Other->Mail") }),
	)
	newItem.ChildMenu = fyne.NewMenu("",
		fyne.NewMenuItem("File", func() { fmt.Println("Menu New->File") }),
		fyne.NewMenuItem("Directory", func() { fmt.Println("Menu New->Directory") }),
		otherItem,
	)
	openSettings := func() {
		w := fyne.CurrentApp().NewWindow("Fyne Settings")
		w.SetContent(settings.NewSettings().LoadAppearanceScreen(w))
		w.Resize(fyne.NewSize(480, 480))
		w.Show()
	}
	settingsItem := fyne.NewMenuItem("Settings", openSettings)
	settingsShortcut := &desktop.CustomShortcut{KeyName: fyne.KeyComma, Modifier: fyne.KeyModifierShortcutDefault}
	settingsItem.Shortcut = settingsShortcut

	mainWindow.Canvas().AddShortcut(settingsShortcut, func(shortcut fyne.Shortcut) {
		openSettings()
	})
	cutShortcut := &fyne.ShortcutCut{Clipboard: mainWindow.Clipboard()}
	cutItem := fyne.NewMenuItem("Cut", func() {
		shortcutFocused(cutShortcut, mainWindow)
	})
	cutItem.Shortcut = cutShortcut
	copyShortcut := &fyne.ShortcutCopy{Clipboard: mainWindow.Clipboard()}
	copyItem := fyne.NewMenuItem("Copy", func() {
		shortcutFocused(copyShortcut, mainWindow)
	})
	copyItem.Shortcut = copyShortcut
	pasteShortcut := &fyne.ShortcutPaste{Clipboard: mainWindow.Clipboard()}
	pasteItem := fyne.NewMenuItem("Paste", func() {
		shortcutFocused(pasteShortcut, mainWindow)
	})
	pasteItem.Shortcut = pasteShortcut
	performFind := func() { fmt.Println("Menu Find") }
	findItem := fyne.NewMenuItem("Find", performFind)
	findItem.Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyF, Modifier: fyne.KeyModifierShortcutDefault | fyne.KeyModifierAlt | fyne.KeyModifierShift | fyne.KeyModifierControl | fyne.KeyModifierSuper}
	mainWindow.Canvas().AddShortcut(findItem.Shortcut, func(shortcut fyne.Shortcut) {
		performFind()
	})
	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("Documentation", func() {
			u, _ := url.Parse("https://developer.fyne.io")
			_ = fyne.CurrentApp().OpenURL(u)
		}),
		fyne.NewMenuItem("Support", func() {
			u, _ := url.Parse("https://fyne.io/support/")
			_ = fyne.CurrentApp().OpenURL(u)
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Sponsor", func() {
			u, _ := url.Parse("https://fyne.io/sponsor/")
			_ = fyne.CurrentApp().OpenURL(u)
		}))
	// a quit item will be appended to our first (File) menu
	file := fyne.NewMenu("File", newItem, checkedItem, disabledItem)
	device := fyne.CurrentDevice()
	if !device.IsMobile() && !device.IsBrowser() {
		file.Items = append(file.Items, fyne.NewMenuItemSeparator(), settingsItem)
	}
	return fyne.NewMainMenu(
		file,
		fyne.NewMenu("Edit", cutItem, copyItem, pasteItem, fyne.NewMenuItemSeparator(), findItem),
		helpMenu,
	)
}