package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	_ "image/jpeg"
	_ "image/png"

	"charm.land/ssh"
	"charm.land/wish/v2"
	"charm.land/wish/v2/activeterm"
	btm "charm.land/wish/v2/bubbletea"
	"charm.land/wish/v2/logging"

	"github.com/muesli/termenv"
)

type focusPane int

const (
	listPane focusPane = iota
	viewportPane
	pluginPane
)

const (
	// host = "localhost"
	sshPort  = "2222"
	httpPort = "8080"
)

var (
	activeBorderStyle   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("39"))
	inactiveBorderStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("238"))
)

type model struct {
	files          []string
	session        ssh.Session
	serverHost     string
	list           list.Model
	downloadList   HorizontalList
	viewport       viewport.Model
	pluginViewport viewport.Model
	lastSelected   string
	choice         string
	downloadLink   string
	styles         listStyle
	quitting       bool
	focus          focusPane
	header         string
	pluginHeader   string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:

		m.updateHeaders(msg.Width)

		headerHeight := m.getHeaderLineCount()
		pluginHeaderHeight := m.getPluginHeaderLineCount()

		bodyHeight := max(0, msg.Height-headerHeight)
		pluginBodyHeight := max(0, msg.Height-pluginHeaderHeight)

		borderV := activeBorderStyle.GetVerticalBorderSize()
		borderH := activeBorderStyle.GetHorizontalBorderSize()
		arrowHeight := 2

		listWidth := msg.Width / 3
		viewportWidth := msg.Width - listWidth

		listHeight := max(0, bodyHeight-borderV)
		viewportHeight := max(0, bodyHeight-borderV-arrowHeight)

		m.list.SetWidth(max(0, listWidth-borderH))
		m.list.SetHeight(listHeight)

		m.viewport.SetWidth(max(0, viewportWidth-borderH))
		m.viewport.SetHeight(viewportHeight)
		m.pluginViewport.SetWidth(msg.Width)
		m.pluginViewport.SetHeight(max(0, pluginBodyHeight-arrowHeight))

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
				m.focus = pluginPane

				selectedFile := m.files[m.list.Cursor()]
				token := createDownloadToken(selectedFile)
				url := fmt.Sprintf("http://%s:%s/dl/%s", m.serverHost, httpPort, token)

				label := fmt.Sprintf("👉 Ctrl+Click here to download %s", filepath.Base(selectedFile))
				m.downloadLink = termenv.Hyperlink(url, label)

				m.addDownloadLink(m.downloadLink)
				m.setPluginViewportContent()
				return m, nil
			}

		case "left":
			if m.choice != "" {
				m.downloadList.MoveLeft()
				m.setPluginViewportContent()
			}

		case "right":
			if m.choice != "" {
				m.downloadList.MoveRight()
				m.setPluginViewportContent()
			}

		case "backspace", "esc":
			m.choice = ""
			m.focus = listPane
			return m, nil

		case "tab":
			if m.focus != pluginPane {
				if m.focus == listPane {
					m.focus = viewportPane

				} else {
					m.focus = listPane
				}
			}
			return m, nil
		}

	}

	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch m.focus {
	case listPane:
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	case viewportPane:
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	case pluginPane:
		m.pluginViewport, cmd = m.pluginViewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.syncDescViewport()

	return m, tea.Batch(cmds...)
}

func (m model) View() tea.View {

	if m.quitting {
		return tea.NewView(m.styles.quitText.Render("Thanks for shopping!"))
	}
	if m.choice != "" {
		v := tea.NewView(buildPluginView(m))
		v.AltScreen = true
		return v
	}

	v := tea.NewView(buildMainView(m))
	v.AltScreen = true

	return v
}

func setActiveBorder(focus focusPane) (lipgloss.Style, lipgloss.Style) {
	leftStyle := inactiveBorderStyle
	rightStyle := inactiveBorderStyle

	switch focus {
	case listPane:
		leftStyle = activeBorderStyle
	case viewportPane:
		rightStyle = activeBorderStyle
	}

	return leftStyle, rightStyle

}

func buildMainView(m model) string {

	topArrow, bottomArrow := setViewportArrows(m.viewport)
	desc := lipgloss.JoinVertical(lipgloss.Left, topArrow, m.viewport.View(), bottomArrow)

	leftStyle, rightStyle := setActiveBorder(m.focus)
	borderH := activeBorderStyle.GetHorizontalBorderSize()
	leftBox := leftStyle.Width(m.list.Width() + borderH).Render(m.list.View())
	rightBox := rightStyle.Width(m.viewport.Width() + borderH).Render(desc)

	body := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)
	return lipgloss.JoinVertical(lipgloss.Left, m.header, body)
}

func buildPluginView(m model) string {
	topArrow, bottomArrow := setViewportArrows(m.pluginViewport)
	viewportString := lipgloss.JoinVertical(lipgloss.Left, topArrow, m.pluginViewport.View(), bottomArrow)
	return lipgloss.JoinVertical(lipgloss.Left, m.pluginHeader, viewportString)
}

func initialModel() model {

	m := model{}
	m.setProductList()
	m.updateProductListStyles()
	m.initDownloadList()
	m.updateHeaders(80)
	m.setViewport()
	m.setPluginViewport()

	return m
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {

	m := initialModel()
	m.session = s
	m.serverHost = getPublicBaseURL()
	m.files = []string{"./plugin-files/Disburser.zip", "./plugin-files/Qwikref.zip", "./plugin-files/ChopChop.zip", "./plugin-files/Mangl3r.zip", "./plugin-files/Distork.zip", "./plugin-files/Krush.zip", "./plugin-files/KKV1.zip"}

	return m, []tea.ProgramOption{}

}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/dl/", downloadHandler)
	httpServer := &http.Server{
		Addr:    ":" + httpPort,
		Handler: mux,
	}

	var keyMount string
	prd := os.Getenv("FLY_APP_NAME")

	if prd != "" {
		keyMount = "/data/ssh/id_ed25519"
	} else {
		keyMount = "./ssh/id_ed25519"
	}

	s, err := wish.NewServer(
		wish.WithAddress(":"+sshPort),
		wish.WithHostKeyPath(keyMount),
		wish.WithPublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
			return true
		}),
		wish.WithMiddleware(
			btm.Middleware(teaHandler),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Fatalln("Could not create server", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("Starting SSH server on :%s", sshPort)
	log.Printf("Starting HTTP server on :%s", httpPort)

	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Fatalln("Could not start SSH server", "error", err)
		}
	}()

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln("Could not start HTTP Server", "error", err)
		}
	}()

	<-done

	log.Println("Stopping servers...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_ = s.Shutdown(ctx)
	_ = httpServer.Shutdown(ctx)

}
