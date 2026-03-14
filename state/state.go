package state

import (
	"encoding/json"
	"os"
)

const path = ".revpr"

type config struct {
	CurrentPR string `json:"currentPR"`
}

var currentConfig config

func init() {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	json.Unmarshal(data, &currentConfig)
}

func CurrentPR() string {
	return currentConfig.CurrentPR
}
