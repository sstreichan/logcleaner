package cleaner

import (
	"bufio"
	"fmt"
	"os"

	"github.com/sstreichan/logcleaner/internal/filter"
)

type Cleaner struct {
	filters []*filter.Filter
}

func New(filters []*filter.Filter) *Cleaner {
	return &Cleaner{filters: filters}
}

type Stats struct {
	TotalLines    int
	FilteredLines int
	BytesRead     int64
}

func (c *Cleaner) Clean(inputPath, outputPath string, progressCb func(int, int)) (*Stats, error) {
	inFile, err := os.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open input file: %w", err)
	}
	defer inFile.Close()

	outFile, err := os.Create(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	stats := &Stats{}
	scanner := bufio.NewScanner(inFile)
	writer := bufio.NewWriter(outFile)
	defer writer.Flush()

	// Increase buffer size for large lines
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		stats.TotalLines++
		stats.BytesRead += int64(len(line))

		if progressCb != nil && lineNum%1000 == 0 {
			progressCb(lineNum, stats.FilteredLines)
		}

		shouldKeep := c.shouldKeepLine(line)
		if shouldKeep {
			if _, err := writer.WriteString(line + "\n"); err != nil {
				return stats, fmt.Errorf("failed to write line: %w", err)
			}
		} else {
			stats.FilteredLines++
		}
	}

	if err := scanner.Err(); err != nil {
		return stats, fmt.Errorf("error reading file: %w", err)
	}

	return stats, nil
}

func (c *Cleaner) shouldKeepLine(line string) bool {
	if len(c.filters) == 0 {
		return true
	}

	for _, f := range c.filters {
		matches := f.Matches(line)

		if f.Type == filter.TypeRemove && matches {
			return false
		}
		if f.Type == filter.TypeKeep && !matches {
			return false
		}
	}

	return true
}
