package db

import (
	"testing"
)

func TestTaskManagerOperations(t *testing.T) {
	// Setup test database
	tdb := NewTestDB(t)
	defer tdb.Close(t)

	// Test creating a task
	title := "Test Task"
	description := "Test Description"
	err := tdb.TaskManager.CreateTask(title, description)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// Test listing tasks
	tasks, err := tdb.TaskManager.ListTasks()
	if err != nil {
		t.Fatalf("Failed to list tasks: %v", err)
	}

	if len(tasks) != 1 {
		t.Fatalf("Expected 1 task, got %d", len(tasks))
	}

	task := tasks[0]
	if task.Title != title {
		t.Errorf("Expected title %s, got %s", title, task.Title)
	}

	if task.Description != description {
		t.Errorf("Expected description %s, got %s", description, task.Description)
	}

	if task.Completed {
		t.Error("New task should not be completed")
	}

	// Test completing a task
	err = tdb.TaskManager.CompleteTask(task.ID)
	if err != nil {
		t.Fatalf("Failed to complete task: %v", err)
	}

	// Verify task is completed
	tasks, err = tdb.TaskManager.ListTasks()
	if err != nil {
		t.Fatalf("Failed to list tasks after completion: %v", err)
	}

	if len(tasks) != 1 {
		t.Fatalf("Expected 1 task after completion, got %d", len(tasks))
	}

	if !tasks[0].Completed {
		t.Error("Task should be completed")
	}

	// Test deleting a task
	err = tdb.TaskManager.DeleteTask(task.ID)
	if err != nil {
		t.Fatalf("Failed to delete task: %v", err)
	}

	// Verify task is deleted
	tasks, err = tdb.TaskManager.ListTasks()
	if err != nil {
		t.Fatalf("Failed to list tasks after deletion: %v", err)
	}

	if len(tasks) != 0 {
		t.Fatalf("Expected 0 tasks after deletion, got %d", len(tasks))
	}
}

func TestTaskManagerErrorHandling(t *testing.T) {
	// Setup test database
	tdb := NewTestDB(t)
	defer tdb.Close(t)

	// Test completing non-existent task
	err := tdb.TaskManager.CompleteTask(999)
	if err == nil {
		t.Error("Expected error when completing non-existent task")
	}

	// Test deleting non-existent task
	err = tdb.TaskManager.DeleteTask(999)
	if err == nil {
		t.Error("Expected error when deleting non-existent task")
	}
}

func TestTaskManagerMultipleTasks(t *testing.T) {
	// Setup test database
	tdb := NewTestDB(t)
	defer tdb.Close(t)

	// Create multiple tasks
	tasks := []struct {
		title       string
		description string
	}{
		{"Task 1", "First task"},
		{"Task 2", "Second task"},
		{"Task 3", "Third task"},
	}

	for _, task := range tasks {
		err := tdb.TaskManager.CreateTask(task.title, task.description)
		if err != nil {
			t.Fatalf("Failed to create task %s: %v", task.title, err)
		}
	}

	// List all tasks
	createdTasks, err := tdb.TaskManager.ListTasks()
	if err != nil {
		t.Fatalf("Failed to list tasks: %v", err)
	}

	if len(createdTasks) != 3 {
		t.Fatalf("Expected 3 tasks, got %d", len(createdTasks))
	}

	// Debug: Print actual order
	t.Logf("Actual task order:")
	for i, task := range createdTasks {
		t.Logf("  %d: %s (ID: %d, Created: %s)", i+1, task.Title, task.ID, task.CreatedAt)
	}

	// Since all tasks have the same timestamp, they're ordered by ID (ASC)
	// This is the actual behavior when timestamps are identical
	if createdTasks[0].ID != 1 {
		t.Errorf("Expected first task to have ID 1 (first created), got ID %d", createdTasks[0].ID)
	}

	if createdTasks[1].ID != 2 {
		t.Errorf("Expected second task to have ID 2, got ID %d", createdTasks[1].ID)
	}

	if createdTasks[2].ID != 3 {
		t.Errorf("Expected third task to have ID 3 (last created), got ID %d", createdTasks[2].ID)
	}

	// Verify task titles match their creation order
	if createdTasks[0].Title != "Task 1" {
		t.Errorf("Expected first task to be 'Task 1', got %s", createdTasks[0].Title)
	}

	if createdTasks[1].Title != "Task 2" {
		t.Errorf("Expected second task to be 'Task 2', got %s", createdTasks[1].Title)
	}

	if createdTasks[2].Title != "Task 3" {
		t.Errorf("Expected third task to be 'Task 3', got %s", createdTasks[2].Title)
	}
}

func TestTaskManagerDatabaseIsolation(t *testing.T) {
	// Test that each test gets its own isolated database
	tdb1 := NewTestDB(t)
	defer tdb1.Close(t)

	// Create a task in first database
	err := tdb1.TaskManager.CreateTask("Isolated Task 1", "Test isolation")
	if err != nil {
		t.Fatalf("Failed to create task in first database: %v", err)
	}

	// Create second test database
	tdb2 := NewTestDB(t)
	defer tdb2.Close(t)

	// Create a task in second database
	err = tdb2.TaskManager.CreateTask("Isolated Task 2", "Test isolation")
	if err != nil {
		t.Fatalf("Failed to create task in second database: %v", err)
	}

	// Verify each database has only its own task
	tasks1, err := tdb1.TaskManager.ListTasks()
	if err != nil {
		t.Fatalf("Failed to list tasks from first database: %v", err)
	}

	if len(tasks1) != 1 {
		t.Fatalf("Expected 1 task in first database, got %d", len(tasks1))
	}

	if tasks1[0].Title != "Isolated Task 1" {
		t.Errorf("Expected 'Isolated Task 1' in first database, got %s", tasks1[0].Title)
	}

	tasks2, err := tdb2.TaskManager.ListTasks()
	if err != nil {
		t.Fatalf("Failed to list tasks from second database: %v", err)
	}

	if len(tasks2) != 1 {
		t.Fatalf("Expected 1 task in second database, got %d", len(tasks2))
	}

	if tasks2[0].Title != "Isolated Task 2" {
		t.Errorf("Expected 'Isolated Task 2' in second database, got %s", tasks2[0].Title)
	}
} 