package utils

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/ncruces/zenity"
)

func OpenSaveFile() (string, error) {
	dir := defaultProjectsDir()
	return zenity.SelectFile(
		zenity.Title("Select Game Save File"),
		zenity.Filename(dir),
		zenity.FileFilters{
			{Name: "JSON Files", Patterns: []string{"*.json"}},
		},
	)
}

func SaveOutputFile() (string, error) {
	dir := defaultProjectsDir()
	filename := filepath.Join(dir, "modified_save.json")
	return zenity.SelectFileSave(
		zenity.Title("Save Modified Circuit File"),
		zenity.Filename(filename),
		zenity.FileFilters{
			{Name: "JSON Files", Patterns: []string{"*.json"}},
		},
	)
}

func defaultProjectsDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	var dir string
	switch runtime.GOOS {
	case "windows":
		dir = filepath.Join(home, "AppData", "LocalLow", "WitWeld", "Geareo", "projects")
	case "darwin":
		dir = filepath.Join(home, "Library", "Application Support", "WitWeld", "Geareo", "projects")
	case "linux":
		dir = filepath.Join(home, ".local", "share", "WitWeld", "Geareo", "projects")
	default:
		dir = filepath.Join(home, ".local", "share", "WitWeld", "Geareo", "projects")
	}
	_ = os.MkdirAll(dir, 0755)
	return dir
}

func OpenSaveFiles() ([]string, error) {
	dir := defaultProjectsDir()
	// Use the zenity API's multiple-file option (either SelectFileMultiple
	// or SelectFile with zenity.Multiple). Adjust call to match the package.
	return zenity.SelectFileMultiple(
		zenity.Title("Select CIRCUIT JSON files to import"),
		zenity.Filename(dir),
		zenity.FileFilters{
			{Name: "JSON Files", Patterns: []string{"*.json"}},
		},
	)
}
