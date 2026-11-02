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
		Foreground(lipgloss.Color("245")).
		MarginRight(2)

	selectedStyle := lipgloss.NewStyle().
		Background(lipgloss.Cyan).
		Foreground(lipgloss.White).
		MarginRight(2)

	for i, item := range h.items {
		if i == h.cursor {
			renderedItems = append(renderedItems, selectedStyle.Render(item.Title()))
		} else {
			renderedItems = append(renderedItems, normalStyle.Render(item.Title()))
		}
	}

	// Join items side-by-side instead of top-to-bottom
	return lipgloss.JoinHorizontal(lipgloss.Left, renderedItems...)
}

func (m *model) setDownloadList() {
	m.downloadList.items = []item{
		{"Download", "Download straight from the server", "N/A"},
		{"Gumroad", "Purchase Straight from Gumraod", "N/A"},
	}
}
