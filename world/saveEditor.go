package world

import (
	"GDv2/world/model"
	"encoding/json"
	"fmt"
	"os"
)

type SaveEditor struct {
	SaveFile model.SaveFile
	IsSaved  bool
}

func NewSaveEditor(path string) (*SaveEditor, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read save file: %w", err)
	}

	var sf model.SaveFile
	if err := json.Unmarshal(raw, &sf); err != nil {
		return nil, fmt.Errorf("parse save file: %w", err)
	}
	sf.Path = path

	return &SaveEditor{SaveFile: sf}, nil
}

func (se *SaveEditor) ImportCircuit(circuit *model.Circuit) error {
	se.SaveFile.Hub.Circuits = append(se.SaveFile.Hub.Circuits, *circuit)
	return nil
}

func (se *SaveEditor) ReloadSaveFile(path string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read save file: %w", err)
	}

	var sf model.SaveFile
	if err := json.Unmarshal(raw, &sf); err != nil {
		return fmt.Errorf("parse save file: %w", err)
	}
	sf.Path = path
	se.SaveFile = sf
	return nil
}
