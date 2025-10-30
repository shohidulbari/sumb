package note

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
	"github.com/shohidulbari/sumb/db"
)

func RenderTreeView(root string, notes []db.NoteResponse) *tree.Tree {
	enumeratorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("63")).MarginRight(1)
	rootStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	itemStyle := lipgloss.NewStyle().Width(50)
	t := tree.Root(root)

	for _, note := range notes {
		noteTree := tree.New().Root(fmt.Sprintf("[%s]", note.ID)).Child(note.Body)
		t = t.Child(noteTree)
	}
	t = t.Enumerator(tree.RoundedEnumerator).
		EnumeratorStyle(enumeratorStyle).
		RootStyle(rootStyle).
		ItemStyle(itemStyle)
	return t
}
