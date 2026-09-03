package model

import (
	"encoding/json"
	"fmt"
)

// Faceplate tile entity keys
const (
	FPKeyFill  = "fill"
	FPKeyLight = "light"
	FPKeyPin   = "pin"
)

type Faceplate struct {
	Tiles    Tile          `json:"tiles"` // was []Tile — the file has a single object here, not an array
	Camera   CircuitCamera `json:"camera"`
	Selector Selector      `json:"selector"`
}

type Tile struct {
	LastID   int        `json:"lastId"`
	Entities []FPEntity `json:"entities"`
}

type FPEntity struct {
	Key  string       `json:"key"`
	ID   int          `json:"id"`
	Data FPEntityData `json:"data"`
}

type FPEntityData interface {
	isFPEntityData()
}

type FPEntityDataBase struct {
	Position Position2D `json:"position"`
	Rotation Direction  `json:"rotation"`
	Name     string     `json:"name"`
	Label    string     `json:"label"`
}

// entities and their data bodies, and the isFPEntityData method to satisfy the interface

type FillData struct {
	FPEntityDataBase
	Color string `json:"color"`
}

type FPLightData struct {
	FPEntityDataBase
	Color string `json:"color"`
}

type PinData struct {
	FPEntityDataBase
}

func (FillData) isFPEntityData()    {}
func (FPLightData) isFPEntityData() {}
func (PinData) isFPEntityData()     {}

// --- JSON handling ---

func (e *FPEntity) UnmarshalJSON(raw []byte) error {
	var shell struct {
		Key  string          `json:"key"`
		ID   int             `json:"id"`
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(raw, &shell); err != nil {
		return fmt.Errorf("unmarshal faceplate entity shell: %w", err)
	}
	e.Key = shell.Key
	e.ID = shell.ID

	var data FPEntityData
	switch shell.Key {
	case FPKeyFill:
		var d FillData
		if err := json.Unmarshal(shell.Data, &d); err != nil {
			return fmt.Errorf("unmarshal %s data: %w", shell.Key, err)
		}
		data = d
	case FPKeyLight:
		var d FPLightData
		if err := json.Unmarshal(shell.Data, &d); err != nil {
			return fmt.Errorf("unmarshal %s data: %w", shell.Key, err)
		}
		data = d
	case FPKeyPin:
		var d PinData
		if err := json.Unmarshal(shell.Data, &d); err != nil {
			return fmt.Errorf("unmarshal %s data: %w", shell.Key, err)
		}
		data = d
	default:
		return fmt.Errorf("unknown faceplate entity key: %q", shell.Key)
	}
	e.Data = data
	return nil
}

func (e FPEntity) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Key  string       `json:"key"`
		ID   int          `json:"id"`
		Data FPEntityData `json:"data"`
	}{e.Key, e.ID, e.Data})
}
