package tui

import (
	"os"
	"path/filepath"
	"strings"
)

type Autocomplete struct {
	cache map[string][]string
}

func NewAutocomplete() *Autocomplete {
	return &Autocomplete{
		cache: make(map[string][]string),
	}
}

func (a *Autocomplete) Complete(input string) string {
	if input == "" {
		return ""
	}

	dir := filepath.Dir(input)
	base := filepath.Base(input)

	// Expand ~ to home directory
	if strings.HasPrefix(dir, "~") {
		home, _ := os.UserHomeDir()
		dir = strings.Replace(dir, "~", home, 1)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return input
	}

	var matches []string
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), base) {
			fullPath := filepath.Join(dir, entry.Name())
			if entry.IsDir() {
				fullPath += string(filepath.Separator)
			}
			matches = append(matches, fullPath)
		}
	}

	if len(matches) == 1 {
		return matches[0]
	}

	return input
}
