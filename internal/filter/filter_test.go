package filter

import (
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name       string
		filterName string
		pattern    string
		filterType FilterType
		wantError  bool
	}{
		{"valid remove filter", "test", "^ERROR", TypeRemove, false},
		{"valid keep filter", "test", "INFO", TypeKeep, false},
		{"empty name", "", "pattern", TypeRemove, true},
		{"empty pattern", "test", "", TypeRemove, true},
		{"invalid regex", "test", "[", TypeRemove, true},
		{"invalid type", "test", "pattern", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.filterName, tt.pattern, tt.filterType)
			if (err != nil) != tt.wantError {
				t.Errorf("New() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestFilterMatches(t *testing.T) {
	f, _ := New("test", "^ERROR", TypeRemove)

	if !f.Matches("ERROR: something went wrong") {
		t.Error("Expected match for ERROR line")
	}

	if f.Matches("INFO: all good") {
		t.Error("Expected no match for INFO line")
	}
}
