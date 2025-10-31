package note

import (
	"fmt"
	"os"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/wordwrap"
	"github.com/shohidulbari/sumb/db"
	"github.com/spf13/cobra"
)

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "New note",
	Long:  `Create a new note`,
	RunE: func(cmd *cobra.Command, args []string) error {
		form := RenderForm(nil)
		p := tea.NewProgram(form)
		if _, err := p.Run(); err != nil {
			fmt.Println(GetAlertWithUsageInfo("Failed to render note form", err.Error()))
			return nil
		}

		if form.canceled {
			fmt.Println(Warn.Render("Note creation canceled."))
			return nil
		}

		noteBody := form.textarea.Value()
		note, err := db.Create(noteBody)
		if err != nil {
			return err
		}
		fmt.Println(
			Success.Render(fmt.Sprintf("🎉 Note created successfully with ID %s.\n", note.ID)),
		)

		return nil
	},
}

var EditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit note by ID",
	Long:  `Edit note details by ID`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			fmt.Println(GetAlertWithUsageInfo("Note ID is required", "sumb edit <Note ID>"))
			return nil
		}
		noteID := args[0]
		note, err := db.GetNoteById(noteID)
		if err != nil {
			fmt.Println(GetAlertWithUsageInfo("Invalid ID", "Use valid ID: sumb edit <Note ID>"))
			return nil
		}

		form := RenderForm(&note.Body)
		form.preValue = note.Body
		p := tea.NewProgram(form)
		if _, err := p.Run(); err != nil {
			fmt.Println(GetAlertWithUsageInfo("Failed to render note edit form", err.Error()))
			return nil
		}

		if form.canceled {
			fmt.Println(Warn.Render("Note update canceled."))
			return nil
		}

		noteBody := form.textarea.Value()
		err = db.Update(note.ID, noteBody)
		if err != nil {
			return err
		}
		fmt.Println(
			Success.Render(fmt.Sprintf("🎉 Note with ID %s updated successfully.\n", note.ID)),
		)
		return nil
	},
}

var SearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search notes",
	Long:  `Search notes by text`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			fmt.Println(GetAlertWithUsageInfo("Search string is required", "sumb search <keyword>"))
			return nil
		}
		keyword := args[0]
		if keyword == "" {
			fmt.Println(GetAlertWithUsageInfo("Search string is required", "sumb search <keyword>"))
		}

		notes, err := db.Search(keyword)
		if err != nil {
			return err
		}

		if len(notes) == 0 {
			fmt.Println(Info.Render("No notes found. Try different keywords."))
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
				fmt.Println(
					GetAlertWithUsageInfo("Size must be a valid integer", "sumb list <size>"),
				)
			}
		}
		notes, err := db.ListLatestNotes(size)
		if err != nil {
			return err
		}

		if len(notes) == 0 {
			fmt.Println(Info.Render("No notes found. Create a new note using: sumb create"))
			return nil
		}

		tree := RenderTreeView(fmt.Sprintf("Latest %d notes", len(notes)), notes)
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
			fmt.Println(GetAlertWithUsageInfo("Note ID is required", "sumb show <Note ID>"))
		}
		noteID := args[0]
		note, err := db.GetNoteById(noteID)
		if err != nil {
			fmt.Println(GetAlertWithUsageInfo("Invalid ID", "Use valid ID: sumb show <Note ID>"))
			return nil
		}
		body := wordwrap.String(note.Body, 80)
		p := tea.NewProgram(
			page{content: body, viewportTitle: fmt.Sprintf("Note ID: %s", note.ID)},
			tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
		)

		if _, err := p.Run(); err != nil {
			fmt.Println(GetAlertWithUsageInfo("Could not run viewport", err.Error()))
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
			fmt.Println(GetAlertWithUsageInfo("Note ID is required", "sumb delete <Note ID>"))
			return nil
		}
		noteID := args[0]
		_, err := db.GetNoteById(noteID)
		if err != nil {
			fmt.Println(GetAlertWithUsageInfo("Invalid ID", "Use valid ID: sumb delete <Note ID>"))
			return nil
		}

		err = db.Delete(noteID)
		if err != nil {
			return err
		}
		fmt.Println(
			Success.Render(fmt.Sprintf("🗑 Note with ID %s deleted successfully.\n", noteID)),
		)
		return nil
	},
}
