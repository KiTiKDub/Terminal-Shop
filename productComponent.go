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
	instructions := horzListStyle.Width(m.list.Width() * 3).Render("Download straight from the server or from Gumroad")
	downloadList := horzListStyle.Width(m.list.Width() * 3).Render(m.downloadList.View())

	switch m.choice {
	case "Disburser":
		msg = renderANSIHalfBlocks("Disburser.jpg", m.list.Width()*3) //change width calculation to be straignt from teamsg eventually
	case "QwikRef":
		msg = renderANSIHalfBlocks("QwikRef.png", m.list.Width()*3)
	case "ChopChop":
		msg = renderANSIHalfBlocks("chopchop.jpg", m.list.Width()*3)
	case "Mangl3r":
		msg = renderANSIHalfBlocks("Mangl3r.jpg", m.list.Width()*3)
	case "Distork":
		msg = renderANSIHalfBlocks("Distork.png", m.list.Width()*3) + "\n" + strconv.Itoa(m.list.Width())
	case "Krush":
		msg = renderANSIHalfBlocks("KITIK_LOGO_NO_BKGD.png", m.list.Width()*3)
	case "KiTiK Kollection Vol 1":
		msg = renderANSIHalfBlocks("kol.jpg", m.list.Width()*3)
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
