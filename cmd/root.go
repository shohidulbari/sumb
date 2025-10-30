package cmd

import (
	"github.com/shohidulbari/sumb/cmd/note"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sumb",
	Short: "A terminal-based productivity application",
	Long: `Sumb is a command-line productivity tool that helps you organize and track your tasks, notes, and productivity sessions.

Features:
- Manage tasks, notes and pomodoro sessions using the command line.

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
	rootCmd.AddCommand(note.NoteCmd)
}
