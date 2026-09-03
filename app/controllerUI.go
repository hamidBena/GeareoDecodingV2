package app

type Controller struct {
	worldSvc *Service

	ui UI // interface, not a concrete CLI/GUI type

	// navigation/session state
	currentMenu MenuID
	history     []MenuID // for "back" if you don't want recursion-based back
}

func NewController(worldSvc *Service, ui UI) *Controller {
	return &Controller{
		worldSvc:    worldSvc,
		ui:          ui,
		currentMenu: MenuMain,
	}
}

func (c *Controller) Run() bool {
	return c.showMenuLoop(c.currentMenu)
}

func (c *Controller) showMenuLoop(id MenuID) bool {
	def := menus[id]
	for {
		c.ui.ClearScreen()
		choice := c.ui.ShowMenu(def.Title, def.Options)

		switch id {
		case MenuMain:
			switch choice {
			case "1":
				c.showMenuLoop(MenuSaveFile)
			case "2":
				c.showMenuLoop(MenuCircuits)
			case "8":
				c.saveChanges()
			case "9":
				return true // reload save file
			case "0":
				c.exit()
				return false
			}
		case MenuSaveFile:
			switch choice {
			case "1":
				// Implement export save file logic
			case "2":
				c.restoreSaveFile()
			case "3":
				c.validateSavefile()
				c.ui.ClickAny()
			case "0":
				return false
			}
		case MenuCircuits:
			switch choice {
			case "1":
				c.showCircuits()
				c.ui.ClickAny()
			case "2":
				c.deleteCircuit()
			case "8":
				c.importCircuit()
			case "9":
				c.exportCircuit()
			case "0":
				return false
			}
		case MenuCircuitManager:
			switch choice {
			case "1":
				// Implement export circuit logic
			case "2":
				// Implement rename circuit logic
			case "3":
				// Implement delete circuit logic
			case "0":
				return false
			}
		}
	}
}
