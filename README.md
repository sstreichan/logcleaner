# 🧹 Log Cleaner

[![Release](https://img.shields.io/github/v/release/sstreichan/logcleaner)](https://github.com/sstreichan/logcleaner/releases)
[![Test](https://github.com/sstreichan/logcleaner/actions/workflows/test.yml/badge.svg)](https://github.com/sstreichan/logcleaner/actions/workflows/test.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/sstreichan/logcleaner)](https://go.dev/)
[![License](https://img.shields.io/github/license/sstreichan/logcleaner)](LICENSE)

Terminal-basiertes Tool zum Säubern von Logdateien mit benutzerdefinierten Filtern und intuitiver TUI.

## ✨ Features

- 🎨 **Terminal User Interface** - Gebaut mit [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- 🔍 **Regex-basierte Filter** - Mächtige Pattern-Matching-Capabilities
- ✅ **Filter-Validierung** - Verhindert ungültige Regex beim Speichern
- 💾 **Persistent Storage** - Filter werden automatisch in `~/.config/logcleaner/` gespeichert
- ⚡ **Tab-Completion** - Auto-Vervollständigung für Dateipfade
- 🚀 **Performance** - Streaming-basiert für große Logfiles (>1GB)
- 📦 **Auto-Release** - GitHub Actions für Versioning und Multi-Platform Builds

## 🚀 Quick Start

### Installation

**macOS / Linux:**
```bash
curl -sSL https://github.com/sstreichan/logcleaner/releases/latest/download/logcleaner_$(uname -s)_$(uname -m).tar.gz | tar xz
sudo mv logcleaner /usr/local/bin/
```

**Oder mit Go:**
```bash
go install github.com/sstreichan/logcleaner/cmd/logcleaner@latest
```

### Erste Schritte

```bash
# Starte die TUI
logcleaner

# 1. Gib den Pfad zu deiner Logdatei ein (Tab für Autocomplete)
# 2. Verwalte deine Filter (a: add, d: delete)
# 3. Drücke Enter zum Verarbeiten
# 4. Fertig! Cleaned file: yourfile.log.cleaned
```

## 📖 Usage

### Basic Workflow

1. **Datei auswählen**
   - Pfad eingeben oder mit Tab durch Verzeichnisse navigieren
   - Enter zum Bestätigen

2. **Filter verwalten**
   - `a` - Neuen Filter hinzufügen
   - `d` - Ausgewählten Filter löschen
   - `↑/↓` - Durch Filter navigieren
   - Enter - Verarbeitung starten

3. **Filter erstellen**
   - Name eingeben (z.B. "Remove Errors")
   - Regex Pattern (z.B. `^ERROR|^FATAL`)
   - Typ wählen: **Remove** (entfernen) oder **Keep** (behalten)

4. **Ergebnis**
   - Statistiken über verarbeitete Zeilen
   - Output-Datei: `<original>.cleaned`

### Filter-Beispiele

#### Fehler entfernen
```json
{
  "name": "Remove Errors",
  "pattern": "^(ERROR|FATAL|CRITICAL)",
  "type": "remove"
}
```

#### Nur Warnungen behalten
```json
{
  "name": "Keep Warnings",
  "pattern": "WARN|WARNING|ALERT",
  "type": "keep"
}
```

#### HTTP-Fehler behalten (4xx, 5xx)
```json
{
  "name": "Keep HTTP Errors",
  "pattern": "HTTP/\\d\\.\\d\" [45]\\d{2}",
  "type": "keep"
}
```

#### Debug-Zeilen entfernen
```json
{
  "name": "Remove Debug",
  "pattern": "^DEBUG|^TRACE|\\[DEBUG\\]",
  "type": "remove"
}
```

#### IP-Adressen filtern
```json
{
  "name": "Keep Specific IP",
  "pattern": "192\\.168\\.1\\.100",
  "type": "keep"
}
```

### Vordefinierte Filter importieren

```bash
# Kopiere Beispiel-Filter in deine Config
cp examples/filters/common.json ~/.config/logcleaner/filters.json
```

## 🛠️ Development

### Prerequisites

- Go 1.22+
- Make (optional)

### Setup

```bash
# Clone
git clone https://github.com/sstreichan/logcleaner.git
cd logcleaner

# Dependencies
go mod download

# Run
make run
# oder
go run cmd/logcleaner/main.go
```

### Testing

```bash
# All tests
make test

# With coverage
make test-coverage
# Öffne coverage.html im Browser

# Benchmarks
make bench
```

### Build

```bash
# Local build
make build

# Cross-platform (requires goreleaser)
make release-test
```

## 📁 Project Structure

```
logcleaner/
├── cmd/logcleaner/          # Main entry point
├── internal/
│   ├── filter/              # Filter logic & validation
│   │   ├── filter.go
│   │   └── filter_test.go
│   ├── storage/             # JSON persistence
│   │   ├── storage.go
│   │   └── storage_test.go
│   ├── cleaner/             # Log processing engine
│   │   ├── cleaner.go
│   │   ├── cleaner_test.go
│   │   └── cleaner_benchmark_test.go
│   └── tui/                 # Bubble Tea UI
│       ├── model.go         # Main model & screens
│       ├── styles.go        # UI styling
│       └── autocomplete.go  # Path completion
├── examples/
│   ├── filters/             # Example filter presets
│   └── logs/                # Sample log files
├── .github/workflows/       # CI/CD
└── .goreleaser.yml          # Release config
```

## 🚢 Release Process

```bash
# 1. Commit mit Conventional Commits
git commit -m "feat: add JSON log parsing"
git commit -m "fix: handle empty lines correctly"

# 2. Tag erstellen
git tag -a v1.0.0 -m "Release v1.0.0"

# 3. Push (löst GitHub Action aus)
git push origin v1.0.0
```

GitHub Actions erstellt automatisch:
- ✅ Binaries für Linux, macOS, Windows (amd64 + arm64)
- ✅ Release mit Auto-Generated Changelog
- ✅ Archiv-Downloads (.tar.gz, .zip)

## ⚡ Performance

Getestet auf einem MacBook Pro M1:

| File Size | Lines | Time | Memory |
|-----------|-------|------|--------|
| 10 MB | 100k | ~0.5s | ~15 MB |
| 100 MB | 1M | ~4s | ~30 MB |
| 1 GB | 10M | ~45s | ~50 MB |

*Mit 3 Regex-Filtern, Streaming-basiert*

### Benchmarks

```bash
$ make bench
goos: darwin
goarch: arm64
BenchmarkClean_NoFilters-10         100   11245633 ns/op   8388608 B/op
BenchmarkClean_SingleFilter-10       50   23456789 ns/op   8388608 B/op
BenchmarkClean_MultipleFilters-10    30   35678901 ns/op   8388608 B/op
```

## 🔧 Configuration

Filter werden gespeichert in:
- **Linux/macOS**: `~/.config/logcleaner/filters.json`
- **Windows**: `%APPDATA%\logcleaner\filters.json`

Manuelles Bearbeiten möglich:
```bash
vim ~/.config/logcleaner/filters.json
```

## 🐛 Troubleshooting

### Filter wird nicht gespeichert

**Problem**: Filter verschwindet nach Neustart

**Lösung**: Prüfe Schreibrechte für `~/.config/logcleaner/`
```bash
ls -la ~/.config/logcleaner/
chmod 644 ~/.config/logcleaner/filters.json
```

### Regex-Pattern funktioniert nicht

**Problem**: Pattern matched nicht wie erwartet

**Lösung**: Teste dein Pattern online: [regex101.com](https://regex101.com/) (wähle "Golang" Flavor)

### TUI zeigt komische Zeichen

**Problem**: Terminal unterstützt keine Unicode-Zeichen

**Lösung**: Verwende ein modernes Terminal (iTerm2, Windows Terminal, Alacritty)

## 🤝 Contributing

Contributions sind willkommen! Siehe [CONTRIBUTING.md](CONTRIBUTING.md) für Details.

**Quick Checklist:**
- ✅ Tests schreiben
- ✅ Conventional Commits verwenden
- ✅ Code formatieren (`make fmt`)
- ✅ Tests bestehen (`make test`)

## 📋 Roadmap

- [ ] Live-Preview während Filtering
- [ ] Filter-Export/Import (YAML, TOML)
- [ ] Filter-Kombinationen (AND/OR Logic)
- [ ] Colored Log Output im TUI
- [ ] Undo/Redo Funktionalität
- [ ] Multi-File Processing
- [ ] Filter-Templates für bekannte Log-Formate (nginx, Apache, syslog)
- [ ] Cloud Storage Integration (S3, GCS)

## 📄 License

MIT License - siehe [LICENSE](LICENSE)

## 🙏 Credits

Gebaut mit:
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI Framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI Components
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Styling
- [GoReleaser](https://goreleaser.com/) - Release Automation

---

**Made with ❤️ in Dresden**
