package main

import (
	"strings"

	"charm.land/bubbles/v2/viewport"
	"charm.land/lipgloss/v2"
)

var (
	descriptionStyle = lipgloss.NewStyle().Padding(1, 2)
	titleStyle       = lipgloss.NewStyle().Foreground(lipgloss.BrightRed).Bold(true).Padding(0, 2)
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

func (m model) setViewportContent(item item) string {

	title := titleStyle.Render("PRICE: " + item.price)
	description := descriptionStyle.Width(m.viewport.Width()).Render(item.description)
	view := lipgloss.JoinVertical(lipgloss.Left, title, description)

	return view

}

func setViewportArrows(vp viewport.Model) (string, string) {
	arrowStyle := lipgloss.NewStyle().Width(vp.Width()).Align(lipgloss.Left).Foreground(lipgloss.Color("241"))

	topIndicator := ""
	bottomIndicator := ""

	if vp.TotalLineCount() > vp.Height() {
		if !vp.AtTop() {
			topIndicator = "↑ More Above"
		}
		if !vp.AtBottom() {
			bottomIndicator = "↓ More Below"
		}
	}

	topArrow := arrowStyle.Render(topIndicator)
	bottomArrow := arrowStyle.Render(bottomIndicator)

	return topArrow, bottomArrow
}

func (m *model) syncViewport() {
	if selected := m.list.SelectedItem(); selected != nil {
		if item, ok := selected.(item); ok {
			if m.lastSelected == "" || item.title != m.lastSelected {
				m.lastSelected = item.title
				m.viewport.SetContent(m.setViewportContent(item))
				m.viewport.GotoTop()
			}
		}
	}
}
