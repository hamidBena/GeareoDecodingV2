package world

import (
	"GDv2/world/model"
	"fmt"
)

// circuit lookup functions
func (se *SaveEditor) GetCircuits() ([]*model.Circuit, error) {
	circuits := make([]*model.Circuit, len(se.SaveFile.Hub.Circuits))
	for i := range se.SaveFile.Hub.Circuits {
		circuits[i] = &se.SaveFile.Hub.Circuits[i]
	}
	return circuits, nil
}

func (se *SaveEditor) GetCircuitByName(circuitName string) (*model.Circuit, error) {
	for i := range se.SaveFile.Hub.Circuits {
		if se.SaveFile.Hub.Circuits[i].Name == circuitName {
			return &se.SaveFile.Hub.Circuits[i], nil
		}
	}
	return nil, fmt.Errorf("circuit not found: %s", circuitName)
}

func (se *SaveEditor) GetCircuitByID(circuitID string) (*model.Circuit, error) {
	for i := range se.SaveFile.Hub.Circuits {
		if se.SaveFile.Hub.Circuits[i].ID == circuitID {
			return &se.SaveFile.Hub.Circuits[i], nil
		}
	}
	return nil, fmt.Errorf("circuit not found: %s", circuitID)
}

func (se *SaveEditor) GetCircuitByIndex(index int) (*model.Circuit, error) {
	if index < 0 || index >= len(se.SaveFile.Hub.Circuits) {
		return nil, fmt.Errorf("index out of bounds: %d", index)
	}
	return &se.SaveFile.Hub.Circuits[index], nil
}

func (se *SaveEditor) GetCircuitIndex(circuit *model.Circuit) (int, error) {
	for i, cir := range se.SaveFile.Hub.Circuits {
		if circuit.ID == cir.ID {
			return i, nil
		}
	}

	return 0, fmt.Errorf("Circuit index not found")
}

// circuit management functions
func (se *SaveEditor) DeleteCircuit(circuit *model.Circuit) error {
	index, err := se.GetCircuitIndex(circuit)
	if err != nil {
		return err
	}

	se.SaveFile.Hub.Circuits = append(se.SaveFile.Hub.Circuits[:index], se.SaveFile.Hub.Circuits[index+1:]...)
	return nil
}

func (se *SaveEditor) AddPartToCircuit(circuit *model.Circuit, part *model.CircuitEntity) error {
	circuit.Elements.LastID++
	part.ID = circuit.Elements.LastID

	return nil
	
}
