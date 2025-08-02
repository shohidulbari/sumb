package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sumb/cmd/notes"
	"github.com/sumb/cmd/tasks"
)

var rootCmd = &cobra.Command{
	Use:   "sumb",
	Short: "A terminal-based task and note management application",
	Long: `Sumb is a command-line task and note management tool that helps you organize and track your tasks and notes.

Features:
- Task management with completion tracking
- Note management with simple body content
- SQLite database storage for persistence

Usage:
  sumb task     # Task management
  sumb note     # Note management`,
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
	// Add note subcommand
	rootCmd.AddCommand(notes.NoteCmd)
} 