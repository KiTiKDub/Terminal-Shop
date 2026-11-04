package main

import (
	"fmt"
	"os"

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
	pluginPane
)

var (
	activeBorderStyle   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.BrightRed)
	inactiveBorderStyle = lipgloss.NewStyle().Border(lipgloss.HiddenBorder())
)

type model struct {
	list           list.Model
	downloadList   HorizontalList
	viewport       viewport.Model
	pluginViewport viewport.Model
	lastSelected   string
	choice         string
	downloadMsg    string
	styles         listStyle
	quitting       bool
	focus          focusPane
	header         string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:

		headerHeight := m.getHeaderLineCount()

		listWidth := msg.Width / 3
		listHeight := msg.Height - headerHeight
		viewportWidth := 2 * msg.Width / 3
		viewportHeight := msg.Height - headerHeight

		m.list.SetWidth(listWidth)
		m.list.SetHeight(listHeight)

		m.viewport.SetWidth(viewportWidth - 2)
		m.viewport.SetHeight(viewportHeight)
		m.pluginViewport.SetWidth(msg.Width)
		m.pluginViewport.SetHeight(viewportHeight - 4)

		return m, nil

	case tea.KeyPressMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":

			if m.focus == pluginPane && m.downloadList.SelectedItem().title == "Download" {
				m.downloadMsg = download(m.choice)
				return m, nil
			}

			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i.title)
				m.focus = pluginPane
				m.setPluginViewportContent()
				return m, nil
			}

		case "left":
			if m.choice != "" {
				m.downloadList.MoveLeft()
				m.setPluginViewportContent()
			}

		case "right":
			if m.choice != "" {
				m.downloadList.MoveRight()
				m.setPluginViewportContent()
			}

		case "backspace", "esc":
			m.choice = ""
			m.downloadMsg = ""
			m.focus = listPane
			return m, nil

		case "tab":
			if m.focus != pluginPane {
				if m.focus == listPane {
					m.focus = viewportPane

				} else {
					m.focus = listPane
				}
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
	case pluginPane:
		m.pluginViewport, cmd = m.pluginViewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.syncDescViewport()

	return m, tea.Batch(cmds...)
}

func (m model) View() tea.View {

	if m.quitting {
		return tea.NewView(m.styles.quitText.Render("Thanks for shopping!"))
	}
	if m.downloadMsg != "" {
		v := tea.NewView(m.downloadMsg)
		v.AltScreen = true
		return v
	}
	if m.choice != "" {
		v := tea.NewView(buildPluginView(m))
		v.AltScreen = true
		return v
	}

	v := tea.NewView(buildMainView(m))
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

func buildMainView(m model) string {

	topArrow, bottomArrow := setViewportArrows(m.viewport)
	desc := lipgloss.JoinVertical(lipgloss.Left, topArrow, m.viewport.View(), bottomArrow)

	leftStyle, rightStyle := setActiveBorder(m.focus)

	body := lipgloss.JoinHorizontal(lipgloss.Center, leftStyle.Width(m.list.Width()).Render(m.list.View()),
		rightStyle.Render(desc))
	return lipgloss.JoinVertical(lipgloss.Left, m.header, body)
}

func buildPluginView(m model) string {
	topArrow, bottomArrow := setViewportArrows(m.pluginViewport)
	viewportString := lipgloss.JoinVertical(lipgloss.Left, topArrow, m.pluginViewport.View(), bottomArrow)
	return lipgloss.JoinVertical(lipgloss.Left, m.header, viewportString)
}

func initialModel() model {

	m := model{}
	m.setProductList()
	m.updateProductListStyles()
	m.setDownloadList()
	m.setHeaderString("KiTiK Store")
	m.setViewport()
	m.setPluginViewport()

	return m
}

func main() {
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
