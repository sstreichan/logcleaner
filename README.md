# Log Cleaner

Terminal-basiertes Tool zum Säubern von Logdateien mit benutzerdefinierten Filtern.

## Features

- 🎨 **TUI (Terminal User Interface)** - Intuitive Bedienung mit Bubble Tea
- 🔍 **Custom Filter** - Regex-basierte Filter mit Validierung
- 💾 **Persistent Storage** - Filter werden automatisch gespeichert
- ⚡ **Auto-Completion** - Tab-Completion für Dateipfade
- 📦 **Auto-Release** - GitHub Actions für Versioning und Changelog
- 🚀 **Performance** - Streaming-basiert für große Logfiles

## Installation

### Von GitHub Releases

```bash
# Linux
wget https://github.com/sstreichan/logcleaner/releases/latest/download/logcleaner_Linux_x86_64.tar.gz
tar -xzf logcleaner_Linux_x86_64.tar.gz
sudo mv logcleaner /usr/local/bin/

# macOS
wget https://github.com/sstreichan/logcleaner/releases/latest/download/logcleaner_Darwin_x86_64.tar.gz
tar -xzf logcleaner_Darwin_x86_64.tar.gz
sudo mv logcleaner /usr/local/bin/

# Windows
# Download von https://github.com/sstreichan/logcleaner/releases/latest
```

### Aus dem Source

```bash
go install github.com/sstreichan/logcleaner/cmd/logcleaner@latest
```

## Usage

```bash
logcleaner
```

### Workflow

1. **Datei auswählen**: Pfad eingeben (Tab für Autocomplete)
2. **Filter verwalten**: Filter anzeigen, erstellen oder löschen
3. **Processing**: Logfile wird gefiltert
4. **Ergebnis**: Statistiken und Output-Datei

### Filter Syntax

**Remove Filter**: Entfernt Zeilen, die dem Pattern entsprechen
```json
{
  "name": "Remove Errors",
  "pattern": "^ERROR",
  "type": "remove"
}
```

**Keep Filter**: Behält nur Zeilen, die dem Pattern entsprechen
```json
{
  "name": "Keep Info",
  "pattern": "INFO|WARN",
  "type": "keep"
}
```

Pattern sind Go Regex (siehe [Syntax](https://pkg.go.dev/regexp/syntax)).

## Development

```bash
# Clone
git clone https://github.com/sstreichan/logcleaner.git
cd logcleaner

# Dependencies
go mod download

# Run
go run cmd/logcleaner/main.go

# Test
go test ./...
go test -v ./... -cover

# Build
go build -o logcleaner cmd/logcleaner/main.go
```

## Release Process

```bash
# 1. Commit changes with conventional commits
git commit -m "feat: add new filter type"
git commit -m "fix: autocomplete crash on empty input"

# 2. Create and push tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# 3. GitHub Actions automatically:
#    - Runs tests
#    - Builds binaries (Linux, macOS, Windows)
#    - Generates changelog
#    - Creates GitHub Release
```

## Architecture

```
cmd/logcleaner/        # Entry point
internal/
  ├── filter/          # Filter logic & validation
  ├── storage/         # Persistent filter storage
  ├── cleaner/         # Log processing engine
  └── tui/            # Bubble Tea UI components
```

## Configuration

Filter werden gespeichert in: `~/.config/logcleaner/filters.json`

## Contributing

Pull Requests sind willkommen! Bitte:
- Tests hinzufügen für neue Features
- Conventional Commits verwenden
- Code formatieren mit `go fmt`

## License

MIT License - siehe LICENSE Datei

## Roadmap

- [ ] Filter-Editor im TUI
- [ ] Multiple Filter kombinieren (AND/OR)
- [ ] Live-Preview während Filtering
- [ ] Export/Import von Filter-Sets
- [ ] Colored output für Logs
- [ ] Performance-Benchmarks
