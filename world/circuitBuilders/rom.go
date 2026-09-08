package autoBuilder

import (
	"GDv2/world"
	"GDv2/world/model"
	"fmt"
	"math"
)

type ROMCircuit struct {
	Position  model.Position2D
	CellCount int
	Layer     int
	Data      []*float64
}

func (r *ROMCircuit) Build(targetCircuit *model.Circuit, se *world.SaveEditor) error {
	if targetCircuit == nil {
		return fmt.Errorf("target circuit is nil")
	}

	if se == nil {
		return fmt.Errorf("save editor is nil")
	}

	if r.CellCount <= 0 {
		return fmt.Errorf("Invalid cell count: %d", r.CellCount)
	}

	err := r.build(targetCircuit, se)
	if err != nil {
		return err
	}

	return nil
}

func (r *ROMCircuit) build(targetCircuit *model.Circuit, se *world.SaveEditor) error {

	width := int(math.Ceil(math.Sqrt(float64(r.CellCount))))
	height := int(math.Ceil(float64(r.CellCount) / float64(width)))

	for heightIndex := 0; heightIndex < height; heightIndex++ {
		for widthIndex := 0; widthIndex < width; widthIndex++ {
			if widthIndex+heightIndex*width >= r.CellCount {
				break
			}
			anchor := model.Position2D{
				X: r.Position.X + widthIndex*5,
				Y: r.Position.Y + heightIndex*3,
			}
			index := widthIndex + heightIndex*width

			parts := []model.CircuitEntity{
				model.NewFuncWithLayer(anchor.Add(model.Position2D{X: -1, Y: 1}), model.Right, model.FuncEqual, 2, r.Layer),
				model.NewConstantWithLayer(anchor.Add(model.Position2D{X: -3, Y: 1}), model.Right, float64(index+1), r.Layer),
				model.NewPortWithLayer(anchor.Add(model.Position2D{X: -3, Y: 0}), model.Right, "addr", r.Layer),
				model.NewValveWithLayer(anchor.Add(model.Position2D{X: 1, Y: 1}), model.Up, false, r.Layer),
				model.NewPortWithLayer(anchor.Add(model.Position2D{X: 0, Y: 2}), model.Right, "out", r.Layer),
			}

			if r.Data != nil && index < len(r.Data) && r.Data[index] != nil {
				parts = append(parts, model.NewConstantWithLayer(anchor, model.Right, *r.Data[index], r.Layer))
			}

			for _, part := range parts {
				if err := se.AddPartToCircuit(targetCircuit, part); err != nil {
					return fmt.Errorf("Error adding part to target circuit: %w", err)
				}
			}
		}
	}

	return nil
}
