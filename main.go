package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/cmd/fyne_demo/tutorials"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
)

var mainWindow fyne.Window
var mainMenu *fyne.Menu

func main() {
	a := app.NewWithID("Verify.DiffTrayApp")
	a.SetIcon(resourceCogsPng)

	createMainMenu()
	createTrayIcon()

	mainWindow = a.NewWindow("Main")
	mainWindow.Resize(fyne.NewSize(640, 460))
	mainWindow.SetCloseIntercept(func() {
		//prevent main window from closing
		mainWindow.Hide()
	})

	//content := container.NewMax()
	//title := widget.NewLabel("Component name")
	//intro := widget.NewLabel("An introduction would probably go\nhere, as well as a")
	//intro.Wrapping = fyne.TextWrapWord
	//setTutorial := func(t tutorials.Tutorial) {
	//	if fyne.CurrentDevice().IsMobile() {
	//		child := a.NewWindow(t.Title)
	//		topWindow = child
	//		child.SetContent(t.View(topWindow))
	//		child.Show()
	//		child.SetOnClosed(func() {
	//			topWindow = w
	//		})
	//		return
	//	}
	//	title.SetText(t.Title)
	//	intro.SetText(t.Intro)
	//	content.Objects = []fyne.CanvasObject{t.View(w)}
	//	content.Refresh()
	//}
	//tutorial := container.NewBorder(
	//	container.NewVBox(title, widget.NewSeparator(), intro), nil, nil, nil, content)
	//if fyne.CurrentDevice().IsMobile() {
	//	w.SetContent(makeNav(setTutorial, false))
	//} else {
	//	split := container.NewHSplit(makeNav(setTutorial, true), tutorial)
	//	split.Offset = 0.2
	//	w.SetContent(split)
	//}
	//w.Resize(fyne.NewSize(640, 460))
	//w.ShowAndRun()
	a.Run()
}

func logLifecycle(a fyne.App) {
	a.Lifecycle().SetOnStarted(func() {
		log.Println("Lifecycle: Started")
	})
	a.Lifecycle().SetOnStopped(func() {
		log.Println("Lifecycle: Stopped")
	})
	a.Lifecycle().SetOnEnteredForeground(func() {
		log.Println("Lifecycle: Entered Foreground")
	})
	a.Lifecycle().SetOnExitedForeground(func() {
		log.Println("Lifecycle: Exited Foreground")
	})
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
	deskApp := fyne.CurrentApp().(desktop.App)
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