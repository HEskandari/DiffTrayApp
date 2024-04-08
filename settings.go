package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"testing"
)

type settingsDialog struct {
	title    string
	confirm  string
	dismiss  string
	widgets  map[string]fyne.CanvasObject
	callback func(bool)
	invalid  bool
}

func (d *settingsDialog) Show() {}

func (d *settingsDialog) Hide() {}

func (d *settingsDialog) SetDismissText(string) {}

func (d *settingsDialog) SetOnClosed(func()) {}

func (d *settingsDialog) Refresh() {}

func (d *settingsDialog) Resize(fyne.Size) {}

func (d *settingsDialog) MinSize() fyne.Size {
	return fyne.Size{
		Width:  320,
		Height: 200,
	}
}

func (d *settingsDialog) getTitle() string {
	return d.title
}

func (d *settingsDialog) getConfirm() string {
	return d.confirm
}

func (d *settingsDialog) getDismiss() string {
	return d.dismiss
}

func (d *settingsDialog) isValid() bool {
	return !d.invalid
}

func (d *settingsDialog) tapOk() {
	d.invalid = false
	for _, wi := range d.widgets {
		if w, ok := wi.(fyne.Validatable); ok {
			if e := w.Validate(); e != nil {
				d.invalid = true
				break
			}
		}
	}
	if !d.invalid {
		d.callback(true)
	}
}

func (d *settingsDialog) tapCancel() {
	d.callback(false)
}

func (d *settingsDialog) setText(t *testing.T, name string, text string) {
	wi, ok := d.widgets[name]
	if !ok {
		t.Fail()
		return
	}
	e, ok := wi.(*widget.Entry)
	if !ok {
		t.Fail()
		return
	}
	e.SetText(text)
}

func (d *settingsDialog) setCheck(t *testing.T, name string, check bool) {
	wi, ok := d.widgets[name]
	if !ok {
		t.Fail()
		return
	}
	c, ok := wi.(*widget.Check)
	if !ok {
		t.Fail()
		return
	}
	c.Checked = check
}

func (d *settingsDialog) tapButton(t *testing.T, name string) {
	t.Helper()
	wi, ok := d.widgets[name]
	if !ok {
		t.Errorf("there's no widget with name %s", name)
		return
	}
	b, ok := wi.(*widget.Button)
	if !ok {
		t.Errorf("widget '%s' isn't a button", name)
		return
	}
	b.OnTapped()
}

var lastTestDialog *settingsDialog = nil

func NewOptionsForm(title, confirm, dismiss string, items []*widget.FormItem, callback func(bool), _ fyne.Window) dialog.Dialog {
	widgets := make(map[string]fyne.CanvasObject)
	for _, i := range items {
		widgetsForItem := digWidgets(i.Widget)
		l := len(widgetsForItem)
		if l < 1 {
			continue
		}
		if l == 1 {
			widgets[i.Text] = widgetsForItem[0]
			continue
		}
		for x, wi := range widgetsForItem {
			widgets[fmt.Sprintf("%s-%d", i.Text, x)] = wi
		}
	}
	lastTestDialog = &settingsDialog{title: title, confirm: confirm, dismiss: dismiss, widgets: widgets, callback: callback}
	return lastTestDialog
}

func digWidgets(root fyne.CanvasObject) []fyne.CanvasObject {
	if cnt, ok := root.(*fyne.Container); ok {
		var widgets []fyne.CanvasObject
		for _, o := range cnt.Objects {
			widgets = append(widgets, digWidgets(o)...)
		}
		return widgets
	}
	return []fyne.CanvasObject{root}
}
