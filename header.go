package main

import (
	"charm.land/lipgloss/v2"
	figure "github.com/common-nighthawk/go-figure"
)

func (m *model) updateHeaders(width int) {
	fig := figure.NewFigure("KiTiK Store", "", true)
	w := max(width, 70)

	leftStyle := lipgloss.NewStyle().Width(w)

	centerStyle := lipgloss.NewStyle().
		Width(w).
		Align(lipgloss.Center)

	figLeft := leftStyle.
		Foreground(lipgloss.Color("39")).
		Bold(true).
		Render(fig.String())

	subCentered := centerStyle.
		Foreground(lipgloss.Color("244")).
		Italic(true).
		Render("⚡ AUDIO PLUGINS & CREATIVE TOOLS ⚡")

	hintCentered := centerStyle.
		Foreground(lipgloss.Color("240")).
		Render("[↑/↓] Navigate  •  [Tab] Switch Focus  •  [Enter] Select  •  [q] Quit")

	m.header = lipgloss.JoinVertical(lipgloss.Left, figLeft, subCentered, hintCentered)
	m.pluginHeader = lipgloss.JoinVertical(lipgloss.Left, figLeft, subCentered, "\n")
}

func (m *model) getHeaderLineCount() int {
	return lipgloss.Height(m.header)
}

func (m *model) getPluginHeaderLineCount() int {
	return lipgloss.Height(m.pluginHeader)
}
