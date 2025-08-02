package tasks

import (
	"fmt"

	"github.com/spf13/cobra"
	db "github.com/sumb/db"
)

var TaskCmd = &cobra.Command{
	Use:   "task",
	Short: "Manage tasks",
	Long:  `Manage your tasks with various operations like create, list, complete, and delete.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Handle quick create with -c flag
		title, _ := cmd.Flags().GetString("create")
		description, _ := cmd.Flags().GetString("description")
		
		if title != "" {
			tm, err := db.NewTaskManager()
			if err != nil {
				return fmt.Errorf("failed to initialize task manager: %w", err)
			}
			defer tm.Close()

			if err := tm.CreateTask(title, description); err != nil {
				return fmt.Errorf("failed to create task: %w", err)
			}

			fmt.Printf("✅ Task created successfully!\n")
			fmt.Printf("Title: %s\n", title)
			if description != "" {
				fmt.Printf("Description: %s\n", description)
			}

			return nil
		}
		return cmd.Help()
	},
}

func init() {
	TaskCmd.AddCommand(createCmd)
	TaskCmd.AddCommand(listCmd)
	TaskCmd.AddCommand(completeCmd)
	TaskCmd.AddCommand(deleteCmd)
	
	TaskCmd.Flags().StringP("create", "c", "", "Quick create a task with title")
	TaskCmd.Flags().StringP("description", "d", "", "Task description (for quick create)")
} 