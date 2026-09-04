package UI

import (
	"GDv2/app"
	"GDv2/world/model"

	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type CircuitState struct {
	Circuits      []*model.Circuit
	SelectedIndex int
}

func NewCircuitView(service *app.Service, window fyne.Window) fyne.CanvasObject {
	circuits, err := service.GetCircuits()
	if err != nil {
		return widget.NewLabel(err.Error())
	}

	selectedIndex := -1

	circuitList := widget.NewList(
		func() int {
			return len(circuits)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			object.(*widget.Label).SetText(circuits[id].Name)
		},
	)

	parts := []model.CircuitEntity{}
	partsList := widget.NewList(
		func() int {
			return len(parts)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			entity := parts[id]
			object.(*widget.Label).SetText(fmt.Sprintf("%s #%d", entity.Key, entity.ID))
		},
	)

	refreshParts := func() {
		parts = nil
		if selectedIndex >= 0 && selectedIndex < len(circuits) {
			parts = circuits[selectedIndex].Elements.Entities
		}
		partsList.UnselectAll()
		partsList.Refresh()
	}

	circuitList.OnSelected = func(id widget.ListItemID) {
		selectedIndex = int(id)
		refreshParts()
	}

	refreshCircuits := func() {
		updatedCircuits, err := service.GetCircuits()
		if err != nil {
			dialog.ShowError(err, window)
			return
		}

		circuits = updatedCircuits
		selectedIndex = -1
		circuitList.UnselectAll()
		circuitList.Refresh()
		refreshParts()
	}

	createButton := widget.NewButton("Create Circuit", func() {
		entry := widget.NewEntry()
		entry.SetPlaceHolder("Enter circuit name")

		dialog.NewForm(
			"Create circuit",
			"Create",
			"Cancel",
			[]*widget.FormItem{
				widget.NewFormItem("Name", entry),
			},
			func(confirmed bool) {
				if !confirmed {
					return
				}

				if err := service.CreateCircuit(entry.Text); err != nil {
					dialog.ShowError(err, window)
					return
				}

				refreshCircuits()
			},
			window,
		).Show()
	})

	embedButton := widget.NewButton("Embed Circuit", func() {
		if selectedIndex == -1 {
			dialog.ShowInformation(
				"No circuit selected",
				"Select a destination circuit first.",
				window,
			)
			return
		}

		offsetXEntry := widget.NewEntry()
		offsetXEntry.SetText("0")

		offsetYEntry := widget.NewEntry()
		offsetYEntry.SetText("0")

		dialog.NewForm(
			"Embed Circuit",
			"Embed",
			"Cancel",
			[]*widget.FormItem{
				widget.NewFormItem("Offset X", offsetXEntry),
				widget.NewFormItem("Offset Y", offsetYEntry),
			},
			func(confirmed bool) {
				if !confirmed {
					return
				}

				offsetX, err := strconv.Atoi(offsetXEntry.Text)
				if err != nil {
					dialog.ShowInformation(
						"Invalid offset",
						"Offset X must be a whole number.",
						window,
					)
					return
				}

				offsetY, err := strconv.Atoi(offsetYEntry.Text)
				if err != nil {
					dialog.ShowInformation(
						"Invalid offset",
						"Offset Y must be a whole number.",
						window,
					)
					return
				}

				if err := service.EmbedCircuit(selectedIndex, offsetX, offsetY); err != nil {
					dialog.ShowError(err, window)
					return
				}

				dialog.ShowInformation(
					"Embedded",
					"Circuit embedded successfully.",
					window,
				)
			},
			window,
		).Show()
	})

	importButton := widget.NewButton("Import Circuit", func() {
		err := service.ImportCircuit()
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		refreshCircuits()
	})

	exportButton := widget.NewButton("Export Circuit", func() {
		if selectedIndex == -1 {
			dialog.ShowInformation(
				"No circuit selected",
				"Select a circuit first.",
				window,
			)
			return
		}

		err := service.ExportCircuit(selectedIndex)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
	})

	deleteButton := widget.NewButton("Delete Circuit", func() {
		if selectedIndex == -1 {
			dialog.ShowInformation(
				"No circuit selected",
				"Select a circuit first.",
				window,
			)
			return
		}
		dialog.NewConfirm(
			"Delete circuit",
			"Are you sure you want to delete this circuit?",
			func(confirmed bool) {
				if !confirmed {
					return
				}

				if err := service.DeleteCircuit(selectedIndex); err != nil {
					dialog.ShowError(err, window)
					return
				}

				refreshCircuits()
			},
			window,
		).Show()
	})

	renameButton := widget.NewButton("Rename Circuit", func() {
		if selectedIndex == -1 {
			dialog.ShowInformation(
				"No circuit selected",
				"Select a circuit first.",
				window,
			)
			return
		}

		entry := widget.NewEntry()
		entry.SetPlaceHolder("Enter new circuit name")

		formentry := widget.NewFormItem("name", entry)

		dialog.NewForm(
			"Rename circuit",
			"rename",
			"cancel",
			[]*widget.FormItem{formentry},
			func(b bool) {
				if !b {
					return
				}

				if entry.Text == "" {
					dialog.ShowInformation(
						"Invalid name",
						"Circuit name cannot be empty.",
						window,
					)
					return
				}

				if err := service.RenameCircuit(selectedIndex, entry.Text); err != nil {
					dialog.ShowError(err, window)
					return
				}

				refreshCircuits()
			}, window).Show()
	})

	operations := container.NewVBox(
		widget.NewLabel("Circuit operations"),
		createButton,
		embedButton,
		importButton,
		exportButton,
		deleteButton,
		renameButton,
	)

	rightTabs := container.NewAppTabs(
		container.NewTabItem("Operations", operations),
		container.NewTabItem("Parts", partsList),
	)
	rightTabs.SetTabLocation(container.TabLocationLeading)

	return container.NewHSplit(
		rightTabs,
		circuitList,
	)
}
