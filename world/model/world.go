package model

//metadata
const (
	MatWood    = "wod"
	MatMetal   = "metal"
	MatPlastic = "plastic"
	MatGlass   = "glass"
)

type Parts struct {
	LastID   int          `json:"lastId"`
	Entities []PartEntity `json:"entities"`
}

type PartEntity struct {
	Key  string         `json:"key"`
	ID   int            `json:"id"`
	Data PartEntityData `json:"data"`
}

type PartEntityData interface {
	isPartEntityData()
}

type PartEntityDataBase struct {
	Color     string    `json:"color"`
	Material  string    `json:"material"`
	Name      string    `json:"name"`
	Transform Transform `json:"_c"`
	Unknown_m M         `json:"_m"`
}

type Transform struct {
	X    float64 `json:"0"`
	Y    float64 `json:"1"`
	Z    float64 `json:"2"`
	XRot float64 `json:"3"`
	YRot float64 `json:"4"`
	ZRot float64 `json:"5"`
}

type M struct{}
