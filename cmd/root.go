package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sumb/cmd/notes"
	"github.com/sumb/cmd/pomodoro"
	"github.com/sumb/cmd/tasks"
)

var rootCmd = &cobra.Command{
	Use:   "sumb",
	Short: "A terminal-based task, note, and pomodoro management application",
	Long: `Sumb is a command-line task, note, and pomodoro management tool that helps you organize and track your tasks, notes, and productivity sessions.

Features:
- Task management with completion tracking
- Note management with simple body content
- Pomodoro timer with session tracking
- SQLite database storage for persistence

Usage:
  sumb task      # Task management
  sumb note      # Note management
  sumb pomodoro  # Pomodoro timer management
  sumb pd        # Alias for pomodoro`,
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
	// Add pomodoro subcommand
	rootCmd.AddCommand(pomodoro.PomodoroCmd)
	// Add pomodoro alias
	rootCmd.AddCommand(pomodoro.PdCmd)
} 