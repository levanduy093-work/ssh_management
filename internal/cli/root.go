package cli

import (
	"fmt"
	"os"

	"github.com/levanduy/ssh_management/internal/repo"
	"github.com/levanduy/ssh_management/internal/service"
	"github.com/spf13/cobra"
)

var (
	dbPath      string
	hostService *service.HostService
)

var rootCmd = &cobra.Command{
	Use:   "sshm",
	Short: "SSH Manager - Easy SSH host management",
	Long: `SSH Manager (sshm) is a terminal-based tool for managing SSH connections.
It allows you to store, organize, and quickly connect to SSH hosts.

Features:
- Add, edit, and remove SSH hosts
- Interactive TUI for browsing hosts
- Search and filter hosts by name, hostname, or tags
- Quick connection with usage tracking`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initializeService()
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Check if there are any hosts first
		hosts, err := hostService.GetAllHosts()
		if err != nil {
			fmt.Printf("Error checking hosts: %v\n", err)
			return
		}

		// If no hosts exist, show welcome message and suggest adding one
		if len(hosts) == 0 {
			fmt.Println("🎉 Welcome to SSH Manager!")
			fmt.Println("")
			fmt.Println("You don't have any SSH hosts configured yet.")
			fmt.Println("Let's add your first host:")
			fmt.Println("")
			fmt.Println("💡 Quick start:")
			fmt.Printf("   %s add myserver%s\n", "\033[1;36m", "\033[0m")
			fmt.Printf("   %s sshm%s (to launch interactive mode)\n", "\033[1;36m", "\033[0m")
			fmt.Println("")
			fmt.Println("📖 For help: sshm --help")
			return
		}

		// Default to TUI mode when hosts exist
		tuiCmd.Run(cmd, args)
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
}

func initializeService() {
	repo, err := repo.NewSQLiteRepo(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize database: %v\n", err)
		os.Exit(1)
	}

	hostService = service.NewHostService(repo)
}
