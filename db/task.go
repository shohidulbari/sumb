package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type TaskStatus string

const (
	StatusTODO        TaskStatus = "TODO"
	StatusInProgress  TaskStatus = "IN_PROGRESS"
	StatusComplete    TaskStatus = "COMPLETED"
)

type Task struct {
	ID          int
	Title       string
	Description string
	Status      TaskStatus
	Deadline    *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TaskManager struct {
	db *sql.DB
}

func NewTaskManager() (*TaskManager, error) {
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

	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return &TaskManager{db: db}, nil
}

func createTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		status TEXT DEFAULT 'TODO',
		deadline DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(query)
	return err
}

func (tm *TaskManager) Close() error {
	return tm.db.Close()
}

func (tm *TaskManager) CreateTask(title, description string, deadline *time.Time) error {
	query := `
	INSERT INTO tasks (title, description, status, deadline, created_at, updated_at)
	VALUES (?, ?, 'TODO', ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

	_, err := tm.db.Exec(query, title, description, deadline)
	return err
}

func (tm *TaskManager) ListTasks() ([]Task, error) {
	query := `
	SELECT id, title, description, status, deadline, created_at, updated_at
	FROM tasks
	ORDER BY created_at DESC`

	rows, err := tm.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		var deadline sql.NullTime
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &deadline, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, err
		}
		if deadline.Valid {
			task.Deadline = &deadline.Time
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (tm *TaskManager) ListTasksWithPagination(limit, offset int) ([]Task, error) {
	query := `
	SELECT id, title, description, status, deadline, created_at, updated_at
	FROM tasks
	ORDER BY created_at DESC
	LIMIT ? OFFSET ?`

	rows, err := tm.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		var deadline sql.NullTime
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &deadline, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, err
		}
		if deadline.Valid {
			task.Deadline = &deadline.Time
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (tm *TaskManager) ListTasksByDeadline(deadline time.Time) ([]Task, error) {
	// Normalize deadline to start of day for comparison
	deadlineStart := time.Date(deadline.Year(), deadline.Month(), deadline.Day(), 0, 0, 0, 0, deadline.Location())
	deadlineEnd := deadlineStart.AddDate(0, 0, 1)
	
	query := `
	SELECT id, title, description, status, deadline, created_at, updated_at
	FROM tasks
	WHERE deadline >= ? AND deadline < ?
	ORDER BY created_at DESC`

	rows, err := tm.db.Query(query, deadlineStart, deadlineEnd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		var deadline sql.NullTime
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &deadline, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, err
		}
		if deadline.Valid {
			task.Deadline = &deadline.Time
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (tm *TaskManager) UpdateTaskStatus(id int, status TaskStatus) error {
	query := `
	UPDATE tasks 
	SET status = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?`

	result, err := tm.db.Exec(query, status, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task with ID %d not found", id)
	}

	return nil
}

func (tm *TaskManager) SearchTasks(query string, limit, offset int) ([]Task, error) {
	searchQuery := `
	SELECT id, title, description, status, deadline, created_at, updated_at
	FROM tasks
	WHERE title LIKE ? OR description LIKE ?
	ORDER BY created_at DESC
	LIMIT ? OFFSET ?`

	// Use %query% for partial matching in both title and description
	searchPattern := "%" + query + "%"
	
	rows, err := tm.db.Query(searchQuery, searchPattern, searchPattern, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		var deadline sql.NullTime
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &deadline, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, err
		}
		if deadline.Valid {
			task.Deadline = &deadline.Time
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (tm *TaskManager) DeleteTask(id int) error {
	query := `DELETE FROM tasks WHERE id = ?`

	result, err := tm.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task with ID %d not found", id)
	}
	return nil
} 