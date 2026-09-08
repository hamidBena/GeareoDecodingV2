package autoBuilder

import (
	"GDv2/world"
	"GDv2/world/model"
	"fmt"
)

type DisplayCircuit struct {
	Position model.Position2D
	Width    int
	Height   int
}

type DisplayContext struct {
	AnalogMemory *model.Circuit
	LED          *model.Circuit
}

func (d *DisplayCircuit) Build(targetCircuit *model.Circuit, se *world.SaveEditor) error {
	if targetCircuit == nil {
		return fmt.Errorf("target circuit is nil")
	}

	displayContext, err := d.checkPreconditions(se)
	if err != nil {
		return err
	}

	err = d.build(targetCircuit, displayContext, se)
	if err != nil {
		return err
	}

	return nil
}

func (d *DisplayCircuit) checkPreconditions(se *world.SaveEditor) (DisplayContext, error) {
	analogMem, err := se.GetCircuitByName("Analog Memory")
	if err != nil {
		return DisplayContext{}, fmt.Errorf("Error occurred while fetching Analog Memory circuit: %w", err)
	}

	led, err := se.GetCircuitByName("LED")
	if err != nil {
		return DisplayContext{}, fmt.Errorf("Error occurred while fetching LED circuit: %w", err)
	}

	displayContext := DisplayContext{
		AnalogMemory: analogMem,
		LED:          led,
	}

	if d.Width <= 0 || d.Height <= 0 {
		return DisplayContext{}, fmt.Errorf("Invalid display dimensions: width=%d, height=%d", d.Width, d.Height)
	}

	return displayContext, nil
}

func (d *DisplayCircuit) build(circuit *model.Circuit, ctx DisplayContext, se *world.SaveEditor) error {
	registersYOffset := 30
	viaYOffset := registersYOffset + 50 + d.Height*2
	cellsLayer := -2
	pixelsLayer := 0

	registerOutputPositions := []model.Position2D{}
	pixelPositions := []model.Position2D{}
	for column := 0; column < d.Width; column++ {
		for row := 0; row < d.Height; row++ {
			position := d.Position.Add(model.Position2D{
				X: column * 9,
				Y: row*2 + registersYOffset,
			})

			registerOutputPositions = append(registerOutputPositions, position.Add(model.Position2D{X: 1, Y: 0}))
			storeCellEntities := []model.CircuitEntity{
				model.NewCircuitRefWithLayer(position, model.Right, ctx.AnalogMemory.ID, cellsLayer),
				model.NewFuncWithLayer(
					position.Add(model.Position2D{X: -2, Y: -1}),
					model.Right,
					model.FuncAnd,
					2, cellsLayer,
				),
				model.NewFuncWithLayer(
					position.Add(model.Position2D{X: -4, Y: -1}),
					model.Right,
					model.FuncEqual,
					2, cellsLayer,
				),
				model.NewConstantWithLayer(
					position.Add(model.Position2D{X: -6, Y: -2}),
					model.Right,
					float64(row+1), cellsLayer,
				),
				model.NewWireWithLayer(
					position.Add(model.Position2D{X: -5, Y: -1}),
					position.Add(model.Position2D{X: -6, Y: -1}),
					cellsLayer,
				),
				model.NewWireWithLayer(
					position.Add(model.Position2D{X: -6, Y: -1}),
					position.Add(model.Position2D{X: -7, Y: -2}),
					cellsLayer,
				),
				model.NewWireWithLayer(
					position.Add(model.Position2D{X: -7, Y: -2}),
					position.Add(model.Position2D{X: -7, Y: -4}),
					cellsLayer,
				),
				model.NewWireWithLayer(
					position.Add(model.Position2D{X: -3, Y: -2}),
					position.Add(model.Position2D{X: -4, Y: -3}),
					cellsLayer,
				),
				model.NewWireWithLayer(
					position.Add(model.Position2D{X: -4, Y: -3}),
					position.Add(model.Position2D{X: -4, Y: -5}),
					cellsLayer,
				),
				model.NewWireWithLayer(
					position.Add(model.Position2D{X: -1, Y: 0}),
					position.Add(model.Position2D{X: -1, Y: 2}),
					cellsLayer,
				),
				model.NewWireWithLayer(
					position.Add(model.Position2D{X: 1, Y: -1}),
					position.Add(model.Position2D{X: 1, Y: 1}),
					cellsLayer,
				),
			}

			for _, entity := range storeCellEntities {
				if err := se.AddPartToCircuit(circuit, entity); err != nil {
					return err
				}
			}
		}

		// per collumn data
		perColumnPositions := d.Position.Add(model.Position2D{
			X: column * 9,
			Y: registersYOffset,
		})

		connectors := []model.CircuitEntity{
			model.NewWireWithLayer(
				perColumnPositions.Add(model.Position2D{X: -7, Y: -4}),
				perColumnPositions.Add(model.Position2D{X: 2, Y: -4}),
				cellsLayer,
			),
			model.NewWireWithLayer(
				perColumnPositions.Add(model.Position2D{X: -1, Y: d.Height * 2}),
				perColumnPositions.Add(model.Position2D{X: 8, Y: d.Height * 2}),
				cellsLayer,
			),
			model.NewWireWithLayer(
				perColumnPositions.Add(model.Position2D{X: 1, Y: d.Height*2 - 1}),
				perColumnPositions.Add(model.Position2D{X: 10, Y: d.Height*2 - 1}),
				cellsLayer,
			),

			model.NewFuncWithLayer(
				perColumnPositions.Add(model.Position2D{X: -4, Y: -6}),
				model.Up,
				model.FuncEqual,
				2, cellsLayer,
			),
			model.NewConstantWithLayer(
				perColumnPositions.Add(model.Position2D{X: -3, Y: -8}),
				model.Up,
				float64(column+1),
				cellsLayer,
			),
			model.NewWireWithLayer(
				perColumnPositions.Add(model.Position2D{X: -4, Y: -7}),
				perColumnPositions.Add(model.Position2D{X: 5, Y: -7}),
				cellsLayer,
			),
		}

		for _, entity := range connectors {
			if err := se.AddPartToCircuit(circuit, entity); err != nil {
				return err
			}
		}
	}

	for row := 0; row < d.Width; row++ {
		for column := 0; column < d.Height; column++ {
			pixelPosition := d.Position.Add(model.Position2D{
				X: column,
				Y: row,
			})
			pixelPositions = append(pixelPositions, pixelPosition.Add(model.Position2D{X: 0, Y: viaYOffset}))

			pixelEntities := []model.CircuitEntity{
				model.NewCircuitRefWithLayer(pixelPosition, model.Right, ctx.LED.ID, pixelsLayer),

				model.NewVia(pixelPosition.Add(model.Position2D{X: 0, Y: viaYOffset}), model.Right),
			}

			for _, entity := range pixelEntities {
				if err := se.AddPartToCircuit(circuit, entity); err != nil {
					return err
				}
			}
		}
	}

	for i, pos := range registerOutputPositions {
		registerOutputWire := model.NewWireWithLayer(
			pos,
			pixelPositions[i],
			cellsLayer,
		)

		pixelWire := model.NewWireWithLayer(
			pixelPositions[i],
			pixelPositions[i].Add(model.Position2D{X: 0, Y: -viaYOffset}),
			pixelsLayer,
		)

		if err := se.AddPartToCircuit(circuit, registerOutputWire); err != nil {
			return err
		}
		if err := se.AddPartToCircuit(circuit, pixelWire); err != nil {
			return err
		}
	}

	return nil
}
