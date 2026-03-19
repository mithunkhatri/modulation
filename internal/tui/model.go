package tui

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mithunkhatri/modulation/internal/audio"
	"github.com/mithunkhatri/modulation/internal/radio"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("6")). // Bright Cyan
			Padding(0, 1)

	taglineStyle = lipgloss.NewStyle().
			Italic(true).
			Foreground(lipgloss.Color("13"))

	headerStyle = lipgloss.NewStyle().
			MarginBottom(1).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(lipgloss.Color("#3C3C3C"))

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#7D56F4")).
			Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EEEEEE"))

	footerStyle = lipgloss.NewStyle().
			MarginTop(1).
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(lipgloss.Color("#3C3C3C")).
			Foreground(lipgloss.Color("#888888"))

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Bold(true)
)

var Version = "dev"

type Model struct {
	player       *audio.Player
	stations     []radio.Station
	cursor       int
	selected     int
	err          error
	width        int
	height       int
	loading      bool
	Background   bool
	errorMessage string
	searchInput  textinput.Model
	filterInput  textinput.Model
	searching    bool
	filtering    bool
	favorites    []radio.Station
	viewMode     int // 0: browse, 1: favorites
	minimized    bool
	spinner      spinner.Model
	isBuffering  bool
}

const (
	viewBrowse = 0
	viewFavs   = 1
)

func NewModel(p *audio.Player) Model {
	ti := textinput.New()
	ti.Placeholder = "Search stations..."
	ti.CharLimit = 156
	ti.Width = 20

	fi := textinput.New()
	fi.Placeholder = "Filter by category (tag)..."
	fi.CharLimit = 64
	fi.Width = 20

	favs, _ := radio.LoadFavorites()

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))

	return Model{
		player:      p,
		cursor:      0,
		selected:    -1,
		loading:     true,
		searchInput: ti,
		filterInput: fi,
		favorites:   favs,
		viewMode:    viewBrowse,
		spinner:     s,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		fetchStations("", ""),
		m.spinner.Tick,
	)
}

type stationsMsg []radio.Station
type errorMsg error

func fetchStations(query string, tag string) tea.Cmd {
	return func() tea.Msg {
		stations, err := radio.GetStations(query, tag, 50)
		if err != nil {
			return errorMsg(err)
		}
		return stationsMsg(stations)
	}
}

type playMsg struct {
	err error
}

