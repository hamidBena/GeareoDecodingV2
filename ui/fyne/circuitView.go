package UI

import (
	"GDv2/app"
	"GDv2/utils"
	customWidgets "GDv2/utils/widgets"
	"GDv2/world/model"
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type CircuitState struct {
	Circuits      []*model.Circuit
	SelectedIndex int
}

type circuitView struct {
	service *app.Service
	window  fyne.Window
	state   CircuitState
	parts   []model.CircuitEntity

	circuitList *widget.List
	partsList   *widget.List
}

func NewCircuitView(service *app.Service, window fyne.Window) fyne.CanvasObject {
	circuits, err := service.GetCircuits()
	if err != nil {
		return widget.NewLabel(err.Error())
	}

	view := newCircuitView(service, window, circuits)
	return view.layout()
}

func newCircuitView(service *app.Service, window fyne.Window, circuits []*model.Circuit) *circuitView {
	view := &circuitView{
		service: service,
		window:  window,
		state: CircuitState{
			Circuits:      circuits,
			SelectedIndex: -1,
		},
	}

	view.circuitList = widget.NewList(
		func() int { return len(view.state.Circuits) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, object fyne.CanvasObject) {
			object.(*widget.Label).SetText(view.state.Circuits[id].Name)
		},
	)
	view.circuitList.OnSelected = func(id widget.ListItemID) {
		view.state.SelectedIndex = int(id)
		view.refreshParts()
	}

	view.partsList = widget.NewList(
		func() int {
			if len(view.parts) == 0 {
				return 1
			}
			return len(view.parts)
		},
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, object fyne.CanvasObject) {
			label := object.(*widget.Label)
			if len(view.parts) == 0 {
				label.SetText("No parts available, select a circuit to view its parts.")
				return
			}

			label.SetText(fmt.Sprintf("%s #%d", view.parts[id].Key, view.parts[id].ID))
		},
	)

	return view
}

func (view *circuitView) layout() fyne.CanvasObject {
	displayIcon, err := utils.LoadIcon("builds/assets/display.png")
	if err != nil {
		return widget.NewLabel("Could not load display icon: " + err.Error())
	}

	romIcon, err := utils.LoadIcon("builds/assets/rom.png")
	if err != nil {
		return widget.NewLabel("Could not load ROM icon: " + err.Error())
	}

	operations := container.NewVBox(
		widget.NewLabel("Circuit operations"),
		widget.NewButton("Create Circuit", view.createCircuit),
		widget.NewButton("Embed Circuit", view.embedCircuit),
		widget.NewButton("Import Circuit", view.importCircuit),
		widget.NewButton("Export Circuit", view.exportCircuit),
		widget.NewButton("Delete Circuit", view.deleteCircuit),
		widget.NewButton("Rename Circuit", view.renameCircuit),
	)

	builders := container.NewGridWrap(fyne.NewSize(200, 250),
		customWidgets.ItemCard("RGB Display", "custom size RGB display", displayIcon, view.showDisplayBuilder),
		customWidgets.ItemCard("ROM", "read only memory bank", romIcon, view.showROMBuilder),
	)

	tabs := container.NewAppTabs(
		container.NewTabItem("Operations", operations),
		container.NewTabItem("Parts", view.partsList),
		container.NewTabItem("Auto Builders", builders),
	)
	tabs.SetTabLocation(container.TabLocationLeading)

	return container.NewHSplit(tabs, view.circuitList)
}

func (view *circuitView) refreshParts() {
	view.parts = nil
	index := view.state.SelectedIndex
	if index >= 0 && index < len(view.state.Circuits) {
		view.parts = view.state.Circuits[index].Elements.Entities
	}
	view.partsList.UnselectAll()
	view.partsList.Refresh()
}

func (view *circuitView) refreshCircuits() {
	circuits, err := view.service.GetCircuits()
	if err != nil {
		dialog.ShowError(err, view.window)
		return
	}
	view.state.Circuits = circuits
	view.state.SelectedIndex = -1
	view.circuitList.UnselectAll()
	view.circuitList.Refresh()
	view.refreshParts()
}

func (view *circuitView) selectedCircuit() (int, bool) {
	index := view.state.SelectedIndex
	if index < 0 || index >= len(view.state.Circuits) {
		dialog.ShowInformation("No circuit selected", "Select a circuit first.", view.window)
		return 0, false
	}
	return index, true
}

func (view *circuitView) createCircuit() {
	name := widget.NewEntry()
	name.SetPlaceHolder("Enter circuit name")
	dialog.NewForm("Create circuit", "Create", "Cancel", []*widget.FormItem{
		widget.NewFormItem("Name", name),
	}, func(confirmed bool) {
		if !confirmed {
			return
		}
		if err := view.service.CreateCircuit(name.Text); err != nil {
			dialog.ShowError(err, view.window)
			return
		}
		view.refreshCircuits()
	}, view.window).Show()
}

func (view *circuitView) embedCircuit() {
	index, ok := view.selectedCircuit()
	if !ok {
		return
	}
	x := newNumberEntry("0")
	y := newNumberEntry("0")
	dialog.NewForm("Embed Circuit", "Embed", "Cancel", []*widget.FormItem{
		widget.NewFormItem("Offset X", x),
		widget.NewFormItem("Offset Y", y),
	}, func(confirmed bool) {
		if !confirmed {
			return
		}
		offsetX, err := parseEntryInt(x, "Offset X")
		if err != nil {
			dialog.ShowInformation("Invalid offset", err.Error(), view.window)
			return
		}
		offsetY, err := parseEntryInt(y, "Offset Y")
		if err != nil {
			dialog.ShowInformation("Invalid offset", err.Error(), view.window)
			return
		}
		if err := view.service.EmbedCircuit(index, offsetX, offsetY); err != nil {
			dialog.ShowError(err, view.window)
			return
		}
		dialog.ShowInformation("Embedded", "Circuit embedded successfully.", view.window)
	}, view.window).Show()
}

