package note

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type errMsg error

type model struct {
	textarea textarea.Model
	err      error
	canceled bool
}

func RenderForm() *model {
	ti := textarea.New()
	ti.Focus()
	ti.CharLimit = 1000
	ti.Prompt = "|"
	ti.FocusedStyle.Prompt = ti.FocusedStyle.Prompt.Foreground(lipgloss.Color("#ffffffff"))
	ti.FocusedStyle.CursorLine = ti.FocusedStyle.CursorLine.Background(lipgloss.Color("#242424ff"))
	ti.CharLimit = 4096
	ti.SetHeight(25)
	ti.SetWidth(100)
	return &model{
		textarea: ti,
		err:      nil,
	}
}

func (m *model) Init() tea.Cmd {
	return textarea.Blink
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if m.textarea.Focused() {
				m.textarea.Blur()
			}
		case tea.KeyCtrlS:
			return m, tea.Quit
		case tea.KeyCtrlC:
			m.canceled = true
			return m, tea.Quit
		default:
			if !m.textarea.Focused() {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *model) View() string {
	return fmt.Sprintf(
		"\n\n%s\n\n%s %s",
		m.textarea.View(),
		"(ctrl+s to save and exit)",
		"(ctrl+c to cancel)",
	) + "\n\n"
}
