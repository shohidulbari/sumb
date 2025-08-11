# sumb

A powerful CLI tool for managing notes, tasks, and pomodoros with ease. Built in Go with SQLite for local data storage.

## ✨ Features

- 📝 **Notes**: Create, list, search, and manage notes with JSON formatting support
- ✅ **Tasks**: Track tasks with deadlines, status management, and interactive creation
- ⏱️ **Pomodoros**: Time management with smart completion tracking and auto-stop
- 🔍 **Search**: Powerful search across all content with pagination
- 💾 **SQLite**: Local database for privacy and performance
- 🎯 **Interactive Mode**: User-friendly interactive creation for complex content

## 🚀 Installation

### Prerequisites

- Go 1.19 or later
- Git

### Option 1: Build from Source

```bash
# Clone the repository
git clone https://github.com/YOUR_USERNAME/sumb.git
cd sumb

# Build the binary
make build

# Install locally
make install
```

### Option 2: Quick Install Script

```bash
# Download and run the install script
curl -fsSL https://raw.githubusercontent.com/YOUR_USERNAME/sumb/main/install.sh | bash
```

### Option 3: Manual Build

```bash
# Clone and build manually
git clone https://github.com/YOUR_USERNAME/sumb.git
cd sumb
go build -o sumb cmd/sumb/main.go

# Move to your PATH
sudo mv sumb /usr/local/bin/
```

## 📖 Usage

### Notes Management

#### Quick Create

```bash
# Create a note with content
sumb note -c "Your note content here"

# Create a note with JSON content
sumb note -c '{"title": "Meeting Notes", "attendees": ["John", "Jane"]}'
```

#### Interactive Creation

```bash
# Start interactive note creation
sumb note -i

# Follow the prompts to enter your note content
# Press Ctrl+D when finished
```

#### List and Search

```bash
# List all notes (paginated)
sumb note list

# List with pagination
sumb note list --skip 10

# Search notes
sumb note search "meeting"

# Search with pagination
sumb note search "project" --skip 10

# Format JSON output
sumb note list --jsonify
sumb note search "query" --jsonify
```

#### Note Operations

```bash
# Delete a note
sumb note delete 5

# Delete multiple notes
sumb note delete-many 1 3 7
```

### Task Management

#### Quick Create

```bash
# Create a task with title
sumb task -c "Complete project documentation"

# Create with description
sumb task -c "Review code" -d "Check for bugs and style issues"

# Create with deadline
sumb task -c "Team meeting" -l "2025-01-20"

# Use date macros
sumb task -c "Daily standup" -l "today"
sumb task -c "Weekly review" -l "tomorrow"
```

#### Interactive Creation

```bash
# Start interactive task creation
sumb task -i

# Follow prompts for title, description, and deadline
# Accepts "today", "tomorrow", or YYYY-MM-DD format
```

#### List and Filter

```bash
# List all tasks
sumb task list

# List with pagination
sumb task list --skip 10

# Filter by deadline
sumb task list --today
sumb task list --tomorrow
```

#### Task Operations

```bash
# Update task status
sumb task status 5 IN_PROGRESS
sumb task status 5 COMPLETED

# Delete tasks
sumb task delete 3
sumb task delete-many 1 4 6
```

#### Search Tasks

```bash
# Search in title and description
sumb task search "meeting"

# Search with pagination
sumb task search "project" --skip 10
```

### Pomodoro Management

#### Start Pomodoro

```bash
# Start with title and duration
sumb pomodoro start -t "Work Session" -s 25

# Quick start
sumb pomodoro -c "Study Session" -s 30
```

#### Monitor and Control

```bash
# Check current status
sumb pomodoro status

# View live countdown
sumb pomodoro timer

# Stop current pomodoro
sumb pomodoro stop
```

#### List History

```bash
# View pomodoro history
sumb pomodoro list

# Paginated view
sumb pomodoro list --skip 10
```

## 🗄️ Data Storage

All data is stored locally in SQLite databases:

- **Notes**: `~/.sumb/sumb.db` (notes table)
- **Tasks**: `~/.sumb/sumb.db` (tasks table)
- **Pomodoros**: `~/.sumb/sumb.db` (pomodoros table)

## 🔧 Configuration

The tool automatically creates the necessary database and directory structure in your home directory (`~/.sumb/`).

## 📋 Status Values

### Tasks

- `TODO` - Task is pending
- `IN_PROGRESS` - Task is being worked on
- `COMPLETED` - Task is finished

### Pomodoros

- `ACTIVE` - Pomodoro is currently running
- `COMPLETED` - Pomodoro finished successfully
- `STOPPED` - Pomodoro was stopped manually

## 🎯 Examples

### Workflow Example

```bash
# 1. Start a pomodoro for focused work
sumb pomodoro start -t "Code Review" -s 25

# 2. Create a task for what you're working on
sumb task -c "Review pull request #123" -l "today"

# 3. Take notes during the session
sumb note -i
# Enter your notes interactively...

# 4. Check your progress
sumb task list --today
sumb pomodoro status

# 5. Complete the task when done
sumb task status 5 COMPLETED
```

### Study Session Example

```bash
# 1. Create study tasks
sumb task -c "Read Chapter 5" -l "today"
sumb task -c "Practice exercises" -l "tomorrow"

# 2. Start focused study time
sumb pomodoro start -t "Study Session" -s 30

# 3. Take study notes
sumb note -c "Chapter 5: Key concepts and formulas..."

# 4. Track progress
sumb task list --today
```

## 🚨 Troubleshooting

### Common Issues

**"No active pomodoro found"**

- This usually means no pomodoro is currently running
- Start a new one with `sumb pomodoro start -t "Title" -s 25`

**"Database error"**

- Check if `~/.sumb/` directory exists and is writable
- Try removing the directory and restarting (this will reset all data)

**"Command not found"**

- Ensure the binary is in your PATH
- Try running `which sumb` to locate the installation

### Reset Data

```bash
# Remove all data (use with caution!)
rm -rf ~/.sumb/
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) for CLI functionality
- Uses [SQLite](https://sqlite.org/) for data storage
- Inspired by productivity tools like Todo.txt and Pomodoro Technique

---

**Happy productivity! 🚀**
