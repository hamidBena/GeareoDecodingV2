package logging

import (
	"log"
	"os"
)

var std *log.Logger

func Init(path string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	std = log.New(f, "", log.LstdFlags)
	return nil
}

func Info(msg string, args ...any)  { std.Printf("[INFO] "+msg, args...) }
func Warn(msg string, args ...any)  { std.Printf("[WARN] "+msg, args...) }
func Error(msg string, args ...any) { std.Printf("[ERROR] "+msg, args...) }
