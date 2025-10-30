package note

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
)

type page struct {
	viewportTitle string
	content       string
	ready         bool
	viewport      viewport.Model
}

func (p page) Init() tea.Cmd {
	return nil
}

func (p page) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return p, tea.Quit
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(p.headerView())

		if !p.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			p.viewport = viewport.New(100, 20)
			p.viewport.YPosition = headerHeight
			p.viewport.SetContent(p.content)
			p.ready = true
		} else {
			p.viewport.Width = 100
			p.viewport.Height = 20
		}
	}

	// Handle keyboard and mouse events in the viewport
	p.viewport, cmd = p.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return p, tea.Batch(cmds...)
}

func (p page) View() string {
	if !p.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf(
		"%s\n%s\n%s\n(To copy the content, first exit the viewport)\n(To exit, Press 'q' or 'ctrl+c')\n\n",
		p.headerView(),
		p.viewport.View(),
		p.footerView(),
	)
}

func (p page) headerView() string {
	title := titleStyle.Render(p.viewportTitle)
	line := strings.Repeat("─", max(0, p.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (p page) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", p.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, p.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
