package main

import "charm.land/lipgloss/v2"

type HorizontalList struct {
	items  []item
	cursor int
	width  int
}

func (h *HorizontalList) MoveLeft() {
	if h.cursor > 0 {
		h.cursor--
	}
}

func (h *HorizontalList) MoveRight() {
	if h.cursor < len(h.items)-1 {
		h.cursor++
	}
}

func (h HorizontalList) SelectedItem() item {
	return h.items[h.cursor]
}

func (h HorizontalList) View() string {
	var renderedItems []string

	normalStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("236")).
		Foreground(lipgloss.Color("248")).
		Padding(0, 2).
		MarginRight(2)

	selectedStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("39")).
		Foreground(lipgloss.Color("16")).
		Bold(true).
		Padding(0, 2).
		MarginRight(2)

	for i, item := range h.items {
		label := item.Title()
		switch label {
		case "Download":
			label = "⬇ Download Direct"
		case "Gumroad":
			label = "🛒 Gumroad Store"
		}

		if i == h.cursor {
			renderedItems = append(renderedItems, selectedStyle.Render(label))
		} else {
			renderedItems = append(renderedItems, normalStyle.Render(label))
		}
	}

	buttons := lipgloss.JoinHorizontal(lipgloss.Center, renderedItems...)
	hints := lipgloss.NewStyle().
		Foreground(lipgloss.Color("242")).
		Render("\n[← / →] Switch Option  •  [Enter] Confirm  •  [Esc] Back to Catalog")

	return lipgloss.JoinVertical(lipgloss.Center, buttons, hints)
}

func (m *model) initDownloadList() {
	m.downloadList.items = []item{
		// {"Download", "Download straight from the server", "N/A", false},
		{"Coming Soon", "Purchase Straight from Gumraod", "N/A", false},
	}
}

func (m *model) addDownloadLink(link string) {
	newItem := item{link, "Click here to download!", "N/A", false}

	m.downloadList.items = append(m.downloadList.items, newItem)
}
