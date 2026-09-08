package UI

import (
	"GDv2/app"

	fynecanvas "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type mainView struct {
	service     *app.Service
	window      fynecanvas.Window
	tabsContent *fynecanvas.Container
}

func NewMainView(service *app.Service, window fynecanvas.Window) fynecanvas.CanvasObject {
	view := &mainView{
		service:     service,
		window:      window,
		tabsContent: container.NewStack(),
	}
	return view.layout()
}

func (view *mainView) layout() fynecanvas.CanvasObject {
	toolbar := container.NewHBox(
		widget.NewButton("Save", view.save),
		widget.NewButton("Exit", view.exit),
		widget.NewButton("Reload file", view.reload),
	)
	view.refreshTabs()
	return container.NewBorder(toolbar, nil, nil, nil, view.tabsContent)
}

func (view *mainView) refreshTabs() {
	tabs := container.NewDocTabs(
		container.NewTabItem("Circuits", NewCircuitView(view.service, view.window)),
		container.NewTabItem("File", NewSaveFileView(view.service, view.window, view.refreshTabs)),
		// container.NewTabItem("Validation", NewValidationView(service, window)),
	)

	view.tabsContent.Objects = []fynecanvas.CanvasObject{tabs}
	view.tabsContent.Refresh()
}

func (view *mainView) save() {
	if err := view.service.SaveFile(); err != nil {
		dialog.ShowError(err, view.window)
		return
	}
	dialog.ShowInformation("Saved", "Save file written successfully.", view.window)
}

func (view *mainView) exit() {
	dialog.NewConfirm(
		"Confirm save and exit",
		"Do you want to save changes before exiting?",
		func(confirmed bool) {
			if confirmed {
				if err := view.service.SaveFile(); err != nil {
					dialog.NewConfirm(
						"Save failed",
						"Could not save the file. Exit without saving?",
						func(exitAnyway bool) {
							if exitAnyway {
								view.window.Close()
							}
						},
						view.window,
					).Show()

					return
				}
			}

			view.window.Close()
		},
		view.window,
	).Show()
}

func (view *mainView) reload() {
	dialog.NewConfirm(
		"Confirm reload",
		"Do you want to reload the file? All unsaved changes will be lost.",
		func(confirmed bool) {
			if !confirmed {
				return
			}

			if err := view.service.ReloadFile(); err != nil {
				dialog.ShowError(err, view.window)
				return
			}

			view.refreshTabs()
			dialog.ShowInformation(
				"Reloaded",
				"Save file reloaded successfully.",
				view.window,
			)
		},
		view.window,
	).Show()
}
