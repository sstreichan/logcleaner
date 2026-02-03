package tui

import (
	"os"
	"path/filepath"
	"strings"
)

type Autocomplete struct {
	cache        map[string][]string
	lastMatches  []string
	lastInput    string
	currentIndex int
}

func NewAutocomplete() *Autocomplete {
	return &Autocomplete{
		cache:        make(map[string][]string),
		lastMatches:  []string{},
		lastInput:    "",
		currentIndex: 0,
	}
}

// Complete returns the autocompleted path based on the current input
func (a *Autocomplete) Complete(input string) string {
	// Handle empty input - start with current directory
	if input == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return "./"
		}
		return cwd + string(filepath.Separator)
	}

	// Parse input into directory and base
	dir, base := a.parseInput(input)

	// Get all matching entries
	matches := a.findMatches(dir, base)

	if len(matches) == 0 {
		return input
	}

	if len(matches) == 1 {
		// Single match - return it and reset
		a.lastMatches = matches
		a.lastInput = input
		a.currentIndex = 0
		return matches[0]
	}

	// Multiple matches - cycle through them
	// If input changed, reset to first match
	if input != a.lastInput {
		a.lastMatches = matches
		a.lastInput = input
		a.currentIndex = 0
		
		// Try common prefix first
		commonPrefix := a.findCommonPrefix(matches)
		if commonPrefix != input {
			return commonPrefix
		}
		
		// No common prefix, return first match
		return matches[0]
	}

	// Same input as last time - cycle to next match
	a.currentIndex = (a.currentIndex + 1) % len(matches)
	return matches[a.currentIndex]
}

// parseInput splits input into directory and base name, handling special cases
func (a *Autocomplete) parseInput(input string) (string, string) {
	var dir, base string

	// Handle trailing separator (user is in a directory)
	if strings.HasSuffix(input, string(filepath.Separator)) {
		dir = input
		base = ""
	} else {
		dir = filepath.Dir(input)
		base = filepath.Base(input)
	}

	// Expand ~ to home directory
	if strings.HasPrefix(dir, "~") {
		home, err := os.UserHomeDir()
		if err == nil {
			if dir == "~" {
				dir = home
			} else {
				dir = strings.Replace(dir, "~", home, 1)
			}
		}
	}

	// Handle relative paths
	if dir == "." || dir == "" {
		cwd, err := os.Getwd()
		if err == nil {
			dir = cwd
		} else {
			dir = "."
		}
	}

	return dir, base
}

// findMatches returns all filesystem entries matching the base pattern in dir
func (a *Autocomplete) findMatches(dir, base string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return []string{}
	}

	var matches []string
	for _, entry := range entries {
		// Skip hidden files unless user explicitly typed a dot
		if !strings.HasPrefix(base, ".") && strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		// Check if entry matches the base
		if base == "" || strings.HasPrefix(entry.Name(), base) {
			fullPath := filepath.Join(dir, entry.Name())
			
			// Add trailing separator for directories
			if entry.IsDir() {
				fullPath += string(filepath.Separator)
			}
			
			matches = append(matches, fullPath)
		}
	}

	return matches
}

// findCommonPrefix returns the longest common prefix of all matches
func (a *Autocomplete) findCommonPrefix(matches []string) string {
	if len(matches) == 0 {
		return ""
	}

	if len(matches) == 1 {
		return matches[0]
	}

	// Find common prefix
	prefix := matches[0]
	for _, match := range matches[1:] {
		prefix = commonPrefix(prefix, match)
	}

	return prefix
}

// GetLastMatches returns the matches from the last completion attempt
// Useful for displaying suggestions to the user
func (a *Autocomplete) GetLastMatches() []string {
	return a.lastMatches
}

// GetCurrentIndex returns the currently selected match index
func (a *Autocomplete) GetCurrentIndex() int {
	return a.currentIndex
}

// Reset clears the autocomplete state
func (a *Autocomplete) Reset() {
	a.lastMatches = []string{}
	a.lastInput = ""
	a.currentIndex = 0
}

// commonPrefix returns the common prefix of two strings
func commonPrefix(a, b string) string {
	minLen := len(a)
	if len(b) < minLen {
		minLen = len(b)
	}

	for i := 0; i < minLen; i++ {
		if a[i] != b[i] {
			return a[:i]
		}
	}

	return a[:minLen]
}
