package cli

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/levanduy/ssh_management/internal/repo"
	"github.com/levanduy/ssh_management/internal/service"
	"github.com/levanduy/ssh_management/internal/ui"
	"github.com/spf13/cobra"
)

var (
	dbPath        string
	hostService   *service.HostService
	autoDiscovery bool = true // Enable auto-discovery by default
)

var rootCmd = &cobra.Command{
	Use:   "sshm",
	Short: "SSH Manager - TUI-based SSH host management",
	Long: `SSH Manager (sshm) is a terminal-based tool for managing SSH connections.
It automatically discovers SSH hosts from ~/.ssh/known_hosts and provides
an interactive TUI interface for browsing and connecting to hosts.

Features:
- Auto-discovery from ~/.ssh/known_hosts
- Interactive TUI for browsing hosts
- Quick SSH connection with usage tracking
- Lightweight and simple`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initializeService()
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Auto-discovery: Check for new SSH hosts automatically
		autoDiscoverHosts()

		// Always launch TUI mode
		launchTUI()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&dbPath, "db", service.GetDefaultDatabasePath(), "Database file path")
	rootCmd.PersistentFlags().BoolVar(&autoDiscovery, "auto-discovery", true, "Enable automatic SSH host discovery from known_hosts")
}

func initializeService() {
	repo, err := repo.NewSQLiteRepo(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize database: %v\n", err)
		os.Exit(1)
	}

	hostService = service.NewHostService(repo)
}

// autoDiscoverHosts automatically discovers SSH hosts from known_hosts
func autoDiscoverHosts() {
	if !autoDiscovery {
		return // Auto-discovery disabled
	}

	// Count hosts before discovery
	hostsBefore, _ := hostService.GetAllHosts()
	beforeCount := len(hostsBefore)

	// Detect from known_hosts - silent mode
	detectFromSSHFilesSilent()

	// Count hosts after discovery and show notification if new hosts found
	hostsAfter, _ := hostService.GetAllHosts()
	afterCount := len(hostsAfter)

	if afterCount > beforeCount {
		newHosts := afterCount - beforeCount
		fmt.Printf("🔍 Auto-discovered %d new SSH host(s)\n", newHosts)
	}
}

func launchTUI() {
	// Check if there are any hosts first
	hosts, err := hostService.GetAllHosts()
	if err != nil {
		fmt.Printf("Error checking hosts: %v\n", err)
		return
	}

	// If no hosts exist, show welcome message
	if len(hosts) == 0 {
		fmt.Println("🎉 Welcome to SSH Manager!")
		fmt.Println("")
		fmt.Println("No SSH hosts found yet.")
		fmt.Println("💡 Connect to some SSH hosts first, then run 'sshm' again.")
		fmt.Println("   SSH Manager will auto-discover hosts from ~/.ssh/known_hosts")
		return
	}

	// Launch TUI mode
	launchTUIInterface()
}

func launchTUIInterface() {
	model := ui.NewModel(hostService)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
