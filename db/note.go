package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Note struct {
	ID        int
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type NoteManager struct {
	db *sql.DB
}

func NewNoteManager() (*NoteManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	dbDir := filepath.Join(homeDir, ".sumb")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	dbPath := filepath.Join(dbDir, "sumb.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := createNoteTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return &NoteManager{db: db}, nil
}

func createNoteTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS notes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		body TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(query)
	return err
}

func (nm *NoteManager) Close() error {
	return nm.db.Close()
}

func (nm *NoteManager) CreateNote(body string) error {
	query := `
	INSERT INTO notes (body, created_at, updated_at)
	VALUES (?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

	_, err := nm.db.Exec(query, body)
	return err
}

func (nm *NoteManager) ListNotes() ([]Note, error) {
	query := `
	SELECT id, body, created_at, updated_at
	FROM notes
	ORDER BY created_at DESC`

	rows, err := nm.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var note Note
		err := rows.Scan(&note.ID, &note.Body, &note.CreatedAt, &note.UpdatedAt)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func (nm *NoteManager) DeleteNote(id int) error {
	query := `DELETE FROM notes WHERE id = ?`

	result, err := nm.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("note with ID %d not found", id)
	}

	return nil
}

func (nm *NoteManager) ListNotesWithPagination(limit, offset int) ([]Note, error) {
	query := `
	SELECT id, body, created_at, updated_at
	FROM notes
	ORDER BY created_at DESC
	LIMIT ? OFFSET ?`

	rows, err := nm.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var note Note
		err := rows.Scan(&note.ID, &note.Body, &note.CreatedAt, &note.UpdatedAt)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	return notes, nil
} 