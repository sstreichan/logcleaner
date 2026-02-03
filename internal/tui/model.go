package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

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
	screenFilterAdd
	screenProcessing
	screenResults
)

type processingMsg struct {
	stats *cleaner.Stats
	err   error
}

type Model struct {
	screen       screen
	fileInput    textinput.Model
	filePath     string
	filters      []*filter.Filter
	storage      *storage.Storage
	autocomplete *Autocomplete

	// Filter management
	selectedFilter   int
	newFilterName    textinput.Model
	newFilterPattern textinput.Model
	newFilterType    filter.FilterType
	filterInputFocus int // 0=name, 1=pattern, 2=type

	// Processing
	processing       bool
	progressLines    int
	progressFiltered int
	stats            *cleaner.Stats
	err              error

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
	fileInput.Width = 60

	newFilterName := textinput.New()
	newFilterName.Placeholder = "Filter name"
	newFilterName.Width = 40

	newFilterPattern := textinput.New()
	newFilterPattern.Placeholder = "Regex pattern (e.g. ^ERROR)"
	newFilterPattern.Width = 40

	return &Model{
		screen:           screenFileSelect,
		fileInput:        fileInput,
		filters:          filters,
		storage:          storage,
		autocomplete:     NewAutocomplete(),
		newFilterName:    newFilterName,
		newFilterPattern: newFilterPattern,
		newFilterType:    filter.TypeRemove,
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

	case processingMsg:
		m.processing = false
		m.stats = msg.stats
		m.err = msg.err
		m.screen = screenResults
		return m, nil

	case tea.KeyMsg:
		// Global quit keys
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		// Route to appropriate screen handler
		switch m.screen {
		case screenFileSelect:
			return m.updateFileSelect(msg)
		case screenFilterManage:
			return m.updateFilterManage(msg)
		case screenFilterAdd:
			return m.updateFilterAdd(msg)
		case screenResults:
			return m.updateResults(msg)
		}
	}

	return m, nil
}

func (m Model) updateFileSelect(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		return m, tea.Quit

	case "tab":
		// Handle autocomplete cycling
		currentValue := m.fileInput.Value()
		completion := m.autocomplete.Complete(currentValue)
		
		if completion != "" && completion != currentValue {
			m.fileInput.SetValue(completion)
			m.fileInput.SetCursor(len(completion))
		}
		return m, nil

	case "enter":
		// Validate and move to next screen
		if m.fileInput.Value() != "" {
			if _, err := os.Stat(m.fileInput.Value()); err == nil {
				m.filePath = m.fileInput.Value()
				m.screen = screenFilterManage
				m.autocomplete.Reset()
			}
		}
		return m, nil

	default:
		// For all other keys, let textinput handle them
		// But first, check if this is a typing key (not just navigation)
		oldValue := m.fileInput.Value()
		
		var cmd tea.Cmd
		m.fileInput, cmd = m.fileInput.Update(msg)
		
		// If the value changed, reset autocomplete
		if m.fileInput.Value() != oldValue {
			m.autocomplete.Reset()
		}
		
		return m, cmd
	}
}

func (m Model) updateFilterManage(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		return m, tea.Quit

	case "esc":
		m.screen = screenFileSelect
		return m, nil

	case "up", "k":
		if m.selectedFilter > 0 {
			m.selectedFilter--
		}

	case "down", "j":
		if m.selectedFilter < len(m.filters)-1 {
			m.selectedFilter++
		}

	case "a":
		m.screen = screenFilterAdd
		m.newFilterName.SetValue("")
		m.newFilterPattern.SetValue("")
		m.newFilterName.Focus()
		m.newFilterPattern.Blur()
		m.filterInputFocus = 0
		m.newFilterType = filter.TypeRemove
		return m, textinput.Blink

	case "d":
		if len(m.filters) > 0 && m.selectedFilter < len(m.filters) {
			m.filters = append(m.filters[:m.selectedFilter], m.filters[m.selectedFilter+1:]...)
			if m.selectedFilter >= len(m.filters) && m.selectedFilter > 0 {
				m.selectedFilter--
			}
			m.storage.Save(m.filters)
		}

	case "enter":
		if m.filePath != "" {
			m.screen = screenProcessing
			m.processing = true
			return m, m.processFile()
		}
	}

	return m, nil
}

