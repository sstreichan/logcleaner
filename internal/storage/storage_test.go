package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sstreichan/logcleaner/internal/filter"
)

func TestSaveAndLoad(t *testing.T) {
	tempDir := t.TempDir()
	s := &Storage{
		configPath: filepath.Join(tempDir, "filters.json"),
	}

	filters := []*filter.Filter{
		{Name: "test1", Pattern: "^ERROR", Type: filter.TypeRemove},
		{Name: "test2", Pattern: "INFO", Type: filter.TypeKeep},
	}

	if err := s.Save(filters); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	loaded, err := s.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if len(loaded) != len(filters) {
		t.Errorf("Expected %d filters, got %d", len(filters), len(loaded))
	}
}

func TestLoadNonExistent(t *testing.T) {
	tempDir := t.TempDir()
	s := &Storage{
		configPath: filepath.Join(tempDir, "nonexistent.json"),
	}

	filters, err := s.Load()
	if err != nil {
		t.Fatalf("Load() should not error on nonexistent file: %v", err)
	}

	if len(filters) != 0 {
		t.Errorf("Expected empty filter list, got %d", len(filters))
	}
}
