package notes

import (
	"fmt"
	"strconv"
	"strings"

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
		jsonify, _ := cmd.Flags().GetBool("jsonify")
		
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
		
		for _, note := range noteList {	
			fmt.Printf("[%d] [%s]\n", note.ID, note.CreatedAt.Format("Jan 2, 2006 15:04"))
			if jsonify && isJSON(note.Body) {
				formattedJSON, err := formatJSON(note.Body)
				if err == nil {
					fmt.Printf("%s\n\n", indentJSON(formattedJSON))
				}
			} else {
				fmt.Printf("%s\n\n", note.Body)
			}	
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

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search notes by content",
	Long:  `Search through your notes by matching partial content in the note body.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]
		skip, _ := cmd.Flags().GetInt("skip")
		jsonify, _ := cmd.Flags().GetBool("jsonify")
		
		if query == "" {
			return fmt.Errorf("search query is required")
		}
		
		nm, err := db.NewNoteManager()
		if err != nil {
			return fmt.Errorf("failed to initialize note manager: %w", err)
		}
		defer nm.Close()

		noteList, err := nm.SearchNotes(query, 10, skip)
		if err != nil {
			return fmt.Errorf("failed to search notes: %w", err)
		}

		if len(noteList) == 0 {
			if skip > 0 {
				fmt.Printf("🔍 No more search results found for '%s'. You've reached the end of the results.\n", query)
			} else {
				fmt.Printf("🔍 No notes found matching '%s'\n", query)
			}
			return nil
		}

		fmt.Println()
		fmt.Printf("🔍 Found %d note(s) matching '%s'", len(noteList), query)
		if skip > 0 {
			fmt.Printf(" (skipped %d)", skip)
		}
		fmt.Printf(":\n\n")
		
		for _, note := range noteList {	
			fmt.Printf("[%d] [%s]\n", note.ID, note.CreatedAt.Format("Jan 2, 2006 15:04"))
			if jsonify && isJSON(note.Body) {
				formattedJSON, err := formatJSON(note.Body)
				if err == nil {
					fmt.Printf("%s\n\n", indentJSON(formattedJSON))
				}
			} else {
				fmt.Printf("%s\n\n", note.Body)
			}	
		}

		if len(noteList) == 10 {
			fmt.Printf("💡 To see more search results, use: sumb note search \"%s\" --skip %d\n", query, skip+10)
		}

		return nil
	},
}

func indentJSON(jsonStr string) string {
	lines := strings.Split(jsonStr, "\n")
	var indentedLines []string
	
	for _, line := range lines {
		indentedLines = append(indentedLines, "   "+line)
	}
	
	return strings.Join(indentedLines, "\n")
}

func init() {
	createCmd.Flags().StringP("body", "b", "", "Note body (required)")
	createCmd.MarkFlagRequired("body")	
	listCmd.Flags().IntP("skip", "s", 0, "Number of notes to skip (for pagination)")
	listCmd.Flags().Bool("jsonify", false, "Format JSON output in a pretty format")	
	searchCmd.Flags().IntP("skip", "s", 0, "Number of results to skip (for pagination)")
	searchCmd.Flags().Bool("jsonify", false, "Format JSON output in a pretty format")
} 