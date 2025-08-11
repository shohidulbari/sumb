package tasks

import (
	"fmt"

	db "github.com/shohidulbari/sumb/db"
	"github.com/spf13/cobra"
)

var TaskCmd = &cobra.Command{
	Use:   "task",
	Short: "Manage tasks",
	Long:  `Easily create, track, update and remove tasks`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Handle quick create with -c flag
		title, _ := cmd.Flags().GetString("create")
		description, _ := cmd.Flags().GetString("description")
		deadlineStr, _ := cmd.Flags().GetString("deadline")
		interactive, _ := cmd.Flags().GetBool("interactive")
		
		if interactive {
			return createTaskInteractive()
		}
		
		if title != "" {
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

			fmt.Println("--------------------------------")

			fmt.Printf("Task created!\n")
			fmt.Printf("Title: %s\n", title)
			if description != "" {
				fmt.Printf("Description: %s\n", description)
			}
			if deadline != nil {
				fmt.Printf("Deadline: %s\n", deadline.Format("Jan 2, 2006"))
			}
			fmt.Println("--------------------------------")

			return nil
		}
		return cmd.Help()
	},
}



func init() {
	TaskCmd.AddCommand(createCmd)
	TaskCmd.AddCommand(listCmd)
	TaskCmd.AddCommand(statusCmd)
	TaskCmd.AddCommand(searchCmd)
	TaskCmd.AddCommand(deleteCmd)
	TaskCmd.AddCommand(deleteMultipleCmd)
	
	TaskCmd.Flags().StringP("create", "c", "", "Quick create a task with title")
	TaskCmd.Flags().StringP("description", "d", "", "Task description (for quick create)")
	TaskCmd.Flags().StringP("deadline", "l", "", "Task deadline (YYYY-MM-DD, 'today', 'tomorrow')")
	TaskCmd.Flags().BoolP("interactive", "i", false, "Create task interactively")
} 