package UI

import (
	"GDv2/app"
	"path/filepath"

	fynecanvas "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func NewSaveFileView(service *app.Service, window fynecanvas.Window, refreshTabs func()) fynecanvas.CanvasObject {
	vBox := container.NewVBox()

	exportButton := widget.NewButton("Export Save File", func() {
		if err := service.ExportSaveFile(); err != nil {
			dialog.ShowError(err, window)
			return
		}

		dialog.ShowInformation("Exported", "Save file exported successfully.", window)
	})

	restoreButton := widget.NewButton("Restore Save File", func() {
		backupFiles, err := service.GetBackupFiles()
		if err != nil {
			dialog.ShowError(err, window)
			return
		}

		if len(backupFiles) == 0 {
			dialog.ShowInformation("No Backups", "No backup files found.", window)
			return
		}

		backupOptions := make([]string, len(backupFiles)+1)
		backupOptions[0] = "Latest save backup"
		for i, file := range backupFiles {
			backupOptions[i+1] = filepath.Base(file)
		}

		selectedIndex := 0
		backupList := widget.NewList(
			func() int { return len(backupOptions) },
			func() fynecanvas.CanvasObject { return widget.NewLabel("") },
			func(id widget.ListItemID, object fynecanvas.CanvasObject) {
				object.(*widget.Label).SetText(backupOptions[id])
			},
		)
		backupList.OnSelected = func(id widget.ListItemID) {
			if id == 0 {
				selectedIndex = 0
				return
			}

			selectedIndex = int(id) - 1
		}
		backupList.Select(0)

		content := container.NewBorder(
			widget.NewLabel("Select a backup file to restore:"),
			nil,
			nil,
			nil,
			container.NewGridWrap(fynecanvas.NewSize(450, 300), backupList),
		)

		dialog.NewCustomConfirm(
			"Restore Save File",
			"Restore",
			"Cancel",
			content,
			func(confirmed bool) {
				if !confirmed {
					return
				}

				dialog.NewConfirm(
					"Confirm restore",
					"Restoring this backup will replace the current save file. Continue?",
					func(restoreConfirmed bool) {
						if !restoreConfirmed {
							return
						}

						if err := service.RestoreBackupFile(backupFiles[selectedIndex]); err != nil {
							dialog.ShowError(err, window)
							return
						}

						refreshTabs()
						dialog.ShowInformation(
							"Restored",
							"The save file was restored successfully.",
							window,
						)
					},
					window,
				).Show()
			},
			window,
		).Show()
	})

	vBox.Add(exportButton)
	vBox.Add(restoreButton)

	return vBox
}
