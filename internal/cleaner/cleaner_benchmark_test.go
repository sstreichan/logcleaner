package cleaner

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/sstreichan/logcleaner/internal/filter"
)

func BenchmarkClean_NoFilters(b *testing.B) {
	tempDir := b.TempDir()
	inputPath := filepath.Join(tempDir, "input.log")
	outputPath := filepath.Join(tempDir, "output.log")

	// Create a log file with 10000 lines
	input := generateLogLines(10000)
	if err := os.WriteFile(inputPath, []byte(input), 0644); err != nil {
		b.Fatal(err)
	}

	c := New([]*filter.Filter{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := c.Clean(inputPath, outputPath, nil)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkClean_SingleFilter(b *testing.B) {
	tempDir := b.TempDir()
	inputPath := filepath.Join(tempDir, "input.log")
	outputPath := filepath.Join(tempDir, "output.log")

	input := generateLogLines(10000)
	if err := os.WriteFile(inputPath, []byte(input), 0644); err != nil {
		b.Fatal(err)
	}

	f, _ := filter.New("remove-errors", "^ERROR", filter.TypeRemove)
	c := New([]*filter.Filter{f})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := c.Clean(inputPath, outputPath, nil)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkClean_MultipleFilters(b *testing.B) {
	tempDir := b.TempDir()
	inputPath := filepath.Join(tempDir, "input.log")
	outputPath := filepath.Join(tempDir, "output.log")

	input := generateLogLines(10000)
	if err := os.WriteFile(inputPath, []byte(input), 0644); err != nil {
		b.Fatal(err)
	}

	f1, _ := filter.New("remove-errors", "^ERROR", filter.TypeRemove)
	f2, _ := filter.New("remove-debug", "^DEBUG", filter.TypeRemove)
	f3, _ := filter.New("remove-trace", "^TRACE", filter.TypeRemove)
	c := New([]*filter.Filter{f1, f2, f3})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := c.Clean(inputPath, outputPath, nil)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func generateLogLines(count int) string {
	levels := []string{"INFO", "DEBUG", "WARN", "ERROR", "TRACE"}
	var result string

	for i := 0; i < count; i++ {
		level := levels[i%len(levels)]
		result += fmt.Sprintf("%s: This is log line number %d with some additional text\n", level, i)
	}

	return result
}
