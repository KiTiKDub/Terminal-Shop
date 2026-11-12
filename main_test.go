package main

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func TestListHeight(t *testing.T) {
	m := initialModel()
	m.list.SetWidth(16)
	m.list.SetHeight(7)
	v := m.list.View()
	t.Logf("SetHeight(7) -> list.Height()=%d, lipgloss.Height()=%d", m.list.Height(), lipgloss.Height(v))
	lines := strings.Split(v, "\n")
	for i, l := range lines {
		t.Logf("line %d: %q", i, l)
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
		{name: "Compact 60x15", width: 60, height: 15},
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

			t.Logf("[%s] leftBox height: %d", tc.name, lipgloss.Height(leftBox))
			lines := strings.Split(leftBox, "\n")
			for i, line := range lines {
				t.Logf("leftBox[%d]: %q", i, line)
			}
			listLines := strings.Split(mod.list.View(), "\n")
			for i, line := range listLines {
				t.Logf("listLines[%d]: %q", i, line)
			}
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
		{name: "Compact 60x15", width: 60, height: 15},
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