func (m Model) updateFilterAdd(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		// In filter add screen, 'q' should type 'q', not quit
		// Only ctrl+c quits (handled globally)
		break

	case "esc":
		m.screen = screenFilterManage
		return m, nil

	case "tab", "shift+tab":
		if msg.String() == "tab" {
			m.filterInputFocus = (m.filterInputFocus + 1) % 3
		} else {
			m.filterInputFocus = (m.filterInputFocus + 2) % 3
		}

		if m.filterInputFocus == 0 {
			m.newFilterName.Focus()
			m.newFilterPattern.Blur()
		} else if m.filterInputFocus == 1 {
			m.newFilterName.Blur()
			m.newFilterPattern.Focus()
		} else {
			m.newFilterName.Blur()
			m.newFilterPattern.Blur()
		}
		return m, textinput.Blink

	case "left", "right":
		if m.filterInputFocus == 2 {
			if m.newFilterType == filter.TypeRemove {
				m.newFilterType = filter.TypeKeep
			} else {
				m.newFilterType = filter.TypeRemove
			}
			return m, nil
		}
		// If not on type selector, let textinput handle left/right

	case "enter":
		name := strings.TrimSpace(m.newFilterName.Value())
		pattern := strings.TrimSpace(m.newFilterPattern.Value())

		if name != "" && pattern != "" {
			newFilter, err := filter.New(name, pattern, m.newFilterType)
			if err == nil {
				m.filters = append(m.filters, newFilter)
				m.storage.Save(m.filters)
				m.screen = screenFilterManage
				return m, nil
			}
		}
		return m, nil
	}

	// Let the focused input handle the key
	var cmd tea.Cmd
	if m.filterInputFocus == 0 {
		m.newFilterName, cmd = m.newFilterName.Update(msg)
	} else if m.filterInputFocus == 1 {
		m.newFilterPattern, cmd = m.newFilterPattern.Update(msg)
	}
	return m, cmd
}

func (m Model) updateResults(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		return m, tea.Quit

	case "enter", "esc":
		m.screen = screenFileSelect
		m.fileInput.SetValue("")
		m.stats = nil
		m.err = nil
		m.autocomplete.Reset()
		return m, nil
	}

	return m, nil
}

func (m Model) processFile() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(100 * time.Millisecond) // Small delay for UI

		outputPath := m.filePath + ".cleaned"
		c := cleaner.New(m.filters)

		stats, err := c.Clean(m.filePath, outputPath, func(lines, filtered int) {
			// Progress callback (could be enhanced with tea.Cmd)
			m.progressLines = lines
			m.progressFiltered = filtered
		})

		return processingMsg{stats: stats, err: err}
	}
}

func (m Model) View() string {
	switch m.screen {
	case screenFileSelect:
		return m.fileSelectView()
	case screenFilterManage:
		return m.filterManageView()
	case screenFilterAdd:
		return m.filterAddView()
	case screenProcessing:
		return m.processingView()
	case screenResults:
		return m.resultsView()
	}
	return ""
}

func (m Model) fileSelectView() string {
	title := titleStyle.Render("📝 Log File Cleaner")

	var content strings.Builder
	content.WriteString("\n")
	content.WriteString(subtitleStyle.Render("Enter the path to your log file:"))
	content.WriteString("\n\n")
	content.WriteString(m.fileInput.View())
	content.WriteString("\n\n")

	// Show file validation
	if m.fileInput.Value() != "" {
		if _, err := os.Stat(m.fileInput.Value()); err != nil {
			content.WriteString(errorStyle.Render("⚠ File not found"))
			content.WriteString("\n")
		} else {
			content.WriteString(infoStyle.Render("✓ File exists"))
			content.WriteString("\n")
		}
	}

	// Show autocomplete suggestions
	matches := m.autocomplete.GetLastMatches()
	if len(matches) > 0 {
		content.WriteString("\n")
		
		// Show how many matches there are
		if len(matches) == 1 {
			content.WriteString(dimStyle.Render("1 match:"))
		} else {
			currentIdx := m.autocomplete.GetCurrentIndex()
			content.WriteString(dimStyle.Render(fmt.Sprintf("%d matches (showing %d/%d):", len(matches), currentIdx+1, len(matches))))
		}
		content.WriteString("\n")
		
		// Display up to 10 suggestions
		maxDisplay := 10
		currentIdx := m.autocomplete.GetCurrentIndex()
		
		for i, match := range matches {
			if i >= maxDisplay {
				content.WriteString(dimStyle.Render(fmt.Sprintf("  ... and %d more", len(matches)-maxDisplay)))
				content.WriteString("\n")
				break
			}
			
			// Extract proper display name
			displayName := match
			
			// For directories with trailing separator, remove it to get the name
			if strings.HasSuffix(match, string(filepath.Separator)) {
				// Remove trailing separator
				cleanPath := strings.TrimSuffix(match, string(filepath.Separator))
				// Get the directory name
				displayName = filepath.Base(cleanPath) + "/"
			} else {
				// For files, just use the base name
				displayName = filepath.Base(match)
			}
			
			// If displayName is still empty or just "/", show the full path
			if displayName == "" || displayName == "/" || displayName == "./" {
				displayName = match
			}
			
			// Highlight the currently selected match
			if i == currentIdx {
				content.WriteString(selectedItemStyle.Render(fmt.Sprintf("→ %s", displayName)))
			} else {
				content.WriteString(dimStyle.Render(fmt.Sprintf("  %s", displayName)))
			}
			content.WriteString("\n")
		}
	}

	help := helpStyle.Render("Tab: cycle completions | Enter: continue | Ctrl+C: quit")

	return lipgloss.JoinVertical(lipgloss.Left, title, content.String(), help)
}

