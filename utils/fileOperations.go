package utils

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
)

func WriteFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

func LoadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func LoadIcon(path string) (fyne.Resource, error) {
	data, err := LoadFile(path)
	if err != nil {
		return nil, err
	}

	return fyne.NewStaticResource(path, data), nil
}

func LoadCSV(dataFile string) ([]*float64, error) {
	file, err := os.Open(dataFile)
	if err != nil {
		return nil, fmt.Errorf("open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var result []*float64

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read CSV row: %w", err)
		}

		for _, text := range row {
			text = strings.TrimSpace(text)
			if text == "" {
				continue
			}

			value, err := strconv.ParseFloat(text, 64)
			if err != nil {
				return nil, fmt.Errorf("parse value %q: %w", text, err)
			}

			result = append(result, &value)
		}
	}

	return result, nil
}
