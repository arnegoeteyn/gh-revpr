package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

const path = ".revpr/state.json"

type config struct {
	CurrentPR string `json:"currentPR"`
}

var currentConfig config

func init() {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	if err := json.Unmarshal(data, &currentConfig); err != nil {
		panic(err)
	}
}

func CurrentPR() string {
	return currentConfig.CurrentPR
}

func SetCurrentPR(pr string) error {
	currentConfig.CurrentPR = pr

	var file *os.File
	var fileErr error
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		slog.Debug("creating config file")
		fileErr = os.MkdirAll(filepath.Dir(path), os.ModePerm)
		if fileErr == nil {
			file, fileErr = os.Create(path)
		}
	} else {
		slog.Debug("opening config file")
		file, fileErr = os.Open(path)

	}

	if fileErr != nil {
		return fmt.Errorf("could not open file: %w", fileErr)
	}

	if err := json.NewEncoder(file).Encode(currentConfig); err != nil {
		return fmt.Errorf("could not encode config to file: %w", err)
	}

	return nil
}
