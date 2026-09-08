package UI

import (
	"GDv2/app"
	"path/filepath"

	fynecanvas "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type saveFileView struct {
	service     *app.Service
	window      fynecanvas.Window
	refreshTabs func()
}

func NewSaveFileView(service *app.Service, window fynecanvas.Window, refreshTabs func()) fynecanvas.CanvasObject {
	view := &saveFileView{service: service, window: window, refreshTabs: refreshTabs}
	return view.layout()
}

func (view *saveFileView) layout() fynecanvas.CanvasObject {
	return container.NewVBox(
		widget.NewButton("Export Save File", view.export),
		widget.NewButton("Restore Save File", view.restore),
	)
}

func (view *saveFileView) export() {
	if err := view.service.ExportSaveFile(); err != nil {
		dialog.ShowError(err, view.window)
		return
	}
	dialog.ShowInformation("Exported", "Save file exported successfully.", view.window)
}

func (view *saveFileView) restore() {
	backupFiles, err := view.service.GetBackupFiles()
	if err != nil {
		dialog.ShowError(err, view.window)
		return
	}
	if len(backupFiles) == 0 {
		dialog.ShowInformation("No Backups", "No backup files found.", view.window)
		return
	}

	selectedIndex := 0
	backupOptions := make([]string, len(backupFiles)+1)
	backupOptions[0] = "Latest save backup"
	for i, file := range backupFiles {
		backupOptions[i+1] = filepath.Base(file)
	}
	backupList := widget.NewList(
		func() int { return len(backupOptions) },
		func() fynecanvas.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, object fynecanvas.CanvasObject) {
			object.(*widget.Label).SetText(backupOptions[id])
		},
	)
	backupList.OnSelected = func(id widget.ListItemID) {
		selectedIndex = int(id)
	}
	backupList.Select(0)

	content := container.NewBorder(
		widget.NewLabel("Select a backup file to restore:"),
		nil, nil, nil,
		container.NewGridWrap(fynecanvas.NewSize(450, 300), backupList),
	)
	dialog.NewCustomConfirm("Restore Save File", "Restore", "Cancel", content, func(confirmed bool) {
		if !confirmed {
			return
		}
		view.confirmRestore(backupFiles[selectedIndex])
	}, view.window).Show()
}

func (view *saveFileView) confirmRestore(path string) {
	dialog.NewConfirm(
		"Confirm restore",
		"Restoring this backup will replace the current save file. Continue?",
		func(confirmed bool) {
			if !confirmed {
				return
			}
			if err := view.service.RestoreBackupFile(path); err != nil {
				dialog.ShowError(err, view.window)
				return
			}
			view.refreshTabs()
			dialog.ShowInformation("Restored", "The save file was restored successfully.", view.window)
		},
		view.window,
	).Show()
}
