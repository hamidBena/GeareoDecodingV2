package model

import "encoding/json"

type SaveFile struct {
	Version     string          `json:"version"`
	Parts       json.RawMessage `json:"parts"`
	WorldCamera WorldCamera     `json:"camera"` // was "worldCamera" — didn't match the file, so it always unmarshaled to zero values
	Settings    Settings        `json:"settings"`
	Hub         Hub             `json:"hub"`
	Selector    Selector        `json:"selector"`
	Path        string          `json:"-"` // excluded from JSON — was missing a tag, so it defaulted to marshaling as "Path"
}

type WorldCamera struct {
	Px   float64 `json:"px"`
	Py   float64 `json:"py"`
	Pz   float64 `json:"pz"`
	Rx   float64 `json:"rx"`
	Ry   float64 `json:"ry"`
	Rz   float64 `json:"rz"`
	Zoom float64 `json:"zoom"`
}

type Settings struct {
	Data SettingsData `json:"data"`
}

type SettingsData struct {
	Advanced               bool    `json:"advanced"`
	Scene                  string  `json:"scene"`
	CollisionDetectionMode int     `json:"collisionDetectionMode"`
	PhysicsSteps           int     `json:"physicsSteps"`
	DynamicPhysicsStep     bool    `json:"dynamicPhysicsStep"`
	Limitless              bool    `json:"limitless"`
	HandleGrid             float64 `json:"handleGrid"`
	MoveGrid               float64 `json:"moveGrid"`
	RotateGrid             float64 `json:"rotateGrid"`
}

type Selector struct {
	Items []SelectorItem `json:"items"`
}

type SelectorItem struct {
	Value *int `json:"value,omitempty"`
}
