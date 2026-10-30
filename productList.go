package main

import (
	"fmt"
	"io"
	"strings"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type listStyle struct {
	title        lipgloss.Style
	item         lipgloss.Style
	selectedItem lipgloss.Style
	pagination   lipgloss.Style
	help         lipgloss.Style
	quitText     lipgloss.Style
}

func setListStyle() listStyle {
	var style listStyle
	style.title = lipgloss.NewStyle().MarginLeft(2)
	style.item = lipgloss.NewStyle().PaddingLeft(4)
	style.selectedItem = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	style.quitText = lipgloss.NewStyle().Margin(1, 0, 2, 4)
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

	str := fmt.Sprintf("%s", i.Title())
	fn := d.styles.item.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return d.styles.selectedItem.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type item struct {
	title       string
	description string
	price       string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return "" }

func initialModel() model {
	items := []list.Item{
		item{"Disburser", getDisburserDesc(), "FREE"},
		item{"QwikRef", getQwikRefDesc(), "FREE"},
		item{"ChopChop!", getChopChopDesc(), "FREE"},
		item{"Mangl3r", getMangl3rDesc(), "FREE"},
		item{"Distork", getDistorkDesc(), "FREE"},
		item{"Krush", getKrushDesc(), "FREE"},
		item{"KiTiK Kollection Vol 1", getKollectionV1Desc(), "$4.99"},
	}

	const defaultWidth = 20
	const defaultHeight = 14

	d := list.NewDefaultDelegate()
	d.ShowDescription = false

	l := list.New(items, d, defaultWidth, defaultHeight)
	l.Title = "KiTiK Store"
	m := model{list: l}
	m.updateListStyles()

	return m
}

func (m *model) updateListStyles() {
	m.styles = setListStyle()
	m.list.Styles.Title = m.styles.title
	m.list.Styles.PaginationStyle = m.styles.pagination
	m.list.Styles.HelpStyle = m.styles.help
	m.list.SetDelegate(itemDelegate{styles: &m.styles})
}
