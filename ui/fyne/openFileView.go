package UI

import (
	"GDv2/app"
	"GDv2/utils"

	fynecanvas "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func NewOpenFileView(window fynecanvas.Window) fynecanvas.CanvasObject {
	title := widget.NewLabel("Open a Geareo save file")
	title.Alignment = fynecanvas.TextAlignCenter

	openButton := widget.NewButton("Open Save File", func() {
		path, err := utils.OpenSaveFile()
		if err != nil {
			dialog.ShowError(err, window)
			return
		}

		service, err := app.NewService(path)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}

		window.SetContent(NewMainView(service, window))
	})

	return container.NewCenter(
		container.NewVBox(
			title,
			openButton,
		),
	)
}