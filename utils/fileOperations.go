package utils

import "os"

func WriteFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

func LoadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}
