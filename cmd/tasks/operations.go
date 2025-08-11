package tasks

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/shohidulbari/sumb/cmd/styles"
	db "github.com/shohidulbari/sumb/db"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new task",
	Long:  `Create a new task with a title, optional description, and optional deadline.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		title, _ := cmd.Flags().GetString("title")
		description, _ := cmd.Flags().GetString("description")
		deadlineStr, _ := cmd.Flags().GetString("deadline")

		if title == "" {
			return fmt.Errorf("title is required")
		}

		deadline, err := parseDeadline(deadlineStr)
		if err != nil {
			return fmt.Errorf("invalid deadline: %w", err)
		}

		tm, err := db.NewTaskManager()
		if err != nil {
			return fmt.Errorf("failed to initialize task manager: %w", err)
		}
		defer tm.Close()

		if err := tm.CreateTask(title, description, deadline); err != nil {
			return fmt.Errorf("failed to create task: %w", err)
		}

		fmt.Printf(styles.TaskCreated)
		fmt.Printf("Title: %s\n", title)
		if description != "" {
			fmt.Printf("Description: %s\n", description)
		}
		if deadline != nil {
			fmt.Printf("Deadline: %s\n", deadline.Format("Jan 2, 2006"))
		}

		return nil
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Long:  `View up to 10 latest tasks with status and details. Use –skip to see older tasks.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		skip, _ := cmd.Flags().GetInt("skip")
		today, _ := cmd.Flags().GetBool("today")
		tomorrow, _ := cmd.Flags().GetBool("tomorrow")
		
		tm, err := db.NewTaskManager()
		if err != nil {
			return fmt.Errorf("failed to initialize task manager: %w", err)
		}
		defer tm.Close()

		var taskList []db.Task
		var err2 error

		if today {
			taskList, err2 = tm.ListTasksByDeadline(time.Now())
		} else if tomorrow {
			taskList, err2 = tm.ListTasksByDeadline(time.Now().AddDate(0, 0, 1))
		} else {
			taskList, err2 = tm.ListTasksWithPagination(10, skip)
		}

		if err2 != nil {
			return fmt.Errorf("failed to list tasks: %w", err2)
		}

		if len(taskList) == 0 {
			if today {
				fmt.Println("No tasks due today.")
			} else if tomorrow {
				fmt.Println("No tasks due tomorrow.")
			} else if skip > 0 {
				fmt.Printf("No more tasks found. You've reached the end of your task list.\n")
			} else {
				fmt.Println("No tasks found. Create your first task with 'sumb task create -t \"Your Task\"'")
			}
			return nil
		}

		fmt.Println()
		if today {
			fmt.Printf("Tasks due today (%d):\n\n", len(taskList))
		} else if tomorrow {
			fmt.Printf("Tasks due tomorrow (%d):\n\n", len(taskList))
		} else {
			fmt.Printf("Found %d task(s)", len(taskList))
			if skip > 0 {
				fmt.Printf(" (skipped %d)", skip)
			}
			fmt.Printf(":\n\n")
		}
		
		for _, task := range taskList {	
			fmt.Printf("[%d] [%s]\n", task.ID, task.CreatedAt.Format("Jan 2, 2006 15:04"))
			fmt.Printf("%s\n", task.Title)
			if task.Description != "" {
				fmt.Printf("   Description: %s\n", task.Description)
			}
			if task.Deadline != nil {
				deadlineStr := task.Deadline.Format("Jan 2, 2006")
				if isOverdue(*task.Deadline) {
					fmt.Printf("   Deadline: %s (OVERDUE)\n", deadlineStr)
				} else {
					fmt.Printf("   Deadline: %s\n", deadlineStr)
				}
			}
			fmt.Printf("   Status: %s\n\n", task.Status)
		}

		// Show pagination hint if there might be more tasks (only for regular list)
		if !today && !tomorrow && len(taskList) == 10 {
			fmt.Printf("To see more tasks, use: sumb task list --skip %d\n", skip+10)
		}

		return nil
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Update task status",
	Long:  `Update a task's status. Valid statuses: TODO, IN_PROGRESS, COMPLETED.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid task ID: %s", args[0])
		}

		status := db.TaskStatus(strings.ToUpper(args[1]))
		if !isValidStatus(status) {
			return fmt.Errorf("invalid status: %s. Valid statuses: TODO, IN_PROGRESS, COMPLETED", args[1])
		}

		tm, err := db.NewTaskManager()
		if err != nil {
			return fmt.Errorf("failed to initialize task manager: %w", err)
		}
		defer tm.Close()

		if err := tm.UpdateTaskStatus(id, status); err != nil {
			return fmt.Errorf("failed to update task status: %w", err)
		}

		fmt.Println(styles.Separator)
		fmt.Printf("Task status updated!\n")
		fmt.Printf(" Id: %d\n", id)
		fmt.Printf(" New Status: %s\n", status)
		fmt.Println(styles.Separator)
		return nil
	},
}

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search tasks",
	Long:  `Search tasks by matching partial text in title or description.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]
		skip, _ := cmd.Flags().GetInt("skip")
		
		if query == "" {
			return fmt.Errorf("search query is required")
		}
		
		tm, err := db.NewTaskManager()
		if err != nil {
			return fmt.Errorf("failed to initialize task manager: %w", err)
		}
		defer tm.Close()

		taskList, err := tm.SearchTasks(query, 10, skip)
		if err != nil {
			return fmt.Errorf("failed to search tasks: %w", err)
		}

		if len(taskList) == 0 {
			if skip > 0 {
				fmt.Printf("No more search results found for '%s'. You've reached the end of the results.\n", query)
			} else {
				fmt.Printf("No tasks found matching '%s'\n", query)
			}
			return nil
		}

		fmt.Println()
		fmt.Printf("Found %d task(s) matching '%s'", len(taskList), query)
		if skip > 0 {
			fmt.Printf(" (skipped %d)", skip)
		}
		fmt.Printf(":\n\n")
		
		for _, task := range taskList {	
			fmt.Printf("[%d] [%s]\n", task.ID, task.CreatedAt.Format("Jan 2, 2006 15:04"))
			fmt.Printf("%s\n", task.Title)
			if task.Description != "" {
				fmt.Printf("   Description: %s\n", task.Description)
			}
			if task.Deadline != nil {
				deadlineStr := task.Deadline.Format("Jan 2, 2006")
				if isOverdue(*task.Deadline) {
					fmt.Printf("   Deadline: %s (OVERDUE)\n", deadlineStr)
				} else {
					fmt.Printf("   Deadline: %s\n", deadlineStr)
				}
			}
			fmt.Printf("   Status: %s\n\n", task.Status)
		}

		if len(taskList) == 10 {
			fmt.Printf("To see more search results, use: sumb task search \"%s\" --skip %d\n", query, skip+10)
		}

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

		fmt.Println(styles.Separator)
		fmt.Printf(styles.TaskDeleted)
		fmt.Printf(" Id: %d\n", id)
		fmt.Println(styles.Separator)
		return nil
	},
}

