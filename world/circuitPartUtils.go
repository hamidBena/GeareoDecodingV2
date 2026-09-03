package world

import (
	"GDv2/world/model"
	"fmt"
)

// part lookup functions
func (se *SaveEditor) GetCircuitPartByID(circuit *model.Circuit, partID int) ([]*model.CircuitEntity, error) {
	var parts []*model.CircuitEntity
	for i := range circuit.Elements.Entities {
		if circuit.Elements.Entities[i].ID == partID {
			parts = append(parts, &circuit.Elements.Entities[i])
		}
	}
	if len(parts) == 0 {
		return nil, fmt.Errorf("part not found: %d", partID)
	}
	return parts, nil
}

func (se *SaveEditor) GetCircuitPartsBetweenPositions(circuit *model.Circuit, position1, position2 model.Position2D) ([]*model.CircuitEntity, error) {
	minX, maxX := min(position1.X, position2.X), max(position1.X, position2.X)
	minY, maxY := min(position1.Y, position2.Y), max(position1.Y, position2.Y)

	var parts []*model.CircuitEntity
	for i := range circuit.Elements.Entities {
		part := &circuit.Elements.Entities[i] // pointer into the real slice, not a loop-var copy

		base, ok := part.Data.(model.CircuitEntityDataBase)
		if !ok {
			continue // skip entity types that don't carry a Position (if any)
		}

		pos := base.Position
		if pos.X >= minX && pos.X <= maxX && pos.Y >= minY && pos.Y <= maxY {
			parts = append(parts, part)
		}
	}

	return parts, nil
}

func (se *SaveEditor) GetSubCircuitIDs(circuit *model.Circuit) ([]string, error) {
	var ids []string
	for _, part := range circuit.Elements.Entities {
		if subCircuit, ok := part.Data.(model.CircuitRefData); ok {
			ids = append(ids, subCircuit.CircuitID)
		}
	}
	return ids, nil
}
