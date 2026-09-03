package model

import (
	"encoding/json"
	"fmt"
)

// metadata

// direction
type Direction int

const (
	Right Direction = iota
	Up
	Left
	Down
)

// Key constants for all supported circuit blocks
const (
	KeyToggle   = "toggle"
	KeyButton   = "button"
	KeyFunc     = "func"
	KeyLight    = "light"
	KeyCurve    = "curve"
	KeyPort     = "port"
	KeyBus      = "bus"
	KeyCircuit  = "circuit"
	KeyValve    = "valve"
	KeyTime     = "time"
	KeyClock    = "clock"
	KeyConstant = "constant"
	KeyDisplay  = "display"
	KeyVia      = "via"
	KeyWire     = "wire"
)

// Function type keys for all supported function blocks
type FunctionType string

const (
	// Logic
	FuncAnd  FunctionType = "and"
	FuncOr   FunctionType = "or"
	FuncNot  FunctionType = "not"
	FuncNand FunctionType = "nand"
	FuncNor  FunctionType = "nor"
	FuncXor  FunctionType = "xor"
	FuncXnor FunctionType = "xnor"

	// Arithmetic
	FuncSum      FunctionType = "sum"
	FuncSubtract FunctionType = "sub"
	FuncMultiply FunctionType = "mul"
	FuncDivide   FunctionType = "div"
	FuncNegate   FunctionType = "negate"
	FuncPass     FunctionType = "pass"
	FuncAverage  FunctionType = "avg"
	FuncMax      FunctionType = "max"
	FuncMin      FunctionType = "min"
	FuncAbs      FunctionType = "abs"
	FuncPower    FunctionType = "pow"
	FuncModulo   FunctionType = "mod"
	FuncSqrt     FunctionType = "sqrt"
	FuncRound    FunctionType = "round"
	FuncFloor    FunctionType = "floor"
	FuncCeil     FunctionType = "ceil"

	// Comparison
	FuncEqual        FunctionType = "equal"
	FuncNotEqual     FunctionType = "not equal"
	FuncGreater      FunctionType = "greater"
	FuncGreaterEqual FunctionType = "greater equal"
	FuncLess         FunctionType = "less"
	FuncLessEqual    FunctionType = "less equal"
	FuncSign         FunctionType = "sign"

	// Trigonometry
	FuncSin      FunctionType = "sin"
	FuncCos      FunctionType = "cos"
	FuncAtan     FunctionType = "atan"
	FuncRadToDeg FunctionType = "rad2deg"
	FuncDegToRad FunctionType = "deg2rad"
)

type Hub struct {
	LastID   int       `json:"lastId"`
	Channels []Channel `json:"channels"`
	Circuits []Circuit `json:"circuits"`
}

type Channel struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Color    string `json:"color"`
	Keyboard bool   `json:"keyboard"`
	Global   bool   `json:"global"`
	IsLocked bool   `json:"isLocked"`
}

type Circuit struct {
	ID             string        `json:"id"`
	Name           string        `json:"name"`
	Elements       Elements      `json:"elements"`
	Camera         CircuitCamera `json:"camera"`
	TicksPerSecond int           `json:"ticksPerSecond"`
	Layer          int           `json:"layer"`
	Behind         int           `json:"behind"`
	Ahead          int           `json:"ahead"`
	Selector       Selector      `json:"selector"`
	Faceplate      Faceplate     `json:"faceplate"`
}

type CircuitCamera struct {
	Px   float64 `json:"px"`
	Py   float64 `json:"py"`
	Zoom float64 `json:"zoom"`
}

type Elements struct {
	LastID   int             `json:"lastId"`
	Entities []CircuitEntity `json:"entities"`
}

type CircuitEntity struct {
	Key  string            `json:"key"`
	ID   int               `json:"id"`
	Data CircuitEntityData `json:"data"`
}

type CircuitEntityData interface {
	isCircuitEntityData()
}

type CircuitEntityDataBase struct {
	Layer    int        `json:"layer"`
	Position Position2D `json:"position"`
	Rotation Direction  `json:"rotation"`
	Name     string     `json:"name"`
}

type Position2D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

//entities and their data bodies, and the isCircuitEntityData method to satisfy the interface

type ToggleData struct {
	CircuitEntityDataBase
	Value bool `json:"value"`
}

type ButtonData struct {
	CircuitEntityDataBase
	Color   string `json:"color"`
	IsDebug bool   `json:"isDebug,omitempty"` // only seen on one sample entity
}

type FuncData struct {
	CircuitEntityDataBase
	Func     FunctionType `json:"func"`
	InCount  int          `json:"inCount"`
	OutCount int          `json:"outCount"`
	Flip     bool         `json:"flip"`
}

type LightData struct {
	CircuitEntityDataBase
	Color string `json:"color"`
}

// the curve data
type curvePoint struct {
	Time       float64 `json:"time"`
	Value      float64 `json:"value"`
	InTangent  float64 `json:"inTangent"`
	OutTangent float64 `json:"outTangent"`
	InWeight   float64 `json:"inWeight"`
	OutWeight  float64 `json:"outWeight"`
	Mode       int     `json:"mode"`
}

type Curve struct {
	Points []curvePoint `json:"keys"`
}

type CurveData struct {
	CircuitEntityDataBase
	Curve Curve `json:"curve"`
}

//--------

type PortData struct {
	CircuitEntityDataBase
	Label string `json:"label"`
}

type BusData struct {
	CircuitEntityDataBase
	Size int  `json:"size"`
	Flip bool `json:"flip"`
}

type CircuitRefData struct { // "circuit" block referencing another Circuit by ID
	CircuitEntityDataBase
	CircuitID string `json:"circuitId"`
}

