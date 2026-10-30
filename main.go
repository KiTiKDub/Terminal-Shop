package main

import (
	"fmt"
	"os"
	"strings"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	_ "image/jpeg"
	_ "image/png"
)

type focusPane int

const (
	listPane focusPane = iota
	viewportPane
)

var (
	activeBorderStyle   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.BrightRed)
	inactiveBorderStyle = lipgloss.NewStyle().Border(lipgloss.HiddenBorder())
)

type model struct {
	list          list.Model
	viewport      viewport.Model
	viewportReady bool
	lastSelected  string
	choice        string
	styles        listStyle
	quitting      bool
	focus         focusPane
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:

		listWidth := msg.Width / 3
		listHeight := msg.Height
		viewportWidth := 2 * msg.Width / 3
		viewportHeight := msg.Height

		m.list.SetWidth(listWidth)
		m.list.SetHeight(listHeight - 4)

		if !m.viewportReady {
			m.viewport = createViewport(viewportWidth-2, viewportHeight-4)
			m.lastSelected = ""
			m.viewportReady = true
		} else {
			m.viewport.SetWidth(viewportWidth - 2)
			m.viewport.SetHeight(viewportHeight - 4)
		}

		m.syncViewport()

		return m, nil

	case tea.KeyPressMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i.title)
			}
			return m, nil
		case "backspace":
			m.choice = ""
			return m, nil
		case "tab":
			if m.focus == listPane {
				m.focus = viewportPane

			} else {
				m.focus = listPane
			}
			return m, nil
		}

	}

	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch m.focus {
	case listPane:
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	case viewportPane:
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.syncViewport()

	return m, tea.Batch(cmds...)
}

func (m model) View() tea.View {

	if m.quitting {
		return tea.NewView(m.styles.quitText.Render("Thanks for shopping!"))
	}
	if m.choice != "" {
		v := tea.NewView(openPluginView(m))
		v.AltScreen = true
		return v
	}

	var desc string

	if !m.viewportReady {
		desc = "Loading..."
	} else {
		topArrow, bottomArrow := setViewportArrows(m.viewport)
		// desc = m.viewport.View()
		desc = lipgloss.JoinVertical(lipgloss.Left, topArrow, m.viewport.View(), bottomArrow)
	}

	leftStyle, rightStyle := setActiveBorder(m.focus)

	var s strings.Builder

	s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, leftStyle.Render(m.list.View()), rightStyle.Render(desc)))

	v := tea.NewView("\n" + s.String())
	v.AltScreen = true

	return v
}

func setActiveBorder(focus focusPane) (lipgloss.Style, lipgloss.Style) {
	leftStyle := inactiveBorderStyle
	rightStyle := inactiveBorderStyle

	switch focus {
	case listPane:
		leftStyle = activeBorderStyle
	case viewportPane:
		rightStyle = activeBorderStyle

	}

	return leftStyle, rightStyle

}

func main() {
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
