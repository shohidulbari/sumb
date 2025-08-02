package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sumb/cmd/tasks"
)

var rootCmd = &cobra.Command{
	Use:   "sumb",
	Short: "A terminal-based task management application",
	Long: `Sumb is a command-line task management tool that helps you organize and track your tasks.
	
Features:
- Create tasks with titles and descriptions
- List all tasks
- Mark tasks as complete
- Delete tasks
- SQLite database storage for persistence

Usage:
  sumb task -c "task-title" -d "task-description"  # Quick create
  sumb task create -t "title" -d "description"     # Create task
  sumb task list                                    # List all tasks
  sumb task complete <id>                           # Mark task complete
  sumb task delete <id>                             # Delete task
  sumb task delete-many <id1> <id2> <id3>           # Delete multiple tasks`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// If no subcommand is provided, show help
		return cmd.Help()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add task subcommand
	rootCmd.AddCommand(tasks.TaskCmd)
} 