package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Pomodoro struct {
	ID          int
	Title       string
	Duration    int
	StartedAt   time.Time
	CompletedAt *time.Time
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type PomodoroManager struct {
	db *sql.DB
}

func NewPomodoroManager() (*PomodoroManager, error) {
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

	if err := createPomodoroTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	if err := migratePomodoroStatuses(db); err != nil {
		return nil, fmt.Errorf("failed to migrate pomodoro statuses: %w", err)
	}

	return &PomodoroManager{db: db}, nil
}

func createPomodoroTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS pomodoros (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		duration INTEGER NOT NULL,
		started_at DATETIME NOT NULL,
		completed_at DATETIME,
		status TEXT NOT NULL DEFAULT 'ACTIVE',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(query)
	return err
}

func migratePomodoroStatuses(db *sql.DB) error {
	queries := []string{
		`UPDATE pomodoros SET status = 'ACTIVE' WHERE status = 'active'`,
		`UPDATE pomodoros SET status = 'COMPLETED' WHERE status = 'completed'`,
		`UPDATE pomodoros SET status = 'STOPPED' WHERE status = 'stopped'`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute migration query '%s': %w", query, err)
		}
	}

	return nil
}

func (pm *PomodoroManager) Close() error {
	return pm.db.Close()
}

func (pm *PomodoroManager) StopAllActivePomodoros() error {
	query := `
	UPDATE pomodoros 
	SET status = 'STOPPED', updated_at = CURRENT_TIMESTAMP
	WHERE status = 'ACTIVE'`

	_, err := pm.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (pm *PomodoroManager) CreatePomodoro(title string, duration int) (*Pomodoro, error) {
	query := `
	INSERT INTO pomodoros (title, duration, started_at, status, created_at, updated_at)
	VALUES (?, ?, CURRENT_TIMESTAMP, 'ACTIVE', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

	result, err := pm.db.Exec(query, title, duration)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Pomodoro{
		ID:        int(id),
		Title:     title,
		Duration:  duration,
		StartedAt: time.Now(),
		Status:    "ACTIVE",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (pm *PomodoroManager) GetActivePomodoro() (*Pomodoro, error) {
	query := `
	SELECT id, title, duration, started_at, completed_at, status, created_at, updated_at
	FROM pomodoros
	WHERE status = 'ACTIVE'
	ORDER BY created_at DESC
	LIMIT 1`

	var pomodoro Pomodoro
	var completedAt sql.NullTime

	err := pm.db.QueryRow(query).Scan(
		&pomodoro.ID,
		&pomodoro.Title,
		&pomodoro.Duration,
		&pomodoro.StartedAt,
		&completedAt,
		&pomodoro.Status,
		&pomodoro.CreatedAt,
		&pomodoro.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if completedAt.Valid {
		pomodoro.CompletedAt = &completedAt.Time
	}

	return &pomodoro, nil
}

func (pm *PomodoroManager) CompletePomodoro(id int) error {
	query := `
	UPDATE pomodoros 
	SET status = 'COMPLETED', completed_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?`

	result, err := pm.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("pomodoro with ID %d not found", id)
	}

	return nil
}

func (pm *PomodoroManager) StopPomodoro(id int) error {
	query := `
	UPDATE pomodoros 
	SET status = 'STOPPED', updated_at = CURRENT_TIMESTAMP
	WHERE id = ?`

	result, err := pm.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("pomodoro with ID %d not found", id)
	}

	return nil
}

func (pm *PomodoroManager) ListPomodoros(limit, offset int) ([]Pomodoro, error) {
	query := `
	SELECT id, title, duration, started_at, completed_at, status, created_at, updated_at
	FROM pomodoros
	ORDER BY created_at DESC
	LIMIT ? OFFSET ?`

	rows, err := pm.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pomodoros []Pomodoro
	for rows.Next() {
		var pomodoro Pomodoro
		var completedAt sql.NullTime

		err := rows.Scan(
			&pomodoro.ID,
			&pomodoro.Title,
			&pomodoro.Duration,
			&pomodoro.StartedAt,
			&completedAt,
			&pomodoro.Status,
			&pomodoro.CreatedAt,
			&pomodoro.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if completedAt.Valid {
			pomodoro.CompletedAt = &completedAt.Time
		}

		pomodoros = append(pomodoros, pomodoro)
	}

	return pomodoros, nil
} 