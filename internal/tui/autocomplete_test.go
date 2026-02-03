package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestComplete_EmptyInput(t *testing.T) {
	a := NewAutocomplete()
	result := a.Complete("")

	if result == "" {
		t.Error("Expected non-empty result for empty input")
	}
}

func TestComplete_SingleMatch(t *testing.T) {
	tempDir := t.TempDir()

	// Create a unique file
	testFile := filepath.Join(tempDir, "unique_test_file.log")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	a := NewAutocomplete()
	input := filepath.Join(tempDir, "unique")
	result := a.Complete(input)

	if result != testFile {
		t.Errorf("Expected %s, got %s", testFile, result)
	}
}

func TestComplete_MultipleMatches(t *testing.T) {
	tempDir := t.TempDir()

	// Create multiple files with common prefix
	os.WriteFile(filepath.Join(tempDir, "test1.log"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(tempDir, "test2.log"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(tempDir, "test3.log"), []byte("test"), 0644)

	a := NewAutocomplete()
	input := filepath.Join(tempDir, "test")
	result := a.Complete(input)

	// Should return common prefix
	expected := filepath.Join(tempDir, "test")
	if !strings.HasPrefix(result, expected) {
		t.Errorf("Expected result to start with %s, got %s", expected, result)
	}

	// Check that matches were stored
	matches := a.GetLastMatches()
	if len(matches) != 3 {
		t.Errorf("Expected 3 matches, got %d", len(matches))
	}
}

func TestComplete_Directory(t *testing.T) {
	tempDir := t.TempDir()

	// Create a subdirectory
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	a := NewAutocomplete()
	input := filepath.Join(tempDir, "sub")
	result := a.Complete(input)

	// Should return directory with trailing separator
	expected := subDir + string(filepath.Separator)
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestComplete_NoMatch(t *testing.T) {
	tempDir := t.TempDir()

	a := NewAutocomplete()
	input := filepath.Join(tempDir, "nonexistent")
	result := a.Complete(input)

	// Should return original input when no matches
	if result != input {
		t.Errorf("Expected %s, got %s", input, result)
	}
}

func TestCommonPrefix(t *testing.T) {
	tests := []struct {
		a, b, expected string
	}{
		{"test1.log", "test2.log", "test"},
		{"file.txt", "file.log", "file."},
		{"abc", "xyz", ""},
		{"same", "same", "same"},
	}

	for _, tt := range tests {
		result := commonPrefix(tt.a, tt.b)
		if result != tt.expected {
			t.Errorf("commonPrefix(%s, %s) = %s, want %s", tt.a, tt.b, result, tt.expected)
		}
	}
}
