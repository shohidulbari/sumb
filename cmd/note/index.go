package note

import (
	"github.com/spf13/cobra"
)

var NoteCmd = &cobra.Command{
	Use:   "note",
	Short: "Manage notes",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	NoteCmd.AddCommand(createCmd)
	NoteCmd.AddCommand(searchCmd)
	NoteCmd.AddCommand(listCmd)
	NoteCmd.AddCommand(showCmd)
}
