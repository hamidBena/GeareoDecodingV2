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
	X int `json:"x"`
	Y int `json:"y"`
}

func (p Position2D) Add(other Position2D) Position2D {
	return Position2D{
		X: p.X + other.X,
		Y: p.Y + other.Y,
	}
}
func (p Position2D) Sub(other Position2D) Position2D {
	return Position2D{
		X: p.X - other.X,
		Y: p.Y - other.Y,
	}
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

func (e *CircuitEntity) Offset(X, Y int) {
	offX := X
	offY := Y

	switch data := e.Data.(type) {
	case ToggleData:
		data.Position.X += offX
		data.Position.Y += offY
		e.Data = data

	case ButtonData:
		data.Position.X += offX
		data.Position.Y += offY
		e.Data = data

	case FuncData:
		data.Position.X += offX
		data.Position.Y += offY
		e.Data = data

	case LightData:
		data.Position.X += offX
		data.Position.Y += offY
		e.Data = data

	case CurveData:
		data.Position.X += offX
		data.Position.Y += offY
		e.Data = data

	case PortData:
		data.Position.X += offX
		data.Position.Y += offY
		e.Data = data

	case BusData:
		data.Position.X += offX
		data.Position.Y += offY
		e.Data = data

	case CircuitRefData:
		data.Position.X += offX
		data.Position.Y += offY
		e.Data = data

	case ValveData:
		data.Position.X += offX
		data.Position.Y += offY
		e.Data = data

	case TimeData:
		data.Position.X += offX
		data.Position.Y += offY
		e.Data = data

	case ClockData:
		data.Position.X += offX
		data.Position.Y += offY
		e.Data = data

	case ConstantData:
		data.Position.X += offX
		data.Position.Y += offY
		e.Data = data

	case DisplayData:
		data.Position.X += offX
		data.Position.Y += offY
		e.Data = data

	case ViaData:
		data.Position.X += offX
		data.Position.Y += offY
		e.Data = data

	case WireData:
		data.Position.X += offX
		data.Position.Y += offY
		e.Data = data
	}
}

func newEntity(key string, data CircuitEntityData) CircuitEntity {
	return CircuitEntity{
		Key:  key,
		Data: data,
	}
}

func newBase(position Position2D, rotation Direction) CircuitEntityDataBase {
	return newBaseWithLayer(position, rotation, 0)
}

func newBaseWithLayer(position Position2D, rotation Direction, layer int) CircuitEntityDataBase {
	return CircuitEntityDataBase{
		Layer:    layer,
		Position: position,
		Rotation: rotation,
	}
}

func NewToggle(position Position2D, rotation Direction, value bool) CircuitEntity {
	return NewToggleWithLayer(position, rotation, value, 0)
}

func NewToggleWithLayer(position Position2D, rotation Direction, value bool, layer int) CircuitEntity {
	return newEntity(KeyToggle, ToggleData{
		CircuitEntityDataBase: newBaseWithLayer(position, rotation, layer),
		Value:                 value,
	})
}

func NewButton(position Position2D, rotation Direction, isDebug bool) CircuitEntity {
	return NewButtonWithLayer(position, rotation, isDebug, 0)
}

func NewButtonWithLayer(position Position2D, rotation Direction, isDebug bool, layer int) CircuitEntity {
	return newEntity(KeyButton, ButtonData{
		CircuitEntityDataBase: newBaseWithLayer(position, rotation, layer),
		Color:                 "#DD0000",
		IsDebug:               isDebug,
	})
}

func NewFunc(position Position2D, rotation Direction, function FunctionType, inCount int) CircuitEntity {
	return NewFuncWithLayer(position, rotation, function, inCount, 0)
}

func NewFuncWithLayer(position Position2D, rotation Direction, function FunctionType, inCount, layer int) CircuitEntity {
	return newEntity(KeyFunc, FuncData{
		CircuitEntityDataBase: newBaseWithLayer(position, rotation, layer),
		Func:                  function,
		InCount:               inCount,
		OutCount:              1,
	})
}

func NewLight(position Position2D, rotation Direction, color string) CircuitEntity {
	return NewLightWithLayer(position, rotation, color, 0)
}

func NewLightWithLayer(position Position2D, rotation Direction, color string, layer int) CircuitEntity {
	return newEntity(KeyLight, LightData{
		CircuitEntityDataBase: newBaseWithLayer(position, rotation, layer),
		Color:                 color,
	})
}

func NewCurve(position Position2D, rotation Direction, curve Curve) CircuitEntity {
	return NewCurveWithLayer(position, rotation, curve, 0)
}

func NewCurveWithLayer(position Position2D, rotation Direction, curve Curve, layer int) CircuitEntity {
	return newEntity(KeyCurve, CurveData{
		CircuitEntityDataBase: newBaseWithLayer(position, rotation, layer),
		Curve:                 curve,
	})
}

func NewPort(position Position2D, rotation Direction, label string) CircuitEntity {
	return NewPortWithLayer(position, rotation, label, 0)
}

func NewPortWithLayer(position Position2D, rotation Direction, label string, layer int) CircuitEntity {
	return newEntity(KeyPort, PortData{
		CircuitEntityDataBase: newBaseWithLayer(position, rotation, layer),
		Label:                 label,
	})
}

func NewBus(position Position2D, rotation Direction, size int, flip bool) CircuitEntity {
	return NewBusWithLayer(position, rotation, size, flip, 0)
}

func NewBusWithLayer(position Position2D, rotation Direction, size int, flip bool, layer int) CircuitEntity {
	return newEntity(KeyBus, BusData{
		CircuitEntityDataBase: newBaseWithLayer(position, rotation, layer),
		Size:                  size,
		Flip:                  flip,
	})
}

func NewCircuitRef(position Position2D, rotation Direction, circuitID string) CircuitEntity {
	return NewCircuitRefWithLayer(position, rotation, circuitID, 0)
}

func NewCircuitRefWithLayer(position Position2D, rotation Direction, circuitID string, layer int) CircuitEntity {
	return newEntity(KeyCircuit, CircuitRefData{
		CircuitEntityDataBase: newBaseWithLayer(position, rotation, layer),
		CircuitID:             circuitID,
	})
}

func NewValve(position Position2D, rotation Direction, flip bool) CircuitEntity {
	return NewValveWithLayer(position, rotation, flip, 0)
}

func NewValveWithLayer(position Position2D, rotation Direction, flip bool, layer int) CircuitEntity {
	return newEntity(KeyValve, ValveData{
		CircuitEntityDataBase: newBaseWithLayer(position, rotation, layer),
		Flip:                  flip,
	})
}

func NewTime(position Position2D, rotation Direction, speed float64) CircuitEntity {
	return NewTimeWithLayer(position, rotation, speed, 0)
}

func NewTimeWithLayer(position Position2D, rotation Direction, speed float64, layer int) CircuitEntity {
	return newEntity(KeyTime, TimeData{
		CircuitEntityDataBase: newBaseWithLayer(position, rotation, layer),
		Speed:                 speed,
	})
}

func NewClock(position Position2D, rotation Direction, off, on int) CircuitEntity {
	return NewClockWithLayer(position, rotation, off, on, 0)
}

func NewClockWithLayer(position Position2D, rotation Direction, off, on, layer int) CircuitEntity {
	return newEntity(KeyClock, ClockData{
		CircuitEntityDataBase: newBaseWithLayer(position, rotation, layer),
		Off:                   off,
		On:                    on,
	})
}

func NewConstant(position Position2D, rotation Direction, value float64) CircuitEntity {
	return NewConstantWithLayer(position, rotation, value, 0)
}

func NewConstantWithLayer(position Position2D, rotation Direction, value float64, layer int) CircuitEntity {
	return newEntity(KeyConstant, ConstantData{
		CircuitEntityDataBase: newBaseWithLayer(position, rotation, layer),
		Value:                 value,
	})
}

func NewDisplay(position Position2D, rotation Direction) CircuitEntity {
	return NewDisplayWithLayer(position, rotation, 0)
}

func NewDisplayWithLayer(position Position2D, rotation Direction, layer int) CircuitEntity {
	return newEntity(KeyDisplay, DisplayData{
		CircuitEntityDataBase: newBaseWithLayer(position, rotation, layer),
	})
}

func NewVia(position Position2D, rotation Direction) CircuitEntity {
	return NewViaWithLayer(position, rotation, 0)
}

func NewViaWithLayer(position Position2D, rotation Direction, layer int) CircuitEntity {
	return newEntity(KeyVia, ViaData{
		CircuitEntityDataBase: newBaseWithLayer(position, rotation, layer),
	})
}

func NewWire(start, end Position2D) CircuitEntity {
	return NewWireWithLayer(start, end, 0)
}

func NewWireWithLayer(start, end Position2D, layer int) CircuitEntity {
	diff := Position2D{
		X: end.X - start.X,
		Y: end.Y - start.Y,
	}
	return newEntity(KeyWire, WireData{
		CircuitEntityDataBase: newBaseWithLayer(start, Right, layer),
		End:                   diff,
	})
}
