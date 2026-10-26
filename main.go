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
		m.list.SetWidth(msg.Width / 3)

		if !m.viewportReady {
			m.viewport = createViewport(2*msg.Width/3, msg.Height)
			m.lastSelected = ""
			m.viewportReady = true
		} else {
			m.viewport.SetWidth(2 * msg.Width / 3)
			m.viewport.SetHeight(msg.Height)
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
		desc = m.viewport.View()
	}

	var s strings.Builder
	s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, m.list.View(), desc))

	v := tea.NewView("\n" + s.String())
	v.AltScreen = true

	return v
}

func main() {
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
