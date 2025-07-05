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
	if h.host.IPAddress != "" && h.host.IPAddress != h.host.Hostname {
		return fmt.Sprintf("%s (%s@%s:%d) [%s]", h.host.Name, h.host.Username, h.host.Hostname, h.host.Port, h.host.IPAddress)
	}
	return fmt.Sprintf("%s (%s@%s:%d)", h.host.Name, h.host.Username, h.host.Hostname, h.host.Port)
}

func (h hostItem) Description() string {
	desc := ""
	if h.host.Description != "" {
		desc = h.host.Description
	}
	if h.host.Tags != "" {
		if desc != "" {
			desc += " • "
		}
		desc += "🏷️ " + h.host.Tags
	}
	if h.host.UseCount > 0 {
		if desc != "" {
			desc += " • "
		}
		desc += fmt.Sprintf("Used %d times", h.host.UseCount)
	}
	return desc
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

	// Create list
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "SSH Hosts"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false) // Disable default help to use our custom help
	l.Styles.Title = titleStyle

	// Disable all default key bindings
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
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	messageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B"))
)

func (m Model) Init() tea.Cmd {
	return m.loadHosts()
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
		return fmt.Sprintf(
			"Search SSH Hosts\n\n%s\n\n%s",
			m.searchInput.View(),
			"Press Enter to search, Esc to cancel",
		)
	case confirmDeleteView:
		if m.hostToDelete != nil {
			return fmt.Sprintf(
				"⚠️  Delete Host Confirmation\n\n"+
					"Are you sure you want to delete:\n"+
					"🖥️  %s (%s@%s:%d)\n\n"+
					"⚠️  This will remove the host from:\n"+
					"   • SSH Manager database\n"+
					"   • ~/.ssh/known_hosts file\n\n"+
					"❓ Continue? (y/N)\n\n"+
					"Press 'y' to confirm, 'n' or 'Esc' to cancel",
				m.hostToDelete.Name,
				m.hostToDelete.Username,
				m.hostToDelete.Hostname,
				m.hostToDelete.Port,
			)
		}
		return "Error: No host selected for deletion"
	default:
		content := m.list.View()
		if m.message != "" {
			var msgStyle lipgloss.Style
			if strings.Contains(m.message, "Error") {
				msgStyle = errorStyle
			} else {
				msgStyle = messageStyle
			}
			content += "\n" + msgStyle.Render(m.message)
		}

		// Add custom help
		helpText := "↑/k up • ↓/j down • / search • x delete • r refresh • q quit • ? more"
		content += "\n" + helpText

		return content
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
