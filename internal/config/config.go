package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type InputTarget struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Owner string `json:"owner"`
}

func LoadTargetsFromFile(path string) ([]InputTarget, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Impossible de lire le fichier %s : %w", path, err)
	}

	var targets []InputTarget
	if err := json.Unmarshal(data, &targets); err != nil {
		return nil, fmt.Errorf("Erreur lors de la désérialisation %s : %w", path, err)
	}

	return targets, nil
}

func SaveTargetsToFile(filePath string, targets []InputTarget) error {
	data, err := json.MarshalIndent(targets, "", "  ")
	if err != nil {
		return fmt.Errorf("impossible de lire le fichier %s: %w", filePath, err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("impossible de lire le fichier %s: %w", filePath, err)
	}
	return nil
}
