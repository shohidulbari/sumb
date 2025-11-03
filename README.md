# A simple and fast terminal utility to manage notes with full-text search

**sumb** is a lightweight command-line tool designed to help you quickly create, organize, and manage your notes directly from the terminal with full text search capabilities.

## 🧩 Installation

You can install **sumb** using the following one-liner:

```bash
curl -fsSL https://raw.githubusercontent.com/shohidulbari/sumb/main/install.sh | bash
```

Alternatively, you can clone the repository and build it from source:

```bash
chmod +x install.sh
./install.sh
```

## 🚀 Features

- **Add Notes**: Quickly add new notes with a terminal UI
- **Update Notes**: Edit existing notes effortlessly.
- **List Notes**: View your notes in a clean, organized format.
- **Search Notes**: Find notes using keywords with powerful full-text search capabilities.
- **View Notes**: Read your notes in a scrollable viewport.
- **Delete Notes**: Remove notes you no longer need with ease.

## 📚 Usage

After installation, you can use the following commands:

- `sumb create`: Opens a terminal textarea and allows you to type your note. `ctrl+s` will store the note, `ctrl+c` to cancel the creation.
- `sumb edit <note_id>`: Opens the specified note in a terminal textarea for editing.
- `sumb list <size:number>`: Displays a list of all your latest notes with their IDs. By default shows latest 10 notes.
- `sumb search <keyword:string>`: Searches for notes containing the specified keyword.
- `sumb show <note_id>`: Displays the full content of note in a scrollable viewport.
- `sumb delete <note_id>`: Deletes the note with the specified ID.

## 📂 Data Storage

**sumb** stores your notes in a local database using [bbolt](https://github.com/etcd-io/bbolt) and the database file is located at `~/.sumb/sumb.db` .
The full-text search index by [bleve](https://github.com/blevesearch/bleve) is stored at `~/.sumb/sumb.bleve`.

## 🤝 Contributing

Contributions are welcome!
If you find a bug or have a feature request, please open an issue on the [GitHub repository](https://github.com/shohidulbari/sumb) or submit a pull request with proper details.