func (view *circuitView) importCircuit() {
	if err := view.service.ImportCircuit(); err != nil {
		dialog.ShowError(err, view.window)
		return
	}
	view.refreshCircuits()
}

func (view *circuitView) exportCircuit() {
	index, ok := view.selectedCircuit()
	if !ok {
		return
	}
	if err := view.service.ExportCircuit(index); err != nil {
		dialog.ShowError(err, view.window)
		return
	}
	dialog.ShowInformation("Operation Complete", "The circuit has been exported successfully!", view.window)
}

func (view *circuitView) deleteCircuit() {
	index, ok := view.selectedCircuit()
	if !ok {
		return
	}
	dialog.NewConfirm("Delete circuit", "Are you sure you want to delete this circuit?", func(confirmed bool) {
		if !confirmed {
			return
		}
		if err := view.service.DeleteCircuit(index); err != nil {
			dialog.ShowError(err, view.window)
			return
		}
		view.refreshCircuits()
	}, view.window).Show()
}

func (view *circuitView) renameCircuit() {
	index, ok := view.selectedCircuit()
	if !ok {
		return
	}
	name := widget.NewEntry()
	name.SetPlaceHolder("Enter new circuit name")
	dialog.NewForm("Rename circuit", "Rename", "Cancel", []*widget.FormItem{
		widget.NewFormItem("Name", name),
	}, func(confirmed bool) {
		if !confirmed {
			return
		}
		newName := strings.TrimSpace(name.Text)
		if newName == "" {
			dialog.ShowInformation("Invalid name", "Circuit name cannot be empty.", view.window)
			return
		}
		if err := view.service.RenameCircuit(index, newName); err != nil {
			dialog.ShowError(err, view.window)
			return
		}
		view.refreshCircuits()
	}, view.window).Show()
}

func (view *circuitView) showDisplayBuilder() {
	index, ok := view.selectedCircuit()
	if !ok {
		return
	}
	positionX := newNumberEntry("0")
	positionY := newNumberEntry("0")
	width := newNumberEntry("8")
	height := newNumberEntry("8")
	dialog.NewForm("Build display circuit", "Build", "Cancel", []*widget.FormItem{
		widget.NewFormItem("Position X", positionX),
		widget.NewFormItem("Position Y", positionY),
		widget.NewFormItem("Width", width),
		widget.NewFormItem("Height", height),
	}, func(confirmed bool) {
		if !confirmed {
			return
		}
		values, err := parseEntries(
			[]*widget.Entry{positionX, positionY, width, height},
			[]string{"Position X", "Position Y", "Width", "Height"},
		)
		if err != nil {
			dialog.ShowInformation("Invalid parameter", err.Error(), view.window)
			return
		}
		if err := view.service.BuildDisplayCircuit(index, model.Position2D{X: values[0], Y: values[1]}, values[2], values[3]); err != nil {
			dialog.ShowError(err, view.window)
			return
		}
		dialog.ShowInformation("Build complete", "The display circuit was built successfully.", view.window)
	}, view.window).Show()
}

func (view *circuitView) showROMBuilder() {
	index, ok := view.selectedCircuit()
	if !ok {
		return
	}
	x := newNumberEntry("0")
	y := newNumberEntry("0")
	cells := newNumberEntry("16")
	layer := newNumberEntry("0")
	data := widget.NewEntry()
	data.SetPlaceHolder("Optional data file")
	browse := widget.NewButton("Browse", func() {
		path, err := utils.OpenCSV()
		if err == nil {
			data.SetText(path)
		}
	})
	dataPicker := container.NewBorder(nil, nil, nil, browse, data)
	dialog.NewForm("Build ROM", "Build", "Cancel", []*widget.FormItem{
		widget.NewFormItem("Position X", x),
		widget.NewFormItem("Position Y", y),
		widget.NewFormItem("Cell count", cells),
		widget.NewFormItem("Layer", layer),
		widget.NewFormItem("Data file", dataPicker),
	}, func(confirmed bool) {
		if !confirmed {
			return
		}
		values, err := parseEntries([]*widget.Entry{x, y, cells, layer}, []string{"X position", "Y position", "Cell count", "Layer"})
		if err != nil {
			dialog.ShowError(err, view.window)
			return
		}
		if err := view.service.BuildROMCircuit(index, model.Position2D{X: values[0], Y: values[1]}, values[2], values[3], strings.TrimSpace(data.Text)); err != nil {
			dialog.ShowError(err, view.window)
			return
		}
		dialog.ShowInformation("Build complete", "ROM built successfully.", view.window)
	}, view.window).Show()
}

func newNumberEntry(value string) *widget.Entry {
	entry := widget.NewEntry()
	entry.SetText(value)
	return entry
}

func parseEntryInt(entry *widget.Entry, name string) (int, error) {
	value, err := strconv.Atoi(strings.TrimSpace(entry.Text))
	if err != nil {
		return 0, fmt.Errorf("%s must be a whole number", name)
	}
	return value, nil
}

func parseEntries(entries []*widget.Entry, names []string) ([]int, error) {
	values := make([]int, len(entries))
	for i, entry := range entries {
		value, err := parseEntryInt(entry, names[i])
		if err != nil {
			return nil, err
		}
		values[i] = value
	}
	return values, nil
}
