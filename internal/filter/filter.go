package filter

import (
	"encoding/json"
	"fmt"
	"regexp"
)

type FilterType string

const (
	TypeRemove FilterType = "remove"
	TypeKeep   FilterType = "keep"
)

type Filter struct {
	Name    string     `json:"name"`
	Pattern string     `json:"pattern"`
	Type    FilterType `json:"type"`
	regex   *regexp.Regexp
}

func New(name, pattern string, filterType FilterType) (*Filter, error) {
	if name == "" {
		return nil, fmt.Errorf("filter name cannot be empty")
	}
	if pattern == "" {
		return nil, fmt.Errorf("filter pattern cannot be empty")
	}

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern: %w", err)
	}

	if filterType != TypeRemove && filterType != TypeKeep {
		return nil, fmt.Errorf("invalid filter type: must be 'remove' or 'keep'")
	}

	return &Filter{
		Name:    name,
		Pattern: pattern,
		Type:    filterType,
		regex:   regex,
	}, nil
}

func (f *Filter) Matches(line string) bool {
	if f.regex == nil {
		f.regex, _ = regexp.Compile(f.Pattern)
	}
	return f.regex.MatchString(line)
}

func (f *Filter) MarshalJSON() ([]byte, error) {
	type Alias Filter
	return json.Marshal(&struct{ *Alias }{(*Alias)(f)})
}

func (f *Filter) UnmarshalJSON(data []byte) error {
	type Alias Filter
	aux := &struct{ *Alias }{Alias: (*Alias)(f)}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	regex, err := regexp.Compile(f.Pattern)
	if err != nil {
		return fmt.Errorf("invalid regex in stored filter: %w", err)
	}
	f.regex = regex
	return nil
}
