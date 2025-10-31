package note

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/wordwrap"
	"github.com/shohidulbari/sumb/db"
	"github.com/spf13/cobra"
)

var ErrKeywordRequired = errors.New("keyword is required for searching notes")

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "New note",
	Long:  `Create a new note`,
	RunE: func(cmd *cobra.Command, args []string) error {
		form := RenderForm(nil)
		p := tea.NewProgram(form)
		if _, err := p.Run(); err != nil {
			return fmt.Errorf("failed to run note form: %w", err)
		}

		if form.canceled {
			fmt.Println("Note creation canceled.")
			return nil
		}

		noteBody := form.textarea.Value()
		note, err := db.Create(noteBody)
		if err != nil {
			return err
		}
		fmt.Printf("Note created with ID: %s, %s\n", note.ID, note.Body)

		return nil
	},
}

var EditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit note by ID",
	Long:  `Edit note details by ID`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("note ID is required")
		}
		noteID := args[0]
		note, err := db.GetNoteById(noteID)
		if err != nil {
			return err
		}
		if note == nil {
			fmt.Printf("Note with ID %s not found.\n", noteID)
			return nil
		}

		form := RenderForm(&note.Body)
		form.preValue = note.Body
		p := tea.NewProgram(form)
		if _, err := p.Run(); err != nil {
			return fmt.Errorf("failed to run note form: %w", err)
		}

		if form.canceled {
			fmt.Println("Note editing canceled.")
			return nil
		}

		noteBody := form.textarea.Value()
		err = db.Update(note.ID, noteBody)
		if err != nil {
			return err
		}
		log.Printf("Updated!")
		return nil
	},
}

var SearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search notes",
	Long:  `Search notes by text`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("Search string is required. Use sumb search <keyword>")
		}
		keyword := args[0]
		if keyword == "" {
			return fmt.Errorf("%w", ErrKeywordRequired)
		}

		notes, err := db.Search(keyword)
		if err != nil {
			return err
		}

		if len(notes) == 0 {
			fmt.Println("No notes found.")
			return nil
		}

		tree := RenderTreeView(fmt.Sprintf("Search results for '%s'", keyword), notes)
		fmt.Println(tree)
		return nil
	},
}

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List latest n notes",
	Long:  `List latest n number of notes`,
	RunE: func(cmd *cobra.Command, args []string) error {
		size := 10
		var err error
		if len(args) > 0 {
			size, err = strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("Invalid size value: %s", err.Error())
			}
		}
		notes, err := db.ListLatestNotes(size)
		if err != nil {
			return err
		}

		if len(notes) == 0 {
			fmt.Println("No notes found.")
			return nil
		}

		tree := RenderTreeView(fmt.Sprintf("Latest %d notes", size), notes)
		fmt.Println(tree)
		return nil
	},
}

var ShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show note by ID",
	Long:  `Show note details by ID`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("note ID is required")
		}
		noteID := args[0]
		note, err := db.GetNoteById(noteID)
		if err != nil {
			return err
		}
		if note == nil {
			fmt.Printf("Note with ID %s not found.\n", noteID)
			return nil
		}
		body := wordwrap.String(note.Body, 80)
		p := tea.NewProgram(
			page{content: body, viewportTitle: fmt.Sprintf("Note ID: %s", note.ID)},
			tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
		)

		if _, err := p.Run(); err != nil {
			fmt.Println("could not run program:", err)
			os.Exit(1)
		}
		return nil
	},
}

var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete note by ID",
	Long:  `Delete note details by ID`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("note ID is required")
		}
		noteID := args[0]
		note, err := db.GetNoteById(noteID)
		if err != nil {
			return err
		}
		if note == nil {
			fmt.Printf("Note with ID %s not found.\n", noteID)
			return nil
		}

		err = db.Delete(noteID)
		if err != nil {
			return err
		}
		fmt.Printf("Note with ID %s deleted successfully.\n", noteID)
		return nil
	},
}
