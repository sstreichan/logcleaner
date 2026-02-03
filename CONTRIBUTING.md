# Contributing to Log Cleaner

Vielen Dank für dein Interesse an Log Cleaner! 🎉

## Development Setup

### Prerequisites

- Go 1.22 or higher
- Git
- Make (optional, aber empfohlen)

### Setup

```bash
# Clone repository
git clone https://github.com/sstreichan/logcleaner.git
cd logcleaner

# Install dependencies
go mod download

# Run tests
make test

# Run the app
make run
```

## Development Workflow

### 1. Create a Branch

```bash
git checkout -b feature/your-feature-name
# oder
git checkout -b fix/your-bug-fix
```

### 2. Make Changes

- Schreibe sauberen, idiomatischen Go Code
- Füge Tests für neue Features hinzu
- Halte dich an die bestehende Code-Struktur

### 3. Test Your Changes

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run benchmarks (for performance changes)
make bench

# Format code
make fmt

# Vet code
make vet
```

### 4. Commit Your Changes

Wir verwenden [Conventional Commits](https://www.conventionalcommits.org/):

```bash
# Features
git commit -m "feat: add new filter type for JSON logs"

# Bug fixes
git commit -m "fix: prevent crash on empty log files"

# Documentation
git commit -m "docs: update README with new examples"

# Performance
git commit -m "perf: optimize regex compilation"

# Refactoring
git commit -m "refactor: simplify filter validation logic"

# Tests
git commit -m "test: add edge cases for autocomplete"
```

### 5. Push and Create PR

```bash
git push origin feature/your-feature-name
```

Dann erstelle einen Pull Request auf GitHub.

## Code Guidelines

### Go Style

- Folge [Effective Go](https://golang.org/doc/effective_go.html)
- Verwende `gofmt` für Formatierung
- Schreibe selbsterklärenden Code mit minimalen Kommentaren
- Halte Funktionen klein und fokussiert

### Testing

- Jede neue Funktion braucht Tests
- Strebe >80% Code Coverage an
- Teste Edge Cases und Fehler-Szenarien
- Benchmarks für Performance-kritischen Code

### TUI Development

- Teste UI-Changes manuell in verschiedenen Terminal-Größen
- Stelle sicher, dass Keyboard-Navigation funktioniert
- Verwende konsistente Styles aus `internal/tui/styles.go`

## Project Structure

```
logcleaner/
├── cmd/logcleaner/       # Application entry point
├── internal/
│   ├── filter/           # Filter logic and validation
│   ├── storage/          # Persistent storage
│   ├── cleaner/          # Log processing engine
│   └── tui/             # Terminal UI (Bubble Tea)
├── examples/            # Example files and filters
└── .github/workflows/   # CI/CD
```

## Common Tasks

### Adding a New Filter Type

1. Update `internal/filter/filter.go`
2. Add validation logic
3. Write tests in `internal/filter/filter_test.go`
4. Update TUI if needed
5. Add example to `examples/filters/`

### Adding a New TUI Screen

1. Add screen constant in `internal/tui/model.go`
2. Implement view function
3. Add update logic for keyboard navigation
4. Add styles to `internal/tui/styles.go` if needed

### Performance Optimization

1. Write benchmark first (`*_benchmark_test.go`)
2. Make changes
3. Compare benchmarks:
   ```bash
   go test -bench=. -benchmem ./internal/cleaner/ > old.txt
   # make changes
   go test -bench=. -benchmem ./internal/cleaner/ > new.txt
   benchcmp old.txt new.txt
   ```

## Questions?

Erstelle ein [Issue](https://github.com/sstreichan/logcleaner/issues) oder frage in der PR!