var deleteMultipleCmd = &cobra.Command{
	Use:   "delete-many",
	Short: "Delete multiple tasks",
	Long:  `Delete multiple tasks by providing their IDs separated by spaces.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("at least one task ID is required")
		}

		tm, err := db.NewTaskManager()
		if err != nil {
			return fmt.Errorf("failed to initialize task manager: %w", err)
		}
		defer tm.Close()

		var deletedIDs []int
		var failedIDs []string

		for _, arg := range args {
			id, err := strconv.Atoi(arg)
			if err != nil {
				failedIDs = append(failedIDs, arg)
				continue
			}

			if err := tm.DeleteTask(id); err != nil {
				failedIDs = append(failedIDs, arg)
			} else {
				deletedIDs = append(deletedIDs, id)
			}
		}

		if len(deletedIDs) > 0 {
			fmt.Println(styles.Separator)
			fmt.Println(styles.TaskDeletedMany)
			fmt.Printf(" Deleted: %v\n", deletedIDs)
			fmt.Println(styles.Separator)
		}	

		if len(failedIDs) > 0 {
			fmt.Printf("Some tasks could not be deleted: %v\n", failedIDs)
		}

		return nil
	},
}

// Helper functions
func createTaskInteractive() error {
	fmt.Println("Interactive Task Creation")
	fmt.Println("Enter task details (press Enter to skip optional fields):")
	fmt.Println(styles.Separator)
	
	scanner := bufio.NewScanner(os.Stdin)
	
	// Get title
	fmt.Print("Title: ")
	if !scanner.Scan() {
		return fmt.Errorf("failed to read title")
	}
	title := strings.TrimSpace(scanner.Text())
	if title == "" {
		return fmt.Errorf("title is required")
	}
	
	// Get description
	fmt.Print("Description (optional): ")
	if !scanner.Scan() {
		return fmt.Errorf("failed to read description")
	}
	description := strings.TrimSpace(scanner.Text())
	
	// Get deadline
	fmt.Print("Deadline (YYYY-MM-DD, 'today', 'tomorrow', or press Enter to skip): ")
	if !scanner.Scan() {
		return fmt.Errorf("failed to read deadline")
	}
	deadlineStr := strings.TrimSpace(scanner.Text())
	
	deadline, err := parseDeadline(deadlineStr)
	if err != nil {
		return fmt.Errorf("invalid deadline: %w", err)
	}
	
	tm, err := db.NewTaskManager()
	if err != nil {
		return fmt.Errorf("failed to initialize task manager: %w", err)
	}
	defer tm.Close()

	if err := tm.CreateTask(title, description, deadline); err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	fmt.Println(styles.Separator)
	fmt.Printf("Task created!\n")
	fmt.Printf("Title: %s\n", title)
	if description != "" {
		fmt.Printf("Description: %s\n", description)
	}
	if deadline != nil {
		fmt.Printf("Deadline: %s\n", deadline.Format("Jan 2, 2006"))
	}
	fmt.Println(styles.Separator)

	return nil
}

func parseDeadline(deadlineStr string) (*time.Time, error) {
	if deadlineStr == "" {
		return nil, nil
	}
	
	deadlineStr = strings.ToLower(strings.TrimSpace(deadlineStr))
	
	switch deadlineStr {
	case "today":
		now := time.Now()
		return &now, nil
	case "tomorrow":
		tomorrow := time.Now().AddDate(0, 0, 1)
		return &tomorrow, nil
	default:
		deadline, err := time.Parse("2006-01-02", deadlineStr)
		if err != nil {
			return nil, fmt.Errorf("invalid date format. Use YYYY-MM-DD, 'today', or 'tomorrow'")
		}
		return &deadline, nil
	}
}



func isValidStatus(status db.TaskStatus) bool {
	return status == db.StatusTODO || status == db.StatusInProgress || status == db.StatusComplete
}

func isOverdue(deadline time.Time) bool {
	// Only consider the date part, not the time
	deadlineDate := time.Date(deadline.Year(), deadline.Month(), deadline.Day(), 0, 0, 0, 0, deadline.Location())
	today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location())
	return today.After(deadlineDate)
}

func init() {
	createCmd.Flags().StringP("title", "t", "", "Task title (required)")
	createCmd.MarkFlagRequired("body")
	createCmd.Flags().StringP("description", "d", "", "Task description")
	createCmd.Flags().StringP("deadline", "l", "", "Task deadline (YYYY-MM-DD, 'today', 'tomorrow')")
	
	listCmd.Flags().IntP("skip", "s", 0, "Number of tasks to skip (for pagination)")
	listCmd.Flags().Bool("today", false, "Show only tasks due today")
	listCmd.Flags().Bool("tomorrow", false, "Show only tasks due tomorrow")
	
	searchCmd.Flags().IntP("skip", "s", 0, "Number of results to skip (for pagination)")
} 