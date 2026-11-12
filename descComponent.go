package main

import (
	"strings"

	"charm.land/bubbles/v2/viewport"
	"charm.land/lipgloss/v2"
)

var (
	descHeroTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("39")).
				Bold(true).
				Padding(0, 2)

	descMetaStyle = lipgloss.NewStyle().
			Padding(0, 2)

	descPlatformBadge = lipgloss.NewStyle().
				Background(lipgloss.Color("238")).
				Foreground(lipgloss.Color("250")).
				Padding(0, 1)

	descDividerStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("238")).
				Padding(0, 2)

	descriptionBodyStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("252")).
				Padding(1, 2)

	descActionPromptStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("39")).
				Bold(true).
				Padding(1, 2)
)

func (m *model) setViewport() {
	m.viewport = viewport.New(viewport.WithWidth(0), viewport.WithHeight(0))
	m.viewport.YPosition = m.getHeaderLineCount()
	m.viewport.HighlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Background(lipgloss.Color("34"))
	m.viewport.SelectedHighlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Background(lipgloss.Color("47"))
	m.viewport.HighlightNext()
}

func (m model) setDescViewportContent(item item) string {
	heroTitle := descHeroTitleStyle.Render(item.title)

	var priceBadge string
	if item.price == "FREE" {
		priceBadge = m.styles.freeBadge.Render("FREE")
	} else {
		priceBadge = m.styles.paidBadge.Render(item.price)
	}

	var platformBadge string
	if item.isPlugin == true {
		platformBadge = descPlatformBadge.Render("VST3 / AU / LV2")
	}

	metaRow := descMetaStyle.Render(lipgloss.JoinHorizontal(lipgloss.Center, priceBadge, "  ", platformBadge))
	dividerWidth := max(10, m.viewport.Width()-4)
	divider := descDividerStyle.Render(strings.Repeat("─", dividerWidth))

	description := descriptionBodyStyle.Width(m.viewport.Width()).Render(item.description)
	actionPrompt := descActionPromptStyle.Render("↵ Press Enter to View & Download")

	view := lipgloss.JoinVertical(lipgloss.Left, heroTitle, metaRow, divider, description, actionPrompt)
	return view
}

func setViewportArrows(vp viewport.Model) (string, string) {
	arrowStyle := lipgloss.NewStyle().Width(vp.Width()).Align(lipgloss.Left).Foreground(lipgloss.Color("244"))

	topIndicator := ""
	bottomIndicator := ""

	if vp.TotalLineCount() > vp.Height() {
		if !vp.AtTop() {
			topIndicator = "  ▲ More above"
		}
		if !vp.AtBottom() {
			bottomIndicator = "  ▼ More below"
		}
	}

	topArrow := arrowStyle.Render(topIndicator)
	bottomArrow := arrowStyle.Render(bottomIndicator)

	return topArrow, bottomArrow
}

func (m *model) syncDescViewport() {
	if selected := m.list.SelectedItem(); selected != nil {
		if item, ok := selected.(item); ok {
			if m.lastSelected == "" || item.title != m.lastSelected {
				m.lastSelected = item.title
				m.viewport.SetContent(m.setDescViewportContent(item))
				m.viewport.GotoTop()
			}
		}
	}
}
