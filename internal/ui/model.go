package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/levanduy/ssh_management/internal/domain"
	"github.com/levanduy/ssh_management/internal/service"
	"github.com/levanduy/ssh_management/pkg/ssh"
)

type state int

const (
	listView state = iota
	searchView
	connectingView
	confirmDeleteView
)

type Model struct {
	state        state
	list         list.Model
	searchInput  textinput.Model
	hosts        []*domain.Host
	hostService  *service.HostService
	width        int
	height       int
	message      string
	hostToDelete *domain.Host // Host pending deletion
}

type hostItem struct {
	host *domain.Host
}

func (h hostItem) FilterValue() string {
	return h.host.Name + " " + h.host.Hostname + " " + h.host.IPAddress + " " + h.host.Description + " " + h.host.Tags
}

func (h hostItem) Title() string {
	// Host name in white
	name := h.host.Name

	// Connection info in cyan (like in image)
	connInfo := fmt.Sprintf("(%s@%s:%d)", h.host.Username, h.host.Hostname, h.host.Port)

	// IP address in green (like in image) if available and different from hostname
	if h.host.IPAddress != "" && h.host.IPAddress != h.host.Hostname {
		return fmt.Sprintf("%s %s [%s]", name, connInfo, h.host.IPAddress)
	}

	return fmt.Sprintf("%s %s", name, connInfo)
}

func (h hostItem) Description() string {
	var parts []string

	// Description
	if h.host.Description != "" {
		parts = append(parts, h.host.Description)
	}

	// Tags
	if h.host.Tags != "" {
		parts = append(parts, h.host.Tags)
	}

	// Usage count
	if h.host.UseCount > 0 {
		parts = append(parts, fmt.Sprintf("Used %d times", h.host.UseCount))
	}

	return strings.Join(parts, " • ")
}

type keyMap struct {
	Search  key.Binding
	Connect key.Binding
	Delete  key.Binding
	Refresh key.Binding
	Back    key.Binding
	Quit    key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Search, k.Connect, k.Delete, k.Refresh, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Search, k.Connect, k.Delete},
		{k.Refresh, k.Back, k.Quit},
	}
}