func (m Model) filterManageView() string {
	var sb strings.Builder

	sb.WriteString(titleStyle.Render("🔧 Filter Management"))
	sb.WriteString("\n\n")
	sb.WriteString(infoStyle.Render(fmt.Sprintf("File: %s", filepath.Base(m.filePath))))
	sb.WriteString("\n\n")

	if len(m.filters) == 0 {
		sb.WriteString(subtitleStyle.Render("No filters configured yet."))
		sb.WriteString("\n")
		sb.WriteString(dimStyle.Render("Press 'a' to add your first filter."))
	} else {
		sb.WriteString(subtitleStyle.Render("Active Filters:"))
		sb.WriteString("\n\n")

		for i, f := range m.filters {
			var style lipgloss.Style
			var prefix string

			if i == m.selectedFilter {
				style = selectedItemStyle
				prefix = "→ "
			} else {
				style = itemStyle
				prefix = "  "
			}

			var typeIcon string
			if f.Type == filter.TypeRemove {
				typeIcon = "🗑️  Remove"
			} else {
				typeIcon = "✅ Keep"
			}

			line := fmt.Sprintf("%s%s [%s]: %s", prefix, f.Name, typeIcon, f.Pattern)
			sb.WriteString(style.Render(line))
			sb.WriteString("\n")
		}
	}

	sb.WriteString("\n")
	sb.WriteString(helpStyle.Render("↑/↓: navigate | a: add filter | d: delete | Enter: process | Esc: back | Ctrl+C: quit"))

	return sb.String()
}

func (m Model) filterAddView() string {
	var sb strings.Builder

	sb.WriteString(titleStyle.Render("➕ Add New Filter"))
	sb.WriteString("\n\n")

	// Name input
	var nameLabel string
	if m.filterInputFocus == 0 {
		nameLabel = focusedLabelStyle.Render("Name:")
	} else {
		nameLabel = labelStyle.Render("Name:")
	}
	sb.WriteString(nameLabel)
	sb.WriteString("\n")
	sb.WriteString(m.newFilterName.View())
	sb.WriteString("\n\n")

	// Pattern input
	var patternLabel string
	if m.filterInputFocus == 1 {
		patternLabel = focusedLabelStyle.Render("Pattern (Regex):")
	} else {
		patternLabel = labelStyle.Render("Pattern (Regex):")
	}
	sb.WriteString(patternLabel)
	sb.WriteString("\n")
	sb.WriteString(m.newFilterPattern.View())
	sb.WriteString("\n\n")

	// Type selector
	var typeLabel string
	if m.filterInputFocus == 2 {
		typeLabel = focusedLabelStyle.Render("Type:")
	} else {
		typeLabel = labelStyle.Render("Type:")
	}
	sb.WriteString(typeLabel)
	sb.WriteString("\n")

	if m.newFilterType == filter.TypeRemove {
		sb.WriteString(selectedButtonStyle.Render("[Remove]"))
		sb.WriteString(" ")
		sb.WriteString(buttonStyle.Render(" Keep "))
	} else {
		sb.WriteString(buttonStyle.Render("Remove"))
		sb.WriteString(" ")
		sb.WriteString(selectedButtonStyle.Render("[Keep]"))
	}

	sb.WriteString("\n\n")
	sb.WriteString(dimStyle.Render("Remove: Filter out matching lines | Keep: Only keep matching lines"))
	sb.WriteString("\n\n")
	sb.WriteString(helpStyle.Render("Tab: next field | ←/→: toggle type | Enter: save | Esc: cancel"))

	return sb.String()
}

func (m Model) processingView() string {
	var sb strings.Builder

	sb.WriteString(titleStyle.Render("⚙️  Processing"))
	sb.WriteString("\n\n")
	sb.WriteString(infoStyle.Render("Cleaning log file..."))
	sb.WriteString("\n\n")

	spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	frame := int(time.Now().UnixMilli()/100) % len(spinner)
	sb.WriteString(spinner[frame])
	sb.WriteString(" Processing...")

	return sb.String()
}

func (m Model) resultsView() string {
	var sb strings.Builder

	if m.err != nil {
		sb.WriteString(titleStyle.Render("❌ Error"))
		sb.WriteString("\n\n")
		sb.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
	} else if m.stats != nil {
		sb.WriteString(titleStyle.Render("✅ Complete"))
		sb.WriteString("\n\n")

		outputPath := m.filePath + ".cleaned"

		// Statistics box
		statsBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1, 2).
			Width(50)

		statsContent := fmt.Sprintf(
			"📊 Statistics\n\n"+
				"Total Lines:     %d\n"+
				"Filtered Lines:  %d\n"+
				"Remaining Lines: %d\n"+
				"Bytes Processed: %.2f MB\n\n"+
				"Output: %s",
			m.stats.TotalLines,
			m.stats.FilteredLines,
			m.stats.TotalLines-m.stats.FilteredLines,
			float64(m.stats.BytesRead)/(1024*1024),
			filepath.Base(outputPath),
		)

		sb.WriteString(statsBox.Render(statsContent))
	}

	sb.WriteString("\n\n")
	sb.WriteString(helpStyle.Render("Enter: process another file | Ctrl+C: quit"))

	return sb.String()
}
