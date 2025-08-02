package notes

import (
	"fmt"

	"github.com/spf13/cobra"
	db "github.com/sumb/db"
)

var NoteCmd = &cobra.Command{
	Use:   "note",
	Short: "Manage notes",
	Long:  `Manage your notes with various operations like create, list, and delete.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Handle quick create with -c flag
		body, _ := cmd.Flags().GetString("create")
		
		if body != "" {
			nm, err := db.NewNoteManager()
			if err != nil {
				return fmt.Errorf("failed to initialize note manager: %w", err)
			}
			defer nm.Close()

			if err := nm.CreateNote(body); err != nil {
				return fmt.Errorf("failed to create note: %w", err)
			}

			fmt.Println("--------------------------------")

			fmt.Printf("🌟 Note created!\n")
			fmt.Printf("%s\n", body)
			fmt.Println("--------------------------------")

			return nil
		}
		return cmd.Help()
	},
}

func init() {
	NoteCmd.AddCommand(createCmd)
	NoteCmd.AddCommand(listCmd)
	NoteCmd.AddCommand(deleteCmd)
	NoteCmd.AddCommand(deleteMultipleCmd)
	
	NoteCmd.Flags().StringP("create", "c", "", "Quick create a note with body")
} 