var keys = keyMap{
	Search: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "search"),
	),
	Connect: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "connect"),
	),
	Delete: key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("x", "delete"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

func NewModel(hostService *service.HostService) Model {
	// Create search input
	searchInput := textinput.New()
	searchInput.Placeholder = "Search hosts..."
	searchInput.Focus()
	searchInput.CharLimit = 50

	// Create delegate with no background highlight
	delegate := list.NewDefaultDelegate()

	// Selected item - no background, just different text color
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(accentColor).        // Green for selected
		Background(lipgloss.Color("")). // No background
		Bold(true)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(cyanColor).         // Cyan for selected description
		Background(lipgloss.Color("")) // No background

	// Normal items
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.
		Foreground(textColor)
	delegate.Styles.NormalDesc = delegate.Styles.NormalDesc.
		Foreground(mutedColor)

	// Create list with custom delegate
	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = "SSH Hosts"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.Styles.Title = titleStyle

	// Configure key bindings
	l.KeyMap.CursorUp.SetEnabled(true)
	l.KeyMap.CursorDown.SetEnabled(true)
	l.KeyMap.NextPage.SetEnabled(true)
	l.KeyMap.PrevPage.SetEnabled(true)
	l.KeyMap.GoToStart.SetEnabled(true)
	l.KeyMap.GoToEnd.SetEnabled(true)

	// Disable conflicting keys
	l.KeyMap.Filter.SetEnabled(false)
	l.KeyMap.ClearFilter.SetEnabled(false)
	l.KeyMap.CancelWhileFiltering.SetEnabled(false)
	l.KeyMap.AcceptWhileFiltering.SetEnabled(false)
	l.KeyMap.Quit.SetEnabled(false)
	l.KeyMap.ForceQuit.SetEnabled(false)
	l.KeyMap.ShowFullHelp.SetEnabled(false)
	l.KeyMap.CloseFullHelp.SetEnabled(false)

	m := Model{
		state:       listView,
		list:        l,
		searchInput: searchInput,
		hostService: hostService,
	}

	return m
}

var (
	// Color scheme matching the image - keep green and cyan
	primaryColor = lipgloss.Color("#5B21B6") // Purple for headers only
	accentColor  = lipgloss.Color("#10B981") // Green (like in image)
	cyanColor    = lipgloss.Color("#06B6D4") // Cyan (like in image)
	warningColor = lipgloss.Color("#F59E0B") // Orange
	errorColor   = lipgloss.Color("#DC2626") // Red
	mutedColor   = lipgloss.Color("#6B7280") // Gray
	textColor    = lipgloss.Color("#F3F4F6") // Light gray
	dimTextColor = lipgloss.Color("#9CA3AF") // Dimmed gray

	// Clean header styles
	titleStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Background(primaryColor).
			Padding(0, 1).
			Bold(true)

	// Simple message styles
	messageStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			MarginTop(1)

	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			MarginTop(1)

	warningStyle = lipgloss.NewStyle().
			Foreground(warningColor).
			MarginTop(1)

	// Clean help style
	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			MarginTop(1)

	// Simple search styles
	searchTitleStyle = lipgloss.NewStyle().
				Foreground(textColor).
				Background(primaryColor).
				Padding(0, 1).
				Bold(true)

	// Simple confirmation styles
	confirmTitleStyle = lipgloss.NewStyle().
				Foreground(textColor).
				Background(warningColor).
				Padding(0, 1).
				Bold(true)
)

func (m Model) Init() tea.Cmd {
	return m.refreshWithDiscovery()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width, msg.Height-4)
		return m, nil

	case hostsLoadedMsg:
		m.hosts = msg.hosts
		items := make([]list.Item, len(m.hosts))
		for i, host := range m.hosts {
			items[i] = hostItem{host: host}
		}
		m.list.SetItems(items)
		m.message = fmt.Sprintf("Loaded %d host(s)", len(m.hosts))
		return m, nil

	case hostConnectedMsg:
		m.message = fmt.Sprintf("Connected to %s", msg.hostName)
		m.state = listView
		return m, m.loadHosts() // Refresh to update usage stats

	case discoveryMsg:
		// Update the list with new hosts
		hosts, err := m.hostService.GetAllHosts()
		if err != nil {
			return m, nil
		}
		items := make([]list.Item, len(hosts))
		for i, host := range hosts {
			items[i] = hostItem{host: host}
		}
		m.list.SetItems(items)
		m.hosts = hosts
		m.message = fmt.Sprintf("🔍 Auto-discovered %d new host(s)", msg.newHostsCount)
		return m, nil

	case errorMsg:
		m.message = fmt.Sprintf("Error: %s", msg.error)
		return m, nil

	case tea.KeyMsg:
		switch m.state {
		case listView:
			switch {
			case key.Matches(msg, keys.Quit):
				return m, tea.Quit

			case key.Matches(msg, keys.Search):
				m.state = searchView
				m.searchInput.Focus()
				return m, nil

			case key.Matches(msg, keys.Connect):
				selected := m.list.SelectedItem()
				if selected != nil {
					host := selected.(hostItem).host
					return m, m.connectToHost(host)
				}

			case key.Matches(msg, keys.Delete):
				selected := m.list.SelectedItem()
				if selected != nil {
					m.hostToDelete = selected.(hostItem).host
					m.state = confirmDeleteView
					return m, nil
				}

			case key.Matches(msg, keys.Refresh):
				return m, m.refreshWithDiscovery()
			}

			// Update list only if we're in listView and key wasn't handled above
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)

		case searchView:
			switch {
			case key.Matches(msg, keys.Back), key.Matches(msg, keys.Quit):
				m.state = listView
				m.searchInput.Blur()
				m.searchInput.SetValue("")
				return m, m.loadHosts()

			case msg.Type == tea.KeyEnter:
				query := m.searchInput.Value()
				if query != "" {
					return m, m.searchHosts(query)
				}
				m.state = listView
				m.searchInput.Blur()
				return m, nil
			}

			m.searchInput, cmd = m.searchInput.Update(msg)
			cmds = append(cmds, cmd)

		case confirmDeleteView:
			switch {
			case key.Matches(msg, keys.Back), key.Matches(msg, keys.Quit):
				m.state = listView
				m.hostToDelete = nil
				return m, nil

			case msg.Type == tea.KeyEnter, msg.String() == "y", msg.String() == "Y":
				// Confirm deletion
				if m.hostToDelete != nil {
					host := m.hostToDelete
					m.hostToDelete = nil
					m.state = listView
					return m, m.deleteHost(host)
				}
				m.state = listView
				return m, nil

			case msg.String() == "n", msg.String() == "N":
				// Cancel deletion
				m.state = listView
				m.hostToDelete = nil
				return m, nil
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	switch m.state {
	case searchView:
		header := searchTitleStyle.Render("Search SSH Hosts")
		input := m.searchInput.View()
		help := helpStyle.Render("Press Enter to search • Esc to cancel")

		return fmt.Sprintf("%s\n\n%s\n\n%s", header, input, help)

	case confirmDeleteView:
		if m.hostToDelete != nil {
			title := confirmTitleStyle.Render("Delete Host Confirmation")

			hostInfo := fmt.Sprintf(
				"Host: %s\n"+
					"Connection: %s@%s:%d\n"+
					"IP: %s",
				m.hostToDelete.Name,
				m.hostToDelete.Username,
				m.hostToDelete.Hostname,
				m.hostToDelete.Port,
				m.hostToDelete.IPAddress,
			)

			warning := warningStyle.Render(
				"This will remove the host from:\n" +
					"• SSH Manager database\n" +
					"• ~/.ssh/known_hosts file",
			)

			help := helpStyle.Render("Press 'y' to confirm • 'n' or 'Esc' to cancel")

			return fmt.Sprintf("%s\n\n%s\n\n%s\n\nContinue? (y/N)\n\n%s", title, hostInfo, warning, help)
		}
		return errorStyle.Render("Error: No host selected for deletion")

	default:
		// Main list view
		header := titleStyle.Render("SSH Manager")

		// Status bar
		statusText := fmt.Sprintf("Total hosts: %d", len(m.hosts))
		if len(m.hosts) > 0 {
			statusText += fmt.Sprintf(" • Selected: %d", m.list.Index()+1)
		}
		statusBar := helpStyle.Render(statusText)

		// Main content
		content := m.list.View()

		// Message display
		var message string
		if m.message != "" {
			var msgStyle lipgloss.Style
			if strings.Contains(m.message, "Error") || strings.Contains(m.message, "Failed") {
				msgStyle = errorStyle
			} else if strings.Contains(m.message, "Warning") {
				msgStyle = warningStyle
			} else {
				msgStyle = messageStyle
			}
			message = msgStyle.Render(m.message)
		}

		// Help text
		helpText := helpStyle.Render(
			"↑/k up • ↓/j down • / search • enter connect • x delete • r refresh • q quit",
		)

		// Combine elements
		result := header + "\n" + statusBar + "\n\n" + content

		if message != "" {
			result += "\n\n" + message
		}

		result += "\n\n" + helpText

		return result
	}
}

// Commands
type hostsLoadedMsg struct {
	hosts []*domain.Host
}

type hostConnectedMsg struct {
	hostName string
}

type errorMsg struct {
	error string
}

type discoveryMsg struct {
	newHostsCount int
}

func (m Model) loadHosts() tea.Cmd {
	return func() tea.Msg {
		hosts, err := m.hostService.GetAllHosts()
		if err != nil {
			return errorMsg{error: err.Error()}
		}
		return hostsLoadedMsg{hosts: hosts}
	}
}

func (m Model) refreshWithDiscovery() tea.Cmd {
	return func() tea.Msg {
		// First run auto-discovery
		newHostsCount, err := m.hostService.AutoDiscoverFromKnownHosts()
		if err != nil {
			return errorMsg{error: fmt.Sprintf("Auto-discovery failed: %v", err)}
		}

		// Then load all hosts
		hosts, err := m.hostService.GetAllHosts()
		if err != nil {
			return errorMsg{error: err.Error()}
		}

		// If new hosts were discovered, show discovery message
		if newHostsCount > 0 {
			return discoveryMsg{newHostsCount: newHostsCount}
		}

		return hostsLoadedMsg{hosts: hosts}
	}
}

func (m Model) searchHosts(query string) tea.Cmd {
	return func() tea.Msg {
		hosts, err := m.hostService.SearchHosts(query)
		if err != nil {
			return errorMsg{error: err.Error()}
		}
		return hostsLoadedMsg{hosts: hosts}
	}
}

func (m Model) connectToHost(host *domain.Host) tea.Cmd {
	return func() tea.Msg {
		// Update usage stats
		if err := m.hostService.ConnectToHost(host.ID); err != nil {
			return errorMsg{error: fmt.Sprintf("Failed to update stats: %v", err)}
		}

		// Connect via SSH
		if err := ssh.ConnectToHost(host); err != nil {
			return errorMsg{error: fmt.Sprintf("SSH connection failed: %v", err)}
		}

		return hostConnectedMsg{hostName: host.Name}
	}
}

func (m Model) deleteHost(host *domain.Host) tea.Cmd {
	return func() tea.Msg {
		if err := m.hostService.DeleteHostFromBoth(host.ID); err != nil {
			return errorMsg{error: fmt.Sprintf("Failed to delete host: %v", err)}
		}

		// Reload hosts after deletion
		hosts, err := m.hostService.GetAllHosts()
		if err != nil {
			return errorMsg{error: err.Error()}
		}
		return hostsLoadedMsg{hosts: hosts}
	}
}
