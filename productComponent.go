package main

import (
	"strconv"

	"charm.land/bubbles/v2/viewport"
	"charm.land/lipgloss/v2"
)

var (
	horzListStyle = lipgloss.NewStyle().Align(lipgloss.Center).Foreground(lipgloss.BrightRed)
)

func (m *model) setPluginViewportContent() {

	var msg string
	vpWidth := m.pluginViewport.Width()
	instructions := horzListStyle.Width(vpWidth).Render("Download straight from the server or from Gumroad")
	downloadList := horzListStyle.Width(vpWidth).Render(m.downloadList.View())

	switch m.choice {
	case "Disburser":
		msg = renderANSIHalfBlocks("Disburser.jpg", vpWidth)
	case "QwikRef":
		msg = renderANSIHalfBlocks("QwikRef.png", vpWidth)
	case "ChopChop":
		msg = renderANSIHalfBlocks("chopchop.jpg", vpWidth)
	case "Mangl3r":
		msg = renderANSIHalfBlocks("Mangl3r.jpg", vpWidth)
	case "Distork":
		msg = renderANSIHalfBlocks("Distork.png", vpWidth) + "\n" + strconv.Itoa(vpWidth)
	case "Krush":
		msg = renderANSIHalfBlocks("KITIK_LOGO_NO_BKGD.png", vpWidth)
	case "KiTiK Kollection Vol 1":
		msg = renderANSIHalfBlocks("kol.jpg", vpWidth)
	default:
		msg = "default" + m.choice
	}

	m.pluginViewport.SetContent(instructions + "\n\n" + downloadList + "\n\n" + msg)
	m.pluginViewport.GotoTop()
}

func (m *model) setPluginViewport() {
	m.pluginViewport = viewport.New(viewport.WithWidth(0), viewport.WithHeight(0))
	m.pluginViewport.YPosition = m.getHeaderLineCount()

	m.pluginViewport.HighlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Background(lipgloss.Color("34"))
	m.pluginViewport.SelectedHighlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Background(lipgloss.Color("47"))
	m.pluginViewport.HighlightNext()
}
