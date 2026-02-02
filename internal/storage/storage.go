package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sstreichan/logcleaner/internal/filter"
)

type Storage struct {
	configPath string
}

func New() (*Storage, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".config", "logcleaner")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	return &Storage{
		configPath: filepath.Join(configDir, "filters.json"),
	}, nil
}

func (s *Storage) Load() ([]*filter.Filter, error) {
	data, err := os.ReadFile(s.configPath)
	if os.IsNotExist(err) {
		return []*filter.Filter{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read filters: %w", err)
	}

	var filters []*filter.Filter
	if err := json.Unmarshal(data, &filters); err != nil {
		return nil, fmt.Errorf("failed to parse filters: %w", err)
	}

	return filters, nil
}

func (s *Storage) Save(filters []*filter.Filter) error {
	data, err := json.MarshalIndent(filters, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal filters: %w", err)
	}

	if err := os.WriteFile(s.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write filters: %w", err)
	}

	return nil
}
