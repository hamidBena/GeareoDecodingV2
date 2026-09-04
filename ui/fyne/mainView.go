package UI

import (
	"GDv2/app"

	fynecanvas "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func NewMainView(service *app.Service, window fynecanvas.Window) fynecanvas.CanvasObject {
	tabsContent := container.NewStack()
	var refreshTabs func()

	refreshTabs = func() {
		tabs := container.NewDocTabs(
			container.NewTabItem("Circuits", NewCircuitView(service, window)),
			container.NewTabItem("File", NewSaveFileView(service, window, refreshTabs)),
			// container.NewTabItem("Validation", NewValidationView(service, window)),
		)

		tabsContent.Objects = []fynecanvas.CanvasObject{tabs}
		tabsContent.Refresh()
	}

	saveButton := widget.NewButton("Save", func() {
		if err := service.SaveFile(); err != nil {
			dialog.ShowError(err, window)
			return
		}

		dialog.ShowInformation("Saved", "Save file written successfully.", window)
	})

	exitButton := widget.NewButton("Exit", func() {
		dialog.NewConfirm(
			"Confirm save and exit",
			"Do you want to save changes before exiting?",
			func(confirmed bool) {
				if confirmed {
					if err := service.SaveFile(); err != nil {
						dialog.NewConfirm(
							"Save failed",
							"Could not save the file. Exit without saving?",
							func(exitAnyway bool) {
								if exitAnyway {
									window.Close()
								}
							},
							window,
						).Show()

						return
					}
				}

				window.Close()
			},
			window,
		).Show()
	})

	reloadButton := widget.NewButton("Reload file", func() {
		dialog.NewConfirm(
			"Confirm reload",
			"Do you want to reload the file? All unsaved changes will be lost.",
			func(confirmed bool) {
				if !confirmed {
					return
				}

				if err := service.ReloadFile(); err != nil {
					dialog.ShowError(err, window)
					return
				}

				refreshTabs()
				dialog.ShowInformation(
					"Reloaded",
					"Save file reloaded successfully.",
					window,
				)
			},
			window,
		).Show()
	})

	toolbar := container.NewHBox(saveButton, exitButton, reloadButton)
	refreshTabs()

	return container.NewBorder(toolbar, nil, nil, nil, tabsContent)
}
