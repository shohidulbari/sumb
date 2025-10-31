package cmd

import (
	"github.com/shohidulbari/sumb/cmd/note"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sumb",
	Short: "Note manager",
	Long: `Sumb is a terminal-based application to manage your notes 

	Features:
	- Create notes using a text area form
	- List latest notes in a tree view
	- Search notes by keywords
	- View note content in a scrollable viewport

Usage:
  sumb create      # Note creation
  sumb list        # List latest notes
  sumb search      # Search notes by keyword
  sumb show        # Show note in a viewport`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// If no subcommand is provided, show help
		return cmd.Help()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(note.CreateCmd)
	rootCmd.AddCommand(note.SearchCmd)
	rootCmd.AddCommand(note.ListCmd)
	rootCmd.AddCommand(note.ShowCmd)
	rootCmd.AddCommand(note.EditCmd)
	rootCmd.AddCommand(note.DeleteCmd)
}
