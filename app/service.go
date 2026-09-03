package app

import (
	"GDv2/utils"
	"GDv2/world"
	"GDv2/world/model"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	saveEditor *world.SaveEditor
}

func NewService(path string) (*Service, error) {
	saveEditor, err := world.NewSaveEditor(path)
	if err != nil {
		return nil, fmt.Errorf("create save editor: %w", err)
	}

	return &Service{
		saveEditor: saveEditor,
	}, nil
}

func (s *Service) GetCircuits() ([]*model.Circuit, error) {
	circuits, err := s.saveEditor.GetCircuits()
	if err != nil {
		return nil, fmt.Errorf("get circuits: %w", err)
	}
	return circuits, nil
}

func (s *Service) GetCircuit(input string) (*model.Circuit, error) {
	// 1. Try UUID
	if _, err := uuid.Parse(input); err == nil {
		return s.saveEditor.GetCircuitByID(input)
	}
	// 2. Try index (pure integer)
	if idx, err := strconv.Atoi(input); err == nil {
		return s.saveEditor.GetCircuitByIndex(idx)
	}
	// 3. Fall back to name
	return s.saveEditor.GetCircuitByName(input)
}

func (s *Service) ExportCircuit(circuitIndex int, path string) error {
	// this doesn't solve the diamond relation case but it works, not efficent but if it works do NOT touch it
	circuit, err := s.saveEditor.GetCircuitByIndex(circuitIndex)
	if err != nil {
		return fmt.Errorf("get circuit: %w", err)
	}

	subCircuitIDs, err := s.saveEditor.GetSubCircuitIDs(circuit)
	if err != nil {
		return fmt.Errorf("get citcuit subcircuit IDs: %w", err)
	}

	if len(subCircuitIDs) > 0 {
		for _, id := range subCircuitIDs {
			subCircuit, err := s.saveEditor.GetCircuitByID(id)
			if err != nil {
				return fmt.Errorf("get subcircuit %s: %w", id, err)
			}

			subPath := filepath.Join(filepath.Dir(path), subCircuit.Name+" .json")

			index, err := s.saveEditor.GetCircuitIndex(subCircuit)
			if err != nil {
				return fmt.Errorf("get subcircuit index: %w", err)
			}

			if err := s.ExportCircuit(index, subPath); err != nil {
				return fmt.Errorf("export subcircuit %s: %w", subCircuit.Name, err)
			}
		}
	}

	circuitData, err := json.MarshalIndent(circuit, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal circuit: %w", err)
	}

	return utils.WriteFile(path, circuitData)
}

func (s *Service) ImportCircuit(paths []string) error {
	// Pass 1: load every circuit, decide its final ID, build the old->new ID map
	circuits := make([]*model.Circuit, 0, len(paths))
	idMap := make(map[string]string) // old ID -> new ID

	for _, path := range paths {
		data, err := utils.LoadFile(path)
		if err != nil {
			return fmt.Errorf("load %s: %w", path, err)
		}

		var circuit *model.Circuit
		if err := json.Unmarshal(data, &circuit); err != nil {
			return fmt.Errorf("unmarshal %s: %w", path, err)
		}

		oldID := circuit.ID

		if existing, err := s.saveEditor.GetCircuitByID(circuit.ID); err == nil && existing != nil {
			circuit.ID = utils.GenerateUUID()
		}
		idMap[oldID] = circuit.ID // ALWAYS record the mapping, even if unchanged (oldID -> oldID)

		circuit.Name = fmt.Sprintf("%s (Imported)", circuit.Name)
		circuits = append(circuits, circuit)
	}

	// Pass 2: rewrite every sub-circuit reference using the completed map,
	// THEN actually add each circuit to the save
	for _, circuit := range circuits {
		s.remapSubCircuitReferences(circuit, idMap)
		if err := s.saveEditor.ImportCircuit(circuit); err != nil {
			return fmt.Errorf("import circuit %s: %w", circuit.ID, err)
		}
	}

	return nil
}

func (s *Service) remapSubCircuitReferences(circuit *model.Circuit, idMap map[string]string) {
	for i := range circuit.Elements.Entities {
		ref, ok := circuit.Elements.Entities[i].Data.(model.CircuitRefData)
		if !ok {
			continue
		}
		if newID, exists := idMap[ref.CircuitID]; exists {
			ref.CircuitID = newID
			circuit.Elements.Entities[i].Data = ref // write the modified value back
		}
	}
}

func (s *Service) SaveFile() error {
	rawData, err := json.MarshalIndent(s.saveEditor.SaveFile, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal save file: %w", err)
	}

	savePath := s.saveEditor.SaveFile.Path
	if err := s.writeBackup(savePath); err != nil {
		return fmt.Errorf("write backup: %w", err)
	}

	return utils.WriteFile(savePath, rawData)
}

func (s *Service) writeBackup(savePath string) error {
	backupBaseDir, projectName, err := resolveBackupBaseDirectory(savePath)
	if err != nil {
		return err
	}
	projectBackupDir := filepath.Join(backupBaseDir, projectName)

	originalData, err := utils.LoadFile(savePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Nothing to back up on first save.
			return nil
		}
		return fmt.Errorf("read original save file %s: %w", savePath, err)
	}

	if err := os.MkdirAll(projectBackupDir, 0755); err != nil {
		return fmt.Errorf("create backup directory for project %s: %w", projectName, err)
	}

	timestamp := time.Now().Format("20060102_150405")
	backupFileName := fmt.Sprintf("%s_%s.json", projectName, timestamp)
	backupPath := filepath.Join(projectBackupDir, backupFileName)

	if err := utils.WriteFile(backupPath, originalData); err != nil {
		return fmt.Errorf("write backup file %s: %w", backupPath, err)
	}

	return nil
}

