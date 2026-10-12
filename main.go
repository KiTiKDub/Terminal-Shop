package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	_ "image/jpeg"
	_ "image/png"
)

var (
	descriptionStyle = lipgloss.NewStyle().Padding(1, 2).MarginTop(1)
)

type styles struct {
	title        lipgloss.Style
	item         lipgloss.Style
	selectedItem lipgloss.Style
	pagination   lipgloss.Style
	help         lipgloss.Style
	quitText     lipgloss.Style
}

func setStyle(darkBG bool) styles {
	var style styles
	style.title = lipgloss.NewStyle().MarginLeft(2)
	style.item = lipgloss.NewStyle().PaddingLeft(4)
	style.selectedItem = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	style.pagination = list.DefaultStyles(darkBG).PaginationStyle.PaddingLeft(4)
	style.help = list.DefaultStyles(darkBG).HelpStyle.PaddingLeft(4).PaddingBottom(1)
	style.quitText = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	return style
}

type itemDelegate struct {
	styles *styles
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

type model struct {
	list     list.Model
	choice   string
	styles   styles
	quitting bool
}

type item struct {
	title       string
	description string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return "" }

func initialModel() model {
	items := []list.Item{
		item{"Disburser", getDisburserDesc()},
		item{"QwikRef", getQwikRefDesc()},
		item{"ChopChop!", getChopChopDesc()},
		item{"Mangl3r", getMangl3rDesc()},
		item{"Distork", getDistorkDesc()},
		item{"Krush", getKrushDesc()},
		item{"KiTiK Kollection Vol 1", getKollectionV1Desc()},
	}

	const defaultWidth = 20
	const listHeight = 14

	d := list.NewDefaultDelegate()
	d.ShowDescription = false

	l := list.New(items, d, defaultWidth, listHeight)
	l.Title = "KiTiK Store"
	m := model{list: l}
	m.updateStyles(true)

	return m
}

func (m *model) updateStyles(isDark bool) {
	m.styles = setStyle(isDark)
	m.list.Styles.Title = m.styles.title
	m.list.Styles.PaginationStyle = m.styles.pagination
	m.list.Styles.HelpStyle = m.styles.help
	m.list.SetDelegate(itemDelegate{styles: &m.styles})
}

func (m model) descView() string {
	selectedItem := m.list.SelectedItem()
	if selectedItem == nil {
		return descriptionStyle.Render("No Item selected")
	}

	item, ok := selectedItem.(item)
	if !ok {
		return ""
	}

	contentWidth := m.list.Width() * 2

	return descriptionStyle.Width(contentWidth).Render(item.description)
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width / 3)
		return m, nil

	case tea.KeyPressMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i.title)
			}
			return m, nil
		case "backspace":
			m.choice = ""
			return m, nil
		}

	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() tea.View {

	if m.quitting {
		return tea.NewView(m.styles.quitText.Render("Thanks for shopping!"))
	}
	if m.choice != "" {
		v := tea.NewView(openPluginView(m))
		v.AltScreen = true
		return v
	}

	var s strings.Builder
	s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, m.list.View(), m.descView()))

	v := tea.NewView("\n" + s.String())
	v.AltScreen = true

	return v
}

func openPluginView(m model) string {

	var msg string

	switch m.choice {
	case "Disburser":
		msg = "Price: FREE\n\n" + renderANSIHalfBlocks("Disburser.jpg", m.list.Width()*3) //change width calculation to be straignt from teamsg eventually
	case "QwikRef":
		msg = "Price: FREE\n\n" + renderANSIHalfBlocks("QwikRef.png", m.list.Width()*3)
	case "ChopChop!":
		msg = "Price: FREE\n\n" + renderANSIHalfBlocks("chopchop.jpg", m.list.Width()*3)
	case "Mangl3r":
		msg = "Price: FREE\n\n" + renderANSIHalfBlocks("Mangl3r.jpg", m.list.Width()*3)
	case "Distork":
		msg = "Price: FREE\n\n" + renderANSIHalfBlocks("Distork.png", m.list.Width()*3) + "\n" + strconv.Itoa(m.list.Width())
	case "Krush":
		msg = "Price: FREE\n\n" + renderANSIHalfBlocks("KITIK_LOGO_NO_BKGD.png", m.list.Width()*3)
	case "KiTiK Kollection Vol 1":
		msg = "Price: FREE\n\n" + renderANSIHalfBlocks("kol.jpg", m.list.Width()*3)
	default:
		msg = "default" + m.choice
	}

	return msg
}

func main() {
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
