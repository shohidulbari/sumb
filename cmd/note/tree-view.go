package note

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
	"github.com/muesli/reflow/wordwrap"
	"github.com/shohidulbari/sumb/db"
)

var (
	NOTE_BODY_TAKE_LIMIT = 200
	HIGHLIGHT_COLOR      = lipgloss.Color("63")
)

func formatNoteBody(body string) string {
	if (len(body)) > NOTE_BODY_TAKE_LIMIT {
		remainingCharacters := len(body) - NOTE_BODY_TAKE_LIMIT
		take := body[:NOTE_BODY_TAKE_LIMIT]
		wrapped := wordwrap.String(take, 50)

		message := Warn.Render(
			fmt.Sprintf("... (truncated, %d more characters)", remainingCharacters),
		)
		return fmt.Sprintf(
			"%s... %s",
			wrapped,
			message,
		)
	}
	return wordwrap.String(body, 50)
}

func RenderTreeView(root string, notes []db.NoteResponse) *tree.Tree {
	enumeratorStyle := lipgloss.NewStyle().Foreground(HIGHLIGHT_COLOR).MarginRight(1)
	rootStyle := lipgloss.NewStyle().Foreground(HIGHLIGHT_COLOR).Bold(true)
	t := tree.Root(root)

	for _, note := range notes {
		noteTree := tree.New().Root(fmt.Sprintf("[%s]", note.ID)).Child(formatNoteBody(note.Body))
		t = t.Child(noteTree)
	}
	t = t.Enumerator(tree.RoundedEnumerator).
		EnumeratorStyle(enumeratorStyle).
		RootStyle(rootStyle)
	return t
}
