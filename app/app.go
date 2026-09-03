package app

import (
	"GDv2/utils"
)

func Run(ui UI) {
	for {
		ui.ShowMessage("Please Select a save file to edit")
		path, err := utils.OpenSaveFile()
		if err != nil {
			ui.ShowError(err)
			continue
		}

		worldSvc, err := NewService(path)
		if err != nil {
			ui.ShowError(err)
			continue
		}
		controller := NewController(worldSvc, ui)
		ok := controller.Run()
		if !ok {
			ui.ShowMessage("Exiting...")
			break
		}
	}
}
