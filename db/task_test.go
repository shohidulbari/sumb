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
	err := tdb.TaskManager.CreateTask(title, description, nil, StatusTODO)
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

	if task.Status != StatusTODO {
		t.Error("New task should have TODO status")
	}

	// Test updating task status
	err = tdb.TaskManager.UpdateTaskStatus(task.ID, StatusComplete)
	if err != nil {
		t.Fatalf("Failed to update task status: %v", err)
	}

	// Verify task status is updated
	tasks, err = tdb.TaskManager.ListTasks()
	if err != nil {
		t.Fatalf("Failed to list tasks after status update: %v", err)
	}

	if len(tasks) != 1 {
		t.Fatalf("Expected 1 task after status update, got %d", len(tasks))
	}

	if tasks[0].Status != StatusComplete {
		t.Error("Task should have COMPLETE status")
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

	// Test updating status of non-existent task
	err := tdb.TaskManager.UpdateTaskStatus(999, StatusComplete)
	if err == nil {
		t.Error("Expected error when updating status of non-existent task")
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
		err := tdb.TaskManager.CreateTask(task.title, task.description, nil, StatusTODO)
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
	err := tdb1.TaskManager.CreateTask("Isolated Task 1", "Test isolation", nil, StatusTODO)
	if err != nil {
		t.Fatalf("Failed to create task in first database: %v", err)
	}

	// Create second test database
	tdb2 := NewTestDB(t)
	defer tdb2.Close(t)

	// Create a task in second database
	err = tdb2.TaskManager.CreateTask("Isolated Task 2", "Test isolation", nil, StatusTODO)
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

func TestTaskManagerPagination(t *testing.T) {
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
		{"Task 4", "Fourth task"},
		{"Task 5", "Fifth task"},
		{"Task 6", "Sixth task"},
		{"Task 7", "Seventh task"},
		{"Task 8", "Eighth task"},
		{"Task 9", "Ninth task"},
		{"Task 10", "Tenth task"},
		{"Task 11", "Eleventh task"},
		{"Task 12", "Twelfth task"},
	}

	for _, task := range tasks {
		err := tdb.TaskManager.CreateTask(task.title, task.description, nil, StatusTODO)
		if err != nil {
			t.Fatalf("Failed to create task %s: %v", task.title, err)
		}
	}

	// Test first page (limit 5, offset 0)
	firstPage, err := tdb.TaskManager.ListTasksWithPagination(5, 0)
	if err != nil {
		t.Fatalf("Failed to get first page: %v", err)
	}

	if len(firstPage) != 5 {
		t.Fatalf("Expected 5 tasks on first page, got %d", len(firstPage))
	}

	// Since all tasks have the same timestamp, they're ordered by ID (ASC)
	// First page should have tasks with IDs 1-5
	expectedIDs := []int{1, 2, 3, 4, 5}
	for i, task := range firstPage {
		if task.ID != expectedIDs[i] {
			t.Errorf("Expected task %d to have ID %d, got ID %d", i+1, expectedIDs[i], task.ID)
		}
	}

	// Test second page (limit 5, offset 5)
	secondPage, err := tdb.TaskManager.ListTasksWithPagination(5, 5)
	if err != nil {
		t.Fatalf("Failed to get second page: %v", err)
	}

	if len(secondPage) != 5 {
		t.Fatalf("Expected 5 tasks on second page, got %d", len(secondPage))
	}

	// Second page should have tasks with IDs 6-10
	expectedIDs2 := []int{6, 7, 8, 9, 10}
	for i, task := range secondPage {
		if task.ID != expectedIDs2[i] {
			t.Errorf("Expected task %d to have ID %d, got ID %d", i+1, expectedIDs2[i], task.ID)
		}
	}

	// Test third page (limit 5, offset 10)
	thirdPage, err := tdb.TaskManager.ListTasksWithPagination(5, 10)
	if err != nil {
		t.Fatalf("Failed to get third page: %v", err)
	}

	if len(thirdPage) != 2 {
		t.Fatalf("Expected 2 tasks on third page, got %d", len(thirdPage))
	}

	// Third page should have tasks with IDs 11-12
	expectedIDs3 := []int{11, 12}
	for i, task := range thirdPage {
		if task.ID != expectedIDs3[i] {
			t.Errorf("Expected task %d to have ID %d, got ID %d", i+1, expectedIDs3[i], task.ID)
		}
	}

	// Test empty page (limit 5, offset 15)
	emptyPage, err := tdb.TaskManager.ListTasksWithPagination(5, 15)
	if err != nil {
		t.Fatalf("Failed to get empty page: %v", err)
	}

	if len(emptyPage) != 0 {
		t.Fatalf("Expected 0 tasks on empty page, got %d", len(emptyPage))
	}
}

func TestTaskManagerDeleteMultiple(t *testing.T) {
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
		{"Task 4", "Fourth task"},
		{"Task 5", "Fifth task"},
	}

	for _, task := range tasks {
		err := tdb.TaskManager.CreateTask(task.title, task.description, nil, StatusTODO)
		if err != nil {
			t.Fatalf("Failed to create task %s: %v", task.title, err)
		}
	}

	// Verify all tasks exist
	allTasks, err := tdb.TaskManager.ListTasks()
	if err != nil {
		t.Fatalf("Failed to list tasks: %v", err)
	}

	if len(allTasks) != 5 {
		t.Fatalf("Expected 5 tasks, got %d", len(allTasks))
	}

	// Test deleting multiple tasks (IDs 2, 3, 4)
	idsToDelete := []int{2, 3, 4}
	for _, id := range idsToDelete {
		err := tdb.TaskManager.DeleteTask(id)
		if err != nil {
			t.Fatalf("Failed to delete task %d: %v", id, err)
		}
	}

	// Verify tasks were deleted
	remainingTasks, err := tdb.TaskManager.ListTasks()
	if err != nil {
		t.Fatalf("Failed to list remaining tasks: %v", err)
	}

	if len(remainingTasks) != 2 {
		t.Fatalf("Expected 2 remaining tasks, got %d", len(remainingTasks))
	}

	// Verify specific tasks remain (IDs 1 and 5)
	expectedIDs := []int{1, 5}
	for i, task := range remainingTasks {
		if task.ID != expectedIDs[i] {
			t.Errorf("Expected remaining task %d to have ID %d, got ID %d", i+1, expectedIDs[i], task.ID)
		}
	}

	// Test deleting non-existent task
	err = tdb.TaskManager.DeleteTask(999)
	if err == nil {
		t.Error("Expected error when deleting non-existent task, got nil")
	}
} 