func resolveBackupBaseDirectory(savePath string) (string, string, error) {
	cleanPath := filepath.Clean(savePath)
	current := filepath.Dir(cleanPath)

	for {
		parent := filepath.Dir(current)
		if parent == current {
			break
		}

		if samePathComponent(filepath.Base(parent), "projects") {
			projectName := filepath.Base(current)
			if projectName == "." || projectName == string(filepath.Separator) || projectName == "" {
				return "", "", fmt.Errorf("project name is empty in save path: %s", savePath)
			}

			backupBaseDir := filepath.Join(parent, "backup")
			return backupBaseDir, projectName, nil
		}

		current = parent
	}

	return "", "", fmt.Errorf("save path does not include projects/<project-name>: %s", savePath)
}

func samePathComponent(actual string, expected string) bool {
	if runtime.GOOS == "windows" {
		return strings.EqualFold(actual, expected)
	}
	return actual == expected
}

func (s *Service) GetBackupFiles() ([]string, error) {
	backupBaseDir, projectName, err := resolveBackupBaseDirectory(s.saveEditor.SaveFile.Path)
	if err != nil {
		return nil, err
	}
	projectBackupDir := filepath.Join(backupBaseDir, projectName)

	entries, err := os.ReadDir(projectBackupDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("read backup directory %s: %w", projectBackupDir, err)
	}

	backupFiles := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if strings.ToLower(filepath.Ext(name)) != ".json" {
			continue
		}

		backupFiles = append(backupFiles, filepath.Join(projectBackupDir, name))
	}

	// Timestamped file names sort lexicographically; descending gives newest first.
	sort.Slice(backupFiles, func(i, j int) bool {
		return filepath.Base(backupFiles[i]) > filepath.Base(backupFiles[j])
	})

	return backupFiles, nil
}

func (s *Service) RestoreBackupFile(path string) error {
	if path == "" {
		return nil
	}

	backupData, err := utils.LoadFile(path)
	if err != nil {
		return fmt.Errorf("load backup file %s: %w", path, err)
	}

	utils.WriteFile(s.saveEditor.SaveFile.Path, backupData)

	err = s.saveEditor.ReloadSaveFile(s.saveEditor.SaveFile.Path)
	if err != nil {
		return fmt.Errorf("reload save file: %w", err)
	}

	return nil
}

type ValidationIssue struct {
	Severity string
	Message  string
}

func (s *Service) ValidateSavefile() ([]ValidationIssue, error) {
	var issues []ValidationIssue

	// circuit to circuit checks
	circuits := s.saveEditor.SaveFile.Hub.Circuits
	circuitNames := make(map[string]int, len(circuits))
	circuitIDs := make(map[string]int, len(circuits))

	for i := range circuits {
		c := &circuits[i]

		// Game LastID check
		if c.Elements.LastID < len(c.Elements.Entities) {
			issues = append(issues, ValidationIssue{
				Severity: "error",
				Message:  fmt.Sprintf("circuit %q has LastID %d but contains %d entities", c.Name, c.Elements.LastID, len(c.Elements.Entities)),
			})
		}

		// Circuit ID format check
		if c.ID == "" {
			issues = append(issues, ValidationIssue{
				Severity: "error",
				Message:  fmt.Sprintf("circuit %q has empty ID", c.Name),
			})
		}
		if _, err := uuid.Parse(c.ID); err != nil {
			issues = append(issues, ValidationIssue{
				Severity: "error",
				Message:  fmt.Sprintf("circuit %q has invalid ID %q: %v", c.Name, c.ID, err),
			})
		}

		// Run per-circuit validation exactly once
		issues = append(issues, s.validateCircuit(c)...)

		// Track name and ID duplicates
		circuitNames[c.Name]++
		circuitIDs[c.ID]++
	}

	for name, count := range circuitNames {
		if count > 1 {
			issues = append(issues, ValidationIssue{
				Severity: "error",
				Message:  fmt.Sprintf("duplicate circuit name %q appears %d times", name, count),
			})
		}
	}
	for id, count := range circuitIDs {
		if count > 1 {
			issues = append(issues, ValidationIssue{
				Severity: "error",
				Message:  fmt.Sprintf("duplicate circuit ID %s appears %d times", id, count),
			})
		}
	}

	return issues, nil
}

func (s *Service) validateCircuit(circuit *model.Circuit) []ValidationIssue {
	var issues []ValidationIssue
	seenIDs := make(map[int]int, len(circuit.Elements.Entities))

	for i := range circuit.Elements.Entities {
		entity := &circuit.Elements.Entities[i]
		seenIDs[entity.ID]++

		// sub circuit reference check
		if refData, ok := entity.Data.(model.CircuitRefData); ok {
			if _, err := s.saveEditor.GetCircuitByID(refData.CircuitID); err != nil {
				issues = append(issues, ValidationIssue{
					Severity: "error",
					Message:  fmt.Sprintf("circuit %q references non-existent sub-circuit ID %v", circuit.Name, refData.CircuitID),
				})
			}
		}
	}

	// Duplicate part ID check
	for id, count := range seenIDs {
		if count > 1 {
			issues = append(issues, ValidationIssue{
				Severity: "error",
				Message:  fmt.Sprintf("circuit %q has duplicate part ID %d (found %d times)", circuit.Name, id, count),
			})
		}
	}

	return issues
}

func (s *Service) DeleteCircuit(circuitIndex int) error {
	circuit, err := s.saveEditor.GetCircuitByIndex(circuitIndex)
	if err != nil {
		return fmt.Errorf("get circuit: %w", err)
	}

	if err := s.saveEditor.DeleteCircuit(circuit); err != nil {
		return fmt.Errorf("delete circuit: %w", err)
	}
	return nil
}
