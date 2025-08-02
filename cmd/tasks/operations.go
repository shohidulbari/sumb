package tasks

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	db "github.com/sumb/db"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new task",
	Long:  `Create a new task with a title and optional description.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		title, _ := cmd.Flags().GetString("title")
		description, _ := cmd.Flags().GetString("description")

		if title == "" {
			return fmt.Errorf("title is required")
		}

		tm, err := db.NewTaskManager()
		if err != nil {
			return fmt.Errorf("failed to initialize task manager: %w", err)
		}
		defer tm.Close()

		if err := tm.CreateTask(title, description); err != nil {
			return fmt.Errorf("failed to create task: %w", err)
		}

		fmt.Printf("✅ Task created!\n")
		fmt.Printf("Title: %s\n", title)
		if description != "" {
			fmt.Printf("Description: %s\n", description)
		}

		return nil
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Long:  `Display tasks in the database with their status and details. Shows max 10 latest tasks by default. Use --skip to paginate through older tasks.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		skip, _ := cmd.Flags().GetInt("skip")
		
		tm, err := db.NewTaskManager()
		if err != nil {
			return fmt.Errorf("failed to initialize task manager: %w", err)
		}
		defer tm.Close()

		taskList, err := tm.ListTasksWithPagination(10, skip)
		if err != nil {
			return fmt.Errorf("failed to list tasks: %w", err)
		}

		if len(taskList) == 0 {
			if skip > 0 {
				fmt.Printf("📝 No more tasks found. You've reached the end of your task list.\n")
			} else {
				fmt.Println("📝 No tasks found. Create your first task with 'sumb task create -t \"Your Task\"'")
			}
			return nil
		}

		fmt.Printf("📋 Found %d task(s)", len(taskList))
		if skip > 0 {
			fmt.Printf(" (skipped %d)", skip)
		}
		fmt.Printf(":\n\n")
		
		for _, task := range taskList {
			status := "⏳"
			if task.Completed {
				status = "✅"
			}

			fmt.Printf("%s [%d] %s\n", status, task.ID, task.Title)
			if task.Description != "" {
				fmt.Printf("   Description: %s\n", task.Description)
			}
			fmt.Printf("   Created: %s\n", task.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Println()
		}

		// Show pagination hint if there might be more tasks
		if len(taskList) == 10 {
			fmt.Printf("💡 To see more tasks, use: sumb task list --skip %d\n", skip+10)
		}

		return nil
	},
}

var completeCmd = &cobra.Command{
	Use:   "complete",
	Short: "Mark a task as complete",
	Long:  `Mark a task as complete by providing its ID.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("task ID is required")
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid task ID: %s", args[0])
		}

		tm, err := db.NewTaskManager()
		if err != nil {
			return fmt.Errorf("failed to initialize task manager: %w", err)
		}
		defer tm.Close()

		if err := tm.CompleteTask(id); err != nil {
			return fmt.Errorf("failed to complete task: %w", err)
		}

		fmt.Printf("✅ Task %d marked as complete!\n", id)
		return nil
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a task",
	Long:  `Delete a task by providing its ID.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("task ID is required")
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid task ID: %s", args[0])
		}

		tm, err := db.NewTaskManager()
		if err != nil {
			return fmt.Errorf("failed to initialize task manager: %w", err)
		}
		defer tm.Close()

		if err := tm.DeleteTask(id); err != nil {
			return fmt.Errorf("failed to delete task: %w", err)
		}

		fmt.Printf("🗑️  Task %d deleted successfully!\n", id)
		return nil
	},
}

func init() {
	createCmd.Flags().StringP("title", "t", "", "Task title (required)")
	createCmd.Flags().StringP("description", "d", "", "Task description")
	createCmd.MarkFlagRequired("title")
	
	// Add pagination flag for list command
	listCmd.Flags().IntP("skip", "s", 0, "Number of tasks to skip (for pagination)")
} 