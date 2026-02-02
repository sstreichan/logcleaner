package cleaner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sstreichan/logcleaner/internal/filter"
)

func TestClean(t *testing.T) {
	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, "input.log")
	outputPath := filepath.Join(tempDir, "output.log")

	input := `ERROR: first error
INFO: some info
ERROR: second error
DEBUG: debug message
`

	if err := os.WriteFile(inputPath, []byte(input), 0644); err != nil {
		t.Fatal(err)
	}

	f, _ := filter.New("remove-errors", "^ERROR", filter.TypeRemove)
	c := New([]*filter.Filter{f})

	stats, err := c.Clean(inputPath, outputPath, nil)
	if err != nil {
		t.Fatalf("Clean() error = %v", err)
	}

	if stats.TotalLines != 4 {
		t.Errorf("Expected 4 total lines, got %d", stats.TotalLines)
	}

	if stats.FilteredLines != 2 {
		t.Errorf("Expected 2 filtered lines, got %d", stats.FilteredLines)
	}

	output, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(string(output), "ERROR") {
		t.Error("Output should not contain ERROR lines")
	}
}