func (m Model) playCmd(url string, name string) tea.Cmd {
	return func() tea.Msg {
		err := m.player.Play(url, name)
		return playMsg{err: err}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.minimized {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "b" || msg.String() == "esc" {
				m.minimized = false
			}
			if msg.String() == "q" {
				return m, tea.Quit
			}
		}
		return m, nil
	}

	if m.searching {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				m.searching = false
				m.loading = true
				m.cursor = 0
				return m, fetchStations(m.searchInput.Value(), m.filterInput.Value())
			case "esc":
				m.searching = false
				m.searchInput.Blur()
				return m, nil
			}
		}

		var cmd tea.Cmd
		m.searchInput, cmd = m.searchInput.Update(msg)
		return m, cmd
	}

	if m.filtering {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				m.filtering = false
				m.loading = true
				m.cursor = 0
				return m, fetchStations(m.searchInput.Value(), m.filterInput.Value())
			case "esc":
				m.filtering = false
				m.filterInput.Blur()
				return m, nil
			}
		}

		var cmd tea.Cmd
		m.filterInput, cmd = m.filterInput.Update(msg)
		return m, cmd
	}

	// Active list based on view mode
	activeList := m.stations
	if m.viewMode == viewFavs {
		activeList = m.favorites
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "/", "f3":
			m.searching = true
			m.filtering = false
			m.searchInput.Focus()
			return m, nil
		case "f":
			m.filtering = true
			m.searching = false
			m.filterInput.Focus()
			return m, nil
		case "q", "ctrl+c", "f10":
			m.player.Stop()
			return m, tea.Quit
		case "enter", " ":
			m.selected = m.cursor
			if len(activeList) > 0 {
				s := activeList[m.selected]
				if m.player.IsCurrentStation(s.URL) {
					m.player.Pause()
				} else {
					m.isBuffering = true
					return m, tea.Batch(
						m.playCmd(s.URL, s.Name),
						m.spinner.Tick,
					)
				}
			}
		case "up", "k":
			m.errorMessage = ""
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			m.errorMessage = ""
			if m.cursor < len(activeList)-1 {
				m.cursor++
			}
		case "p":
			m.player.Pause()
		case "s":
			m.player.Stop()
			m.selected = -1
		case "b":
			m.minimized = true
			m.searching = false
			m.filtering = false
			return m, nil
		case "v":
			m.viewMode = (m.viewMode + 1) % 2
			m.cursor = 0
			m.selected = -1
		case "a":
			if len(activeList) > 0 {
				s := activeList[m.cursor]
				m.toggleFavorite(s)
			}
		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			idx, _ := strconv.Atoi(msg.String())
			if idx-1 < len(activeList) {
				m.cursor = idx - 1
				m.selected = m.cursor
				s := activeList[m.selected]
				m.player.Play(s.URL, s.Name)
			}
		case "c":
			m.errorMessage = ""
			m.err = nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case stationsMsg:
		m.stations = msg
		m.loading = false

	case errorMsg:
		m.errorMessage = fmt.Sprintf("Error: %v", msg)
		m.loading = false
		// Don't set m.err = msg as that's used for fatal errors in our current View logic

	case playMsg:
		m.isBuffering = false
		if msg.err != nil {
			m.errorMessage = fmt.Sprintf("Playback Error: %v", msg.err)
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	if m.minimized {
		spinnerStr := ""
		if m.isBuffering {
			spinnerStr = m.spinner.View() + " "
		}
		return fmt.Sprintf("\n  %s[Modulation] Playing: %s\n  (Minimized Mode) Press 'b' to return to TUI, 'q' to quit.", spinnerStr, m.getPlayerStatus())
	}

	if m.err != nil {
		return fmt.Sprintf("\n  Error: %v\n\n  Press 'c' to continue and choose a different station.", m.err)
	}

	var s strings.Builder
	activeList := m.stations
	if m.viewMode == viewFavs {
		activeList = m.favorites
	}

	title := "MODULATION"
	if m.viewMode == viewFavs {
		title = "MODULATION (FAVORITES)"
	}

	status := m.getPlayerStatus()
	if m.isBuffering {
		status = m.spinner.View() + " Buffering..."
	}

	titleStr := titleStyle.Render(title)
	taglineStr := taglineStyle.Render(" Resonating through your favorite terminal")
	statusLine := fmt.Sprintf("  Playing: %s", status)
	errorLine := m.getErrorMessage()

	// Simple header building
	s.WriteString("\n")
	s.WriteString(titleStr + "\n")
	s.WriteString(taglineStr + "\n\n")
	s.WriteString(statusLine + "\n")
	if errorLine != "" {
		s.WriteString(errorLine + "\n")
	}

	if m.searching {
		s.WriteString(fmt.Sprintf("  Search: %s\n", m.searchInput.View()))
	}
	if m.filtering {
		s.WriteString(fmt.Sprintf("  Filter (Category): %s\n", m.filterInput.View()))
	}

	s.WriteString("────────────────────────────────────────────────\n")

	// Body
	if m.loading && m.viewMode == viewBrowse {
		s.WriteString(fmt.Sprintf("  %s Loading stations...\n", m.spinner.View()))
	} else {
		stations := activeList
		if len(stations) == 0 {
			if m.viewMode == viewFavs {
				s.WriteString("  No favorites yet. Add some with 'a' while browsing.\n")
			} else {
				s.WriteString("  No stations found.\n")
			}
		} else {
			for i, station := range stations {
				isFav := m.isFavorite(station)
				favChar := " "
				if isFav {
					favChar = "★"
					if runtime.GOOS == "windows" {
						favChar = "*" // Safer for standard Windows CMD/PowerShell
					}
				}

				cursor := " "
				if m.cursor == i {
					cursor = ">"
				}

				indexStr := fmt.Sprintf("%2d.", i+1)
				line := fmt.Sprintf("%s %s %s %-30s | %s | %s", cursor, indexStr, favChar, station.Name, station.Country, station.Tags)

				if m.cursor == i {
					s.WriteString(selectedStyle.Render(line) + "\n")
				} else {
					if m.selected == i {
						s.WriteString(statusStyle.Render(line) + "\n")
					} else {
						s.WriteString(normalStyle.Render(line) + "\n")
					}
				}

				// Limit list to height
				if i > m.height-10 {
					break
				}
			}
		}
	}

	// Footer
	help := " q: quit | enter/space: play | v: toggle favs | a: fav/unfav | p: pause | s: stop | 1-9: select | /: search | f: filter | c: clear | b: background "
	goStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Bold(true)
	credits := fmt.Sprintf(" Modulation (%s) | Made with %s by %s ", Version, goStyle.Render("Go"), goStyle.Render("Mithun Khatri"))

	footer := lipgloss.JoinVertical(lipgloss.Left,
		help,
		credits,
	)
	s.WriteString(footerStyle.Render(footer))

	return s.String()
}

func (m Model) getPlayerStatus() string {
	if m.player.IsPlaying() {
		return statusStyle.Render(m.player.CurrentStation())
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Render("Stopped")
}

func (m Model) getErrorMessage() string {
	if m.errorMessage != "" {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Bold(true).Render(fmt.Sprintf("\n  %s", m.errorMessage))
	}
	return ""
}

func (m Model) isFavorite(s radio.Station) bool {
	for _, f := range m.favorites {
		if f.URL == s.URL {
			return true
		}
	}
	return false
}

func (m *Model) toggleFavorite(s radio.Station) {
	isFav := false
	idx := -1
	for i, f := range m.favorites {
		if f.URL == s.URL {
			isFav = true
			idx = i
			break
		}
	}

	if isFav {
		m.favorites = append(m.favorites[:idx], m.favorites[idx+1:]...)
	} else {
		m.favorites = append(m.favorites, s)
	}
	radio.SaveFavorites(m.favorites)
}
