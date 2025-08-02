# Sumb - Terminal Task Management

A fast, lightweight terminal-based task management application written in Go with SQLite storage.

## Features

- ✅ Create tasks with titles and descriptions
- 📋 List all tasks with status indicators
- ✅ Mark tasks as complete
- 🗑️ Delete tasks
- 💾 SQLite database for persistent storage
- 🚀 Fast and lightweight

## Installation

### From Source

1. Clone the repository:

```bash
git clone https://github.com/yourusername/sumb.git
cd sumb
```

2. Build the application:

```bash
make build
```

3. Install to your system:

```bash
make install
```

### Using Homebrew (Future)

```bash
brew install yourusername/sumb/sumb
```

## Usage

### Create a Task

```bash
# Create a task with title and description
sumb create -t "Buy groceries" -d "Milk, bread, eggs"

# Create a task with just title
sumb create -t "Call mom"

# Alternative format (as requested)
sumb -c "Buy groceries" -d "Milk, bread, eggs"
```

### List All Tasks

```bash
sumb list
```

### Mark Task as Complete

```bash
# Mark task with ID 1 as complete
sumb complete 1
```

### Delete a Task

```bash
# Delete task with ID 1
sumb delete 1
```

### Get Help

```bash
sumb --help
sumb create --help
```

## Database

The application stores data in a SQLite database located at:

- **macOS/Linux**: `~/.sumb/sumb.db`
- **Windows**: `%USERPROFILE%\.sumb\sumb.db`

## Development

### Prerequisites

- Go 1.21 or later
- SQLite3

### Building

```bash
# Build for current platform
make build

# Build for multiple platforms
make build-all

# Run tests
make test

# Clean build artifacts
make clean
```

### Project Structure

```
sumb/
├── cmd/           # CLI commands
│   └── tasks/     # Task management commands
├── tasks/         # Task management package
├── internal/      # Internal packages (legacy)
│   └── database/  # Database operations (legacy)
├── main.go        # Application entry point
├── go.mod         # Go module file
├── Makefile       # Build automation
└── README.md      # This file
```

## License

MIT License - see LICENSE file for details.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## Roadmap

- [ ] Add task categories/tags
- [ ] Add due dates
- [ ] Add priority levels
- [ ] Export tasks to various formats
- [ ] Add task search functionality
- [ ] Add task editing
- [ ] Add task templates
