package app

import (
	"GDv2/utils"
	"fmt"
	"path/filepath"
	"strconv"
)

// circuit management Menu
func (c *Controller) showCircuits() {
	c.ui.ClearScreen()
	circuits, err := c.worldSvc.GetCircuits()
	if err != nil {
		c.ui.ShowMessage(fmt.Sprintf("Error: %v", err))
		return
	}
	for i, circuit := range circuits {
		c.ui.ShowMessage(fmt.Sprintf("%d. %s", i+1, circuit.Name))
	}
}

func (c *Controller) selectCircuit() int {
	circuits, err := c.worldSvc.GetCircuits()
	if err != nil {
		c.ui.ShowMessage(fmt.Sprintf("Error: %v", err))
		return -1
	}
	menuOptions := make([]MenuOption, len(circuits))
	for i, circuit := range circuits {
		menuOptions[i] = MenuOption{
			Key:   fmt.Sprintf("%d", i+1),
			Label: circuit.Name,
		}
	}
	menuOptions = append(menuOptions, MenuOption{Key: "0", Label: "Cancel"})
	choice := c.ui.ShowMenu("Select Circuit", menuOptions)
	if choice == "0" {
		return -1
	}

	for i := range circuits {
		if choice == fmt.Sprintf("%d", i+1) {
			return i
		}
	}
	return -1
}

func (c *Controller) deleteCircuit() {
	c.ui.ClearScreen()
	c.ui.ShowMessage("Please select a circuit to delete:")
	selectedIndex := c.selectCircuit()
	if selectedIndex != -1 {
		c.worldSvc.DeleteCircuit(selectedIndex)
	}
}

func (c *Controller) exportCircuit() {
	c.ui.ClearScreen()
	c.ui.ShowMessage("Please select a circuit to export:")
	selectedIndex := c.selectCircuit()

	c.ui.ShowMessage("Please choose where you want to save the exported circuit")
	path, err := utils.SaveOutputFile()
	if err != nil {
		c.ui.ShowMessage(fmt.Sprintf("Error: %v", err))
		return
	}

	err = c.worldSvc.ExportCircuit(selectedIndex, path)
	if err != nil {
		c.ui.ShowMessage(fmt.Sprintf("Error exporting circuit: %v", err))
	} else {
		c.ui.ShowMessage("Circuit exported successfully!")
	}

	c.ui.ClickAny()
}

func (c *Controller) importCircuit() {
	c.ui.ClearScreen()
	c.ui.ShowMessage("Please choose the circuit files to import")
	paths, err := utils.OpenSaveFiles()
	if err != nil {
		c.ui.ShowMessage(fmt.Sprintf("Error: %v", err))
		return
	}

	err = c.worldSvc.ImportCircuit(paths)
	if err != nil {
		c.ui.ShowMessage(fmt.Sprintf("Error importing circuits: %v", err))
	} else {
		c.ui.ShowMessage("Circuit imported successfully!")
	}

	c.ui.ClickAny()
}

// main menu options
func (c *Controller) saveChanges() {
	c.validateSavefile()

	ok := c.ui.Confirm("Are you sure you want to proceed?")
	if !ok {
		return
	}

	err := c.worldSvc.SaveFile()
	if err != nil {
		fmt.Println("error saving file:", err)
	}
}

func (c *Controller) exit() {
	ok := c.ui.Confirm("Do you want to save changes before exiting?")
	if !ok {
		return
	}

	c.saveChanges()
}

// save file management
func (c *Controller) restoreSaveFile() {
	c.ui.ClearScreen()

	backupFiles, err := c.worldSvc.GetBackupFiles()
	if err != nil {
		c.ui.ShowMessage(fmt.Sprintf("Error loading backups: %v", err))
		c.ui.ClickAny()
		return
	}

	if len(backupFiles) == 0 {
		c.ui.ShowMessage("No backup files found")
		c.ui.ClickAny()
		return
	}

	menuOptions := make([]MenuOption, len(backupFiles))
	for i := range backupFiles {
		menuOptions[i] = MenuOption{
			Key:   fmt.Sprintf("%d", i+1),
			Label: filepath.Base(backupFiles[i]),
		}
	}
	menuOptions = append(menuOptions, MenuOption{Key: "0", Label: "Cancel"})

	choice := c.ui.ShowMenu("Backup files found:", menuOptions)
	if choice == "0" {
		return
	}

	index, err := strconv.Atoi(choice)
	if err != nil || index < 1 || index > len(backupFiles) {
		c.ui.ShowMessage("Invalid selection")
		c.ui.ClickAny()
		return
	}

	if err := c.worldSvc.RestoreBackupFile(backupFiles[index-1]); err != nil {
		c.ui.ShowMessage(fmt.Sprintf("Error restoring backup file: %v", err))
	} else {
		c.ui.ShowMessage("Backup file restored successfully!")
	}

	c.ui.ClickAny()
}

func (c *Controller) validateSavefile() {
	c.ui.ClearScreen()
	issues, err := c.worldSvc.ValidateSavefile()
	if err != nil {
		c.ui.ShowMessage(fmt.Sprintf("Error validating save file: %v", err))
		c.ui.ClickAny()
		return
	}

	if len(issues) == 0 {
		c.ui.ShowMessage("No issues found in the save file.")
	} else {
		c.ui.ShowMessage("Issues found in the save file:")
		for _, issue := range issues {
			c.ui.ShowMessage(fmt.Sprintf("%s: %s", issue.Severity, issue.Message))
		}
	}
}