type ValveData struct {
	CircuitEntityDataBase
	Flip bool `json:"flip"`
}

type TimeData struct {
	CircuitEntityDataBase
	Speed float64 `json:"speed"`
}

type ClockData struct {
	CircuitEntityDataBase
	Off int `json:"off"`
	On  int `json:"on"`
}

type ConstantData struct {
	CircuitEntityDataBase
	Value float64 `json:"value"`
}

// display and via carry no extra fields beyond the base in the sample.
type DisplayData struct {
	CircuitEntityDataBase
}

type ViaData struct {
	CircuitEntityDataBase
}

type WireData struct {
	End Position2D `json:"end"`
	CircuitEntityDataBase
}

func (CircuitEntityDataBase) isCircuitEntityData() {}
func (WireData) isCircuitEntityData()              {}
func (ViaData) isCircuitEntityData()               {}
func (DisplayData) isCircuitEntityData()           {}
func (ConstantData) isCircuitEntityData()          {}
func (ClockData) isCircuitEntityData()             {}
func (TimeData) isCircuitEntityData()              {}
func (ValveData) isCircuitEntityData()             {}
func (CircuitRefData) isCircuitEntityData()        {}
func (BusData) isCircuitEntityData()               {}
func (PortData) isCircuitEntityData()              {}
func (CurveData) isCircuitEntityData()             {}
func (LightData) isCircuitEntityData()             {}
func (FuncData) isCircuitEntityData()              {}
func (ButtonData) isCircuitEntityData()            {}
func (ToggleData) isCircuitEntityData()            {}

// json handling for the Entity struct, which has a polymorphic Data field based on the Key.

func (e *CircuitEntity) UnmarshalJSON(raw []byte) error {
	var shell struct {
		Key  string          `json:"key"`
		ID   int             `json:"id"`
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(raw, &shell); err != nil {
		return fmt.Errorf("unmarshal entity shell: %w", err)
	}
	e.Key = shell.Key
	e.ID = shell.ID

	var data CircuitEntityData
	switch shell.Key {
	case KeyToggle:
		var d ToggleData
		if err := json.Unmarshal(shell.Data, &d); err != nil {
			return fmt.Errorf("unmarshal %s data: %w", shell.Key, err)
		}
		data = d
	case KeyButton:
		var d ButtonData
		if err := json.Unmarshal(shell.Data, &d); err != nil {
			return fmt.Errorf("unmarshal %s data: %w", shell.Key, err)
		}
		data = d
	case KeyFunc:
		var d FuncData
		if err := json.Unmarshal(shell.Data, &d); err != nil {
			return fmt.Errorf("unmarshal %s data: %w", shell.Key, err)
		}
		data = d
	case KeyLight:
		var d LightData
		if err := json.Unmarshal(shell.Data, &d); err != nil {
			return fmt.Errorf("unmarshal %s data: %w", shell.Key, err)
		}
		data = d
	case KeyCurve:
		var d CurveData
		if err := json.Unmarshal(shell.Data, &d); err != nil {
			return fmt.Errorf("unmarshal %s data: %w", shell.Key, err)
		}
		data = d
	case KeyPort:
		var d PortData
		if err := json.Unmarshal(shell.Data, &d); err != nil {
			return fmt.Errorf("unmarshal %s data: %w", shell.Key, err)
		}
		data = d
	case KeyBus:
		var d BusData
		if err := json.Unmarshal(shell.Data, &d); err != nil {
			return fmt.Errorf("unmarshal %s data: %w", shell.Key, err)
		}
		data = d
	case KeyCircuit:
		var d CircuitRefData
		if err := json.Unmarshal(shell.Data, &d); err != nil {
			return fmt.Errorf("unmarshal %s data: %w", shell.Key, err)
		}
		data = d
	case KeyValve:
		var d ValveData
		if err := json.Unmarshal(shell.Data, &d); err != nil {
			return fmt.Errorf("unmarshal %s data: %w", shell.Key, err)
		}
		data = d
	case KeyTime:
		var d TimeData
		if err := json.Unmarshal(shell.Data, &d); err != nil {
			return fmt.Errorf("unmarshal %s data: %w", shell.Key, err)
		}
		data = d
	case KeyClock:
		var d ClockData
		if err := json.Unmarshal(shell.Data, &d); err != nil {
			return fmt.Errorf("unmarshal %s data: %w", shell.Key, err)
		}
		data = d
	case KeyConstant:
		var d ConstantData
		if err := json.Unmarshal(shell.Data, &d); err != nil {
			return fmt.Errorf("unmarshal %s data: %w", shell.Key, err)
		}
		data = d
	case KeyDisplay:
		var d DisplayData
		if err := json.Unmarshal(shell.Data, &d); err != nil {
			return fmt.Errorf("unmarshal %s data: %w", shell.Key, err)
		}
		data = d
	case KeyVia:
		var d ViaData
		if err := json.Unmarshal(shell.Data, &d); err != nil {
			return fmt.Errorf("unmarshal %s data: %w", shell.Key, err)
		}
		data = d
	case KeyWire:
		var d WireData
		if err := json.Unmarshal(shell.Data, &d); err != nil {
			return fmt.Errorf("unmarshal %s data: %w", shell.Key, err)
		}
		data = d
	default:
		return fmt.Errorf("unknown entity key: %q", shell.Key)
	}
	e.Data = data
	return nil
}

func (e CircuitEntity) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Key  string            `json:"key"`
		ID   int               `json:"id"`
		Data CircuitEntityData `json:"data"`
	}{e.Key, e.ID, e.Data})
}
