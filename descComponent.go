package main

import (
	"strings"

	"charm.land/bubbles/v2/viewport"
	"charm.land/lipgloss/v2"
)

var (
	descriptionStyle = lipgloss.NewStyle().Padding(1, 2).MarginTop(1)
	titleStyle       = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).Padding(0, 1)
	infoStyle        = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).Padding(0, 1)
)

func headerView() string {
	title := titleStyle.Render("Description")
	return title
}

func footerView(w int) string {
	line := strings.Repeat("-", max(0, w))
	return line
}

func createViewport(w int, h int) viewport.Model {
	headerHeight := lipgloss.Height(headerView())
	footerHeight := lipgloss.Height(footerView(w))

	viewport := viewport.New(viewport.WithWidth(w), viewport.WithHeight(h-(headerHeight+footerHeight)))
	viewport.YPosition = headerHeight

	viewport.HighlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Background(lipgloss.Color("34"))
	viewport.SelectedHighlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Background(lipgloss.Color("47"))
	viewport.HighlightNext()

	return viewport
}

func (m model) descView(item item) string {

	// contentWidth := m.list.Width() * 2

	// return descriptionStyle.Render(headerView(), item.description, footerView(contentWidth))
	return descriptionStyle.Width(m.viewport.Width()).Render(item.description)

}

func (m *model) syncViewport() {
	if selected := m.list.SelectedItem(); selected != nil {
		if item, ok := selected.(item); ok {
			if m.lastSelected == "" || item.title != m.lastSelected {
				m.lastSelected = item.title
				m.viewport.SetContent(m.descView(item))
				m.viewport.GotoTop()
			}
		}
	}
}
