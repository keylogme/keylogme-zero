package utils

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
)

func ParseFromFile(fname string, d any) error {
	content, err := os.ReadFile(fname)
	if err != nil {
		slog.Info(fmt.Sprintf("Could not open file %s\n", fname))
		return err
	}
	err = json.Unmarshal(content, d)
	if err != nil {
		slog.Info(fmt.Sprintf("Could not parse file %s, file corrupted\n", fname))
		return err
	}
	return nil
}
