package main

import (
	"strings"
	"testing"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func TestHeaderCenteringAndViewSpecificHelp(t *testing.T) {
	m := initialModel()
	m.updateHeaders(100)

	if !strings.Contains(m.header, "AUDIO PLUGINS & CREATIVE TOOLS") {
		t.Errorf("expected m.header to contain subtitle")
	}
	if !strings.Contains(m.header, "Navigate") || !strings.Contains(m.header, "Switch Focus") {
		t.Errorf("expected m.header to contain navigation hints")
	}

	if !strings.Contains(m.pluginHeader, "AUDIO PLUGINS & CREATIVE TOOLS") {
		t.Errorf("expected m.pluginHeader to contain subtitle")
	}
	if strings.Contains(m.pluginHeader, "Navigate") || strings.Contains(m.pluginHeader, "Switch Focus") {
		t.Errorf("expected m.pluginHeader NOT to contain main navigation hints")
	}

	styles := setListStyle()
	bg := styles.selectedItem.GetBackground()
	if bg != nil && bg != (lipgloss.NoColor{}) {
		t.Errorf("expected selectedItem to have no background color, got %v", bg)
	}
}

func TestListItemTruncationAndBadges(t *testing.T) {
	styles := setListStyle()
	delegate := itemDelegate{styles: &styles}

	longItem := item{title: "KiTiK Kollection Vol 1", description: "desc", price: "$4.99"}
	shortItem := item{title: "Disburser", description: "desc", price: "FREE"}

	l := list.New([]list.Item{longItem, shortItem}, delegate, 24, 10)

	// Test long item selected at width 24
	l.Select(0)
	var bufSelected strings.Builder
	delegate.Render(&bufSelected, l, 0, longItem)
	outSelected := bufSelected.String()

	if !strings.Contains(outSelected, "...") {
		t.Errorf("expected long item to be truncated with '...', got %q", outSelected)
	}
	if !strings.Contains(outSelected, "$4.99") {
		t.Errorf("expected long item to contain price badge '$4.99', got %q", outSelected)
	}
	if lipgloss.Height(outSelected) != 1 {
		t.Errorf("expected height 1, got %d", lipgloss.Height(outSelected))
	}
	if lipgloss.Width(outSelected) > 24 {
		t.Errorf("expected width <= 24, got %d", lipgloss.Width(outSelected))
	}

	// Test short item unselected
	l.Select(0) // keep index 0 selected, render index 1
	var bufUnselected strings.Builder
	delegate.Render(&bufUnselected, l, 1, shortItem)
	outUnselected := bufUnselected.String()

	if !strings.Contains(outUnselected, "Disburser") {
		t.Errorf("expected unselected item to contain 'Disburser', got %q", outUnselected)
	}
	if !strings.Contains(outUnselected, "FREE") {
		t.Errorf("expected unselected item to contain 'FREE', got %q", outUnselected)
	}
	if strings.Contains(outUnselected, "...") {
		t.Errorf("short item should not contain '...', got %q", outUnselected)
	}
	if lipgloss.Height(outUnselected) != 1 {
		t.Errorf("expected height 1, got %d", lipgloss.Height(outUnselected))
	}
}

func TestBuildMainViewWithinTerminalBounds(t *testing.T) {
	testCases := []struct {
		name   string
		width  int
		height int
	}{
		{name: "Standard 80x24", width: 80, height: 24},
		{name: "Wide 120x40", width: 120, height: 40},
		{name: "Large 160x50", width: 160, height: 50},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := initialModel()
			updatedModel, _ := m.Update(tea.WindowSizeMsg{
				Width:  tc.width,
				Height: tc.height,
			})
			mod := updatedModel.(model)

			viewStr := buildMainView(mod)
			renderedHeight := lipgloss.Height(viewStr)
			renderedWidth := lipgloss.Width(viewStr)

			t.Logf("[%s] headerHeight=%d, headerWidth=%d, listHeight=%d, viewportHeight=%d, listWidth=%d, viewportWidth=%d",
				tc.name, mod.getHeaderLineCount(), lipgloss.Width(mod.header), mod.list.Height(), mod.viewport.Height(), mod.list.Width(), mod.viewport.Width())
			t.Logf("[%s] list.View() height=%d, width=%d", tc.name, lipgloss.Height(mod.list.View()), lipgloss.Width(mod.list.View()))
			t.Logf("[%s] viewport.View() height=%d, width=%d", tc.name, lipgloss.Height(mod.viewport.View()), lipgloss.Width(mod.viewport.View()))

			if renderedHeight > tc.height {
				t.Errorf("rendered height %d exceeds terminal height %d", renderedHeight, tc.height)
			}
			if renderedWidth > tc.width {
				t.Errorf("rendered width %d exceeds terminal width %d", renderedWidth, tc.width)
			}

			// Verify left and right boxes have matching outer heights
			topArrow, bottomArrow := setViewportArrows(mod.viewport)
			desc := lipgloss.JoinVertical(lipgloss.Left, topArrow, mod.viewport.View(), bottomArrow)
			leftStyle, rightStyle := setActiveBorder(mod.focus)
			borderH := activeBorderStyle.GetHorizontalBorderSize()

			leftBox := leftStyle.Width(mod.list.Width() + borderH).Render(mod.list.View())
			rightBox := rightStyle.Width(mod.viewport.Width() + borderH).Render(desc)

			leftBoxHeight := lipgloss.Height(leftBox)
			rightBoxHeight := lipgloss.Height(rightBox)

			if leftBoxHeight != rightBoxHeight {
				t.Errorf("left box height (%d) does not match right box height (%d)", leftBoxHeight, rightBoxHeight)
			}

			headerHeight := mod.getHeaderLineCount()
			expectedBodyHeight := max(0, tc.height-headerHeight)
			if leftBoxHeight != expectedBodyHeight {
				t.Errorf("box height (%d) does not match expected body height (%d)", leftBoxHeight, expectedBodyHeight)
			}
		})
	}
}

func TestBuildPluginViewWithinTerminalBounds(t *testing.T) {
	testCases := []struct {
		name   string
		width  int
		height int
	}{
		{name: "Standard 80x24", width: 80, height: 24},
		{name: "Wide 120x40", width: 120, height: 40},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := initialModel()
			m.choice = "Disburser"
			updatedModel, _ := m.Update(tea.WindowSizeMsg{
				Width:  tc.width,
				Height: tc.height,
			})
			mod := updatedModel.(model)

			viewStr := buildPluginView(mod)
			renderedHeight := lipgloss.Height(viewStr)
			renderedWidth := lipgloss.Width(viewStr)

			if renderedHeight > tc.height {
				t.Errorf("rendered height %d exceeds terminal height %d", renderedHeight, tc.height)
			}
			if renderedWidth > tc.width {
				t.Errorf("rendered width %d exceeds terminal width %d", renderedWidth, tc.width)
			}
		})
	}
}
