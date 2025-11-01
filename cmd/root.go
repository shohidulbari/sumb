package cmd

import (
	"github.com/shohidulbari/sumb/cmd/note"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sumb",
	Short: "Note manager",
	Long: `Sumb is a terminal-based application to manage your notes with full-text search capabilities.

	Features:
	- Create notes using a text area form
	- List latest notes in a tree view
	- Search notes by keywords
	- View note content in a scrollable viewport
	- Edit notes using a text area form
	- Delete notes by ID

Usage:
  sumb create      			# Note creation
  sumb list        			# List latest notes
  sumb search      			# Search notes by keyword
  sumb show <Note ID>   # Show note in a viewport,
	sumb edit <Note ID>   # Edit note by ID
	sumb delete <Note ID> # Delete note by ID`,
	RunE: func(cmd *cobra.Command, args []string) error {
		versionFlag, err := cmd.Flags().GetBool("version")
		if err != nil {
			return err
		}
		if versionFlag {
			cmd.Println("sumb version 1.0.0")
			return nil
		}
		// If no subcommand is provided, show help
		return cmd.Help()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().BoolP("version", "v", false, "Show the current version")
	rootCmd.AddCommand(note.CreateCmd)
	rootCmd.AddCommand(note.SearchCmd)
	rootCmd.AddCommand(note.ListCmd)
	rootCmd.AddCommand(note.ShowCmd)
	rootCmd.AddCommand(note.EditCmd)
	rootCmd.AddCommand(note.DeleteCmd)
}
