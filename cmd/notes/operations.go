package notes

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/sumb/cmd/styles"
	db "github.com/sumb/db"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new note",
	Long:  `Create a new note with a body.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		body, _ := cmd.Flags().GetString("body")

		if body == "" {
			return fmt.Errorf("body is required")
		}

		nm, err := db.NewNoteManager()
		if err != nil {
			return fmt.Errorf("failed to initialize note manager: %w", err)
		}
		defer nm.Close()

		if err := nm.CreateNote(body); err != nil {
			return fmt.Errorf("failed to create note: %w", err)
		}

		fmt.Printf(styles.NoteCreated)
		fmt.Printf("Body: %s\n", body)

		return nil
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all notes",
	Long:  `Display notes in the database with their details. Shows max 10 latest notes by default. Use --skip to paginate through older notes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		skip, _ := cmd.Flags().GetInt("skip")
		
		nm, err := db.NewNoteManager()
		if err != nil {
			return fmt.Errorf("failed to initialize note manager: %w", err)
		}
		defer nm.Close()

		noteList, err := nm.ListNotesWithPagination(10, skip)
		if err != nil {
			return fmt.Errorf("failed to list notes: %w", err)
		}

		if len(noteList) == 0 {
			if skip > 0 {
				fmt.Printf("📝 No more notes found. You've reached the end of your note list.\n")
			} else {
				fmt.Println("📝 No notes found. Create your first note with 'sumb note create -b \"Your note content\"'")
			}
			return nil
		}

		fmt.Println()
		fmt.Printf("📋 Found %d note(s)", len(noteList))
		if skip > 0 {
			fmt.Printf(" (skipped %d)", skip)
		}
		fmt.Printf(":\n\n")
		
		for idx, note := range noteList {
			if idx == 0 {
				fmt.Println(styles.Separator)
			}

			fmt.Printf("📄 [%d] %s\n", note.ID, note.Body)
			fmt.Printf("   Created: %s\n", note.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Println(styles.Separator)
		}

		// Show pagination hint if there might be more notes
		if len(noteList) == 10 {
			fmt.Printf("💡 To see more notes, use: sumb note list --skip %d\n", skip+10)
		}

		return nil
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a note",
	Long:  `Delete a note by providing its ID.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("note ID is required")
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid note ID: %s", args[0])
		}

		nm, err := db.NewNoteManager()
		if err != nil {
			return fmt.Errorf("failed to initialize note manager: %w", err)
		}
		defer nm.Close()

		if err := nm.DeleteNote(id); err != nil {
			return fmt.Errorf("failed to delete note: %w", err)
		}

		fmt.Println(styles.Separator)
		fmt.Printf(styles.NoteDeleted)
		fmt.Printf(" Id: %d\n", id)
		fmt.Println(styles.Separator)
		return nil
	},
}

var deleteMultipleCmd = &cobra.Command{
	Use:   "delete-many",
	Short: "Delete multiple notes",
	Long:  `Delete multiple notes by providing their IDs separated by spaces.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("at least one note ID is required")
		}

		nm, err := db.NewNoteManager()
		if err != nil {
			return fmt.Errorf("failed to initialize note manager: %w", err)
		}
		defer nm.Close()

		var deletedIDs []int
		var failedIDs []string

		for _, arg := range args {
			id, err := strconv.Atoi(arg)
			if err != nil {
				failedIDs = append(failedIDs, arg)
				continue
			}

			if err := nm.DeleteNote(id); err != nil {
				failedIDs = append(failedIDs, arg)
			} else {
				deletedIDs = append(deletedIDs, id)
			}
		}

		if len(deletedIDs) > 0 {
			fmt.Println(styles.Separator)
			fmt.Println(styles.NoteDeletedMany)
			fmt.Printf(" Deleted: %v\n", deletedIDs)
			fmt.Println(styles.Separator)
		}	

		if len(failedIDs) > 0 {
			fmt.Printf("Some notes could not be deleted: %v\n", failedIDs)
		}

		return nil
	},
}

func init() {
	createCmd.Flags().StringP("body", "b", "", "Note body (required)")
	createCmd.MarkFlagRequired("body")
	
	// Add pagination flag for list command
	listCmd.Flags().IntP("skip", "s", 0, "Number of notes to skip (for pagination)")
} 