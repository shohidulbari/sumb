package db

import (
	"os"
	"path/filepath"
	"testing"
)

// TestDB represents a test database setup
type TestDB struct {
	TempDir    string
	OriginalHome string
	TaskManager *TaskManager
}

// NewTestDB creates a new test database with proper isolation
func NewTestDB(t *testing.T) *TestDB {
	// Create a temporary directory for this test
	tempDir := t.TempDir()
	
	// Store original HOME environment
	originalHome := os.Getenv("HOME")
	
	// Set HOME to temp directory for this test
	os.Setenv("HOME", tempDir)
	
	// Create the .sumb directory in temp
	dbDir := filepath.Join(tempDir, ".sumb")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		t.Fatalf("Failed to create test database directory: %v", err)
	}
	
	// Initialize task manager with test database
	tm, err := NewTaskManager()
	if err != nil {
		t.Fatalf("Failed to initialize test task manager: %v", err)
	}
	
	return &TestDB{
		TempDir:      tempDir,
		OriginalHome: originalHome,
		TaskManager:  tm,
	}
}

// Close cleans up the test database
func (tdb *TestDB) Close(t *testing.T) {
	// Close the task manager
	if tdb.TaskManager != nil {
		if err := tdb.TaskManager.Close(); err != nil {
			t.Errorf("Failed to close task manager: %v", err)
		}
	}
	
	// Restore original HOME environment
	os.Setenv("HOME", tdb.OriginalHome)
	
	// The temp directory will be automatically cleaned up by t.TempDir()
}

// GetDBPath returns the path to the test database file
func (tdb *TestDB) GetDBPath() string {
	return filepath.Join(tdb.TempDir, ".sumb", "sumb.db")
}

// CleanupDB removes the database file for manual cleanup if needed
func (tdb *TestDB) CleanupDB(t *testing.T) {
	dbPath := tdb.GetDBPath()
	if err := os.Remove(dbPath); err != nil && !os.IsNotExist(err) {
		t.Errorf("Failed to cleanup test database: %v", err)
	}
} 