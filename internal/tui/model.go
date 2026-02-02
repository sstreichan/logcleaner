package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/sstreichan/logcleaner/internal/cleaner"
	"github.com/sstreichan/logcleaner/internal/filter"
	"github.com/sstreichan/logcleaner/internal/storage"
)

type screen int

const (
	screenFileSelect screen = iota
	screenFilterManage
	screenProcessing
	screenResults
)

type Model struct {
	screen       screen
	fileInput    textinput.Model
	filters      []*filter.Filter
	storage      *storage.Storage
	autocomplete *Autocomplete

	// Filter management
	filterList       list.Model
	newFilterName    textinput.Model
	newFilterPattern textinput.Model

	// Processing
	stats  *cleaner.Stats
	err    error
	width  int
	height int
}

func NewModel() (*Model, error) {
	storage, err := storage.New()
	if err != nil {
		return nil, err
	}

	filters, err := storage.Load()
	if err != nil {
		return nil, err
	}

	fileInput := textinput.New()
	fileInput.Placeholder = "Enter log file path..."
	fileInput.Focus()
	fileInput.CharLimit = 500
	fileInput.Width = 50

	return &Model{
		screen:       screenFileSelect,
		fileInput:    fileInput,
		filters:      filters,
		storage:      storage,
		autocomplete: NewAutocomplete(),
	}, nil
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "tab":
			if m.screen == screenFileSelect {
				completion := m.autocomplete.Complete(m.fileInput.Value())
				if completion != "" {
					m.fileInput.SetValue(completion)
				}
			}

		case "enter":
			if m.screen == screenFileSelect {
				if _, err := os.Stat(m.fileInput.Value()); err == nil {
					m.screen = screenFilterManage
				}
			}
		}
	}

	var cmd tea.Cmd
	m.fileInput, cmd = m.fileInput.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	switch m.screen {
	case screenFileSelect:
		return m.fileSelectView()
	case screenFilterManage:
		return m.filterManageView()
	case screenProcessing:
		return m.processingView()
	case screenResults:
		return m.resultsView()
	}
	return ""
}

func (m Model) fileSelectView() string {
	title := titleStyle.Render("Log File Cleaner")
	input := fmt.Sprintf("\n%s\n\n%s",
		subtitleStyle.Render("Enter the path to your log file:"),
		m.fileInput.View())

	help := helpStyle.Render("\nTab: autocomplete | Enter: continue | Ctrl+C: quit")

	return lipgloss.JoinVertical(lipgloss.Left, title, input, help)
}

func (m Model) filterManageView() string {
	var sb strings.Builder
	sb.WriteString(titleStyle.Render("Filter Management"))
	sb.WriteString("\n\n")

	if len(m.filters) == 0 {
		sb.WriteString(subtitleStyle.Render("No filters configured yet."))
	} else {
		sb.WriteString("Configured Filters:\n")
		for i, f := range m.filters {
			sb.WriteString(fmt.Sprintf("%d. %s (%s): %s\n", i+1, f.Name, f.Type, f.Pattern))
		}
	}

	sb.WriteString("\n")
	sb.WriteString(helpStyle.Render("Enter: process | Ctrl+C: quit"))

	return sb.String()
}

func (m Model) processingView() string {
	return "Processing (TODO)"
}

func (m Model) resultsView() string {
	return "Results (TODO)"
}
