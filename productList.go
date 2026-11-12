package main

import (
	"fmt"
	"io"
	"strings"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

type listStyle struct {
	title        lipgloss.Style
	item         lipgloss.Style
	selectedItem lipgloss.Style
	pagination   lipgloss.Style
	help         lipgloss.Style
	quitText     lipgloss.Style
	cursor       lipgloss.Style
	freeBadge    lipgloss.Style
	paidBadge    lipgloss.Style
}

func setListStyle() listStyle {
	var style listStyle
	style.title = lipgloss.NewStyle().MarginLeft(2).Bold(true).Foreground(lipgloss.Color("39"))
	style.item = lipgloss.NewStyle().Padding(0, 1)
	style.selectedItem = lipgloss.NewStyle().
		Foreground(lipgloss.Color("255")).
		Bold(true).
		Padding(0, 1)
	style.cursor = lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true)
	style.freeBadge = lipgloss.NewStyle().
		Background(lipgloss.Color("42")).
		Foreground(lipgloss.Color("16")).
		Bold(true).
		Padding(0, 1)
	style.paidBadge = lipgloss.NewStyle().
		Background(lipgloss.Color("214")).
		Foreground(lipgloss.Color("16")).
		Bold(true).
		Padding(0, 1)
	style.quitText = lipgloss.NewStyle().Margin(1, 0, 2, 4).Foreground(lipgloss.Color("244"))
	return style
}

type itemDelegate struct {
	styles *listStyle
}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	wVal := m.Width()
	style := d.styles.item
	cursor := "  "
	if index == m.Index() {
		style = d.styles.selectedItem
		cursor = d.styles.cursor.Render("▸ ")
	}

	var badge string
	if i.price == "FREE" {
		badge = d.styles.freeBadge.Render("FREE")
	} else {
		badge = d.styles.paidBadge.Render(i.price)
	}
	badgeW := lipgloss.Width(badge)
	frameW := style.GetHorizontalFrameSize()

	availTitle := wVal - frameW - 2 - 1 - badgeW
	var rowContent string
	if availTitle >= 4 {
		t := ansi.Truncate(i.Title(), availTitle, "...")
		gap := wVal - frameW - 2 - lipgloss.Width(t) - badgeW
		if gap < 1 {
			gap = 1
		}
		rowContent = cursor + t + strings.Repeat(" ", gap) + badge
	} else {
		availTitleOnly := max(0, wVal-frameW-2)
		t := ansi.Truncate(i.Title(), availTitleOnly, "...")
		rowContent = cursor + t
	}

	fmt.Fprint(w, style.Width(wVal).Render(rowContent))
}

type item struct {
	title       string
	description string
	price       string
	isPlugin    bool
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return "" }

func (m *model) setProductList() {
	items := []list.Item{
		item{"Disburser", getDisburserDesc(), "FREE", true},
		item{"QwikRef", getQwikRefDesc(), "FREE", true},
		item{"ChopChop", getChopChopDesc(), "FREE", true},
		item{"Mangl3r", getMangl3rDesc(), "FREE", true},
		item{"Distork", getDistorkDesc(), "FREE", true},
		item{"Krush", getKrushDesc(), "FREE", true},
		item{"KiTiK Kollection Vol 1", getKollectionV1Desc(), "$4.99", false},
	}

	const defaultWidth = 20
	const defaultHeight = 14

	d := list.NewDefaultDelegate()
	d.ShowDescription = false

	l := list.New(items, d, defaultWidth, defaultHeight)
	l.Title = "KiTiK Store"
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetShowTitle(false)
	m.list = l
}

func (m *model) updateProductListStyles() {
	m.styles = setListStyle()
	m.list.Styles.Title = m.styles.title
	m.list.Styles.PaginationStyle = m.styles.pagination
	m.list.Styles.HelpStyle = m.styles.help
	m.list.SetDelegate(itemDelegate{styles: &m.styles})
